/*
 * MIT License
 *
 * Copyright (c) 2025 linux.do
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package payment

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/hibiken/asynq"
	"github.com/linux-do/pay/internal/common"
	"github.com/linux-do/pay/internal/db"
	"github.com/linux-do/pay/internal/logger"
	"github.com/linux-do/pay/internal/model"
	"github.com/linux-do/pay/internal/util"
	"gorm.io/gorm"
)

// HandleMerchantPaymentNotify 处理商户支付回调任务
func HandleMerchantPaymentNotify(ctx context.Context, t *asynq.Task) error {
	// 解析任务参数
	var payload struct {
		OrderID  uint64 `json:"order_id"`
		ClientID string `json:"client_id"`
	}
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		logger.ErrorF(ctx, "解析商户回调任务参数失败: %v", err)
		return fmt.Errorf("解析任务参数失败: %w", err)
	}

	// 查询订单信息
	var order model.Order
	if err := db.DB(ctx).Where("id = ? AND status = ?", payload.OrderID, model.OrderStatusSuccess).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.ErrorF(ctx, "订单[ID:%d]不存在，跳过回调", payload.OrderID)
			return nil
		}
		return fmt.Errorf("查询订单失败: %w", err)
	}

	// 查询商户API Key信息
	var apiKey model.MerchantAPIKey
	if err := db.DB(ctx).Where("client_id = ?", payload.ClientID).First(&apiKey).Error; err != nil {
		logger.ErrorF(ctx, "查询商户[ClientID:%s]失败: %v", payload.ClientID, err)
		return fmt.Errorf("查询商户信息失败: %w", err)
	}

	// 构建回调参数
	callbackParams := map[string]string{
		"pid":          payload.ClientID,
		"trade_no":     strconv.FormatUint(order.ID, 10),
		"out_trade_no": order.MerchantOrderNo,
		"type":         common.PayTypeEPay,
		"name":         order.OrderName,
		"money":        order.Amount.Truncate(2).StringFixed(2),
		"trade_status": "TRADE_SUCCESS",
		"sign_type":    "MD5",
	}

	callbackParams["sign"] = GenerateSignature(callbackParams, apiKey.ClientSecret)

	if err := sendCallbackRequest(ctx, apiKey.NotifyURL, callbackParams); err != nil {
		retried, _ := asynq.GetRetryCount(ctx)
		maxRetry := 5

		logger.ErrorF(ctx, "商户回调失败: 订单[ID:%d] 重试次数[%d/%d] 错误: %v",
			payload.OrderID, retried+1, maxRetry, err)

		if retried >= maxRetry-1 {
			logger.ErrorF(ctx, "商户回调达到最大重试次数，回调最终失败: 订单[ID:%d]", payload.OrderID)
			return nil // 任务完成（虽然失败）
		}

		return fmt.Errorf("商户回调失败: %w", err)
	}

	logger.InfoF(ctx, "商户回调成功: 订单[ID:%d] ClientID[%s]", payload.OrderID, payload.ClientID)
	return nil
}

// sendCallbackRequest 发送HTTP回调请求
func sendCallbackRequest(ctx context.Context, callbackURL string, params map[string]string) error {
	vals := url.Values{}
	for k, v := range params {
		vals.Add(k, v)
	}

	// 拼接URL
	separator := "?"
	if strings.Contains(callbackURL, "?") {
		separator = "&"
	}
	targetURL := callbackURL + separator + vals.Encode()

	headers := map[string]string{
		"User-Agent": "LinuxDo-Pay/1.0",
	}

	resp, err := util.Request(ctx, http.MethodGet, targetURL, nil, headers, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("回调返回异常状态码: %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %w", err)
	}

	responseText := strings.TrimSpace(strings.ToLower(string(respBody)))
	if responseText != "success" {
		return fmt.Errorf("回调返回非成功响应: %s", string(respBody))
	}

	logger.InfoF(ctx, "商户回调请求成功: URL[%s] 响应[%s]", callbackURL, string(respBody))
	return nil
}
