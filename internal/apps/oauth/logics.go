/*
Copyright 2025 linux.do

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/linux-do/credit/internal/common"
	"github.com/linux-do/credit/internal/config"
	"github.com/linux-do/credit/internal/db"
	"github.com/linux-do/credit/internal/logger"
	"github.com/linux-do/credit/internal/model"
	"github.com/linux-do/credit/internal/otel_trace"
	"go.opentelemetry.io/otel/codes"
	"gorm.io/gorm"
)

func GetUserIDFromSession(s sessions.Session) uint64 {
	userID, ok := s.Get(UserIDKey).(uint64)
	if !ok {
		return 0
	}
	return userID
}

func GetUserIDFromContext(c *gin.Context) uint64 {
	session := sessions.Default(c)
	return GetUserIDFromSession(session)
}

func doOAuth(ctx context.Context, code string) (*model.User, error) {
	// init trace
	ctx, span := otel_trace.Start(ctx, "OAuth")
	defer span.End()

	// get token
	token, err := oauthConf.Exchange(ctx, code)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	var userInfo *model.OAuthUserInfo

	// 先尝试用 OIDC ID Token
	if rawIDToken, ok := token.Extra("id_token").(string); ok && oidcVerifier != nil {
		idToken, err := oidcVerifier.Verify(ctx, rawIDToken)
		if err != nil {
			logger.WarnF(ctx, "Failed to verify ID token, falling back to UserEndpoint: %v", err)
		} else {
			var claims model.OAuthUserInfo
			if err := idToken.Claims(&claims); err != nil {
				logger.WarnF(ctx, "Failed to parse ID token claims, falling back to UserEndpoint: %v", err)
			} else {
				// 从 sub 填充 ID
				if err := claims.FillIdFromSub(); err != nil {
					logger.WarnF(ctx, "Failed to parse sub claim, falling back to UserEndpoint: %v", err)
				} else {
					userInfo = &claims
					logger.InfoF(ctx, "Successfully authenticated user via OIDC: %s (ID: %d)", userInfo.Username, userInfo.Id)
				}
			}
		}
	}

	// 如果 OIDC 不行，回退到 UserEndpoint 方式
	if userInfo == nil {
		logger.InfoF(ctx, "Using UserEndpoint fallback for authentication")
		client := oauthConf.Client(ctx, token)
		resp, err := client.Get(config.Config.OAuth2.UserEndpoint)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
		defer func(body io.ReadCloser) { _ = resp.Body.Close() }(resp.Body)

		// load user info
		responseData, err := io.ReadAll(resp.Body)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
		var fallbackUserInfo model.OAuthUserInfo
		if err = json.Unmarshal(responseData, &fallbackUserInfo); err != nil {
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
		userInfo = &fallbackUserInfo
	}
	if !userInfo.Active {
		err = errors.New(common.BannedAccount)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// 处理用户信息同步逻辑
	var user model.User

	txByUsername := db.DB(ctx).Where("username = ?", userInfo.Username).First(&user)
	if txByUsername.Error != nil {
		txByID := user.GetByID(db.DB(ctx), userInfo.Id)
		if txByID == nil {
			// ID 存在但 username 不匹配(用户改名)
			if err = user.CheckActive(); err != nil {
				span.SetStatus(codes.Error, err.Error())
				return nil, err
			}
			user.UpdateFromOAuthInfo(userInfo)
			if err = db.DB(ctx).Save(&user).Error; err != nil {
				span.SetStatus(codes.Error, err.Error())
				return nil, err
			}
		} else if errors.Is(txByUsername.Error, gorm.ErrRecordNotFound) {
			// ID 和 username 都不存在(全新用户)
			user = model.User{}
			if err = user.CreateWithInitialCredit(ctx, userInfo); err != nil {
				span.SetStatus(codes.Error, err.Error())
				return nil, err
			}
		} else {
			// query failed
			span.SetStatus(codes.Error, txByUsername.Error.Error())
			return nil, txByUsername.Error
		}
	} else {
		if user.ID != userInfo.Id {
			// username 相同但 ID 不同(账户注销后被新用户占用)
			if err = user.CreateWithInitialCredit(ctx, userInfo); err != nil {
				span.SetStatus(codes.Error, err.Error())
				return nil, err
			}
		} else {
			if err = user.CheckActive(); err != nil {
				span.SetStatus(codes.Error, err.Error())
				return nil, err
			}
			user.UpdateFromOAuthInfo(userInfo)
			if err = db.DB(ctx).Save(&user).Error; err != nil {
				span.SetStatus(codes.Error, err.Error())
				return nil, err
			}
		}
	}
	return &user, nil
}
