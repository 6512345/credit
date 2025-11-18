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
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package listener

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/linux-do/pay/internal/config"
	"github.com/linux-do/pay/internal/db"
	"github.com/linux-do/pay/internal/logger"
	"github.com/linux-do/pay/internal/model"
)

// StartExpireListener 启动过期监听器
func StartExpireListener(ctx context.Context) error {
	if db.Redis == nil {
		return fmt.Errorf("redis client is not initialized")
	}

	model.ExpirePendingOrders(ctx)

	// 确保Redis开启了keyspace notifications
	configResult := db.Redis.ConfigSet(ctx, "notify-keyspace-events", "Ex")
	if configResult.Err() != nil {
		log.Printf("[Order Expire Listener] 警告: 设置Redis keyspace notifications失败: %v", configResult.Err())
		log.Printf("[Order Expire Listener] 请手动执行: CONFIG SET notify-keyspace-events Ex")
		return configResult.Err()
	}

	// 订阅过期事件频道
	expiredChannel := fmt.Sprintf("__keyevent@%d__:expired", config.Config.Redis.DB)
	pubSub := db.Redis.PSubscribe(ctx, expiredChannel)

	go func() {
		defer pubSub.Close()
		log.Printf("[Expire Listener] 过期监听器已启动，监听频道: %s", expiredChannel)

		for {
			msg, err := pubSub.ReceiveMessage(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					log.Println("[Order Expire Listener] 监听器已停止")
					return
				}
				logger.ErrorF(ctx, "接收Redis过期事件失败: %v", err)
				continue
			}

			// 处理过期事件
			handleExpiredKey(ctx, msg.Payload)
		}
	}()

	return nil
}

// handleExpiredKey 处理过期的Redis key
func handleExpiredKey(ctx context.Context, expiredKey string) {
	// 只处理订单过期相关的key
	if !strings.HasPrefix(expiredKey, "payment:order:expire:") {
		return
	}

	orderIDStr := strings.TrimPrefix(expiredKey, "payment:order:expire:")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 64)
	if err != nil {
		logger.ErrorF(ctx, "解析订单ID失败: key=%s, error=%v", expiredKey, err)
		return
	}

	// 更新订单状态为过期
	result := db.DB(ctx).Model(&model.Order{}).
		Where("id = ? AND status = ?", orderID, model.OrderStatusPending).
		Update("status", model.OrderStatusExpired)

	if result.Error != nil {
		logger.ErrorF(ctx, "更新订单状态为过期失败: order_id=%d, error=%v", orderID, result.Error)
	} else if result.RowsAffected > 0 {
		logger.InfoF(ctx, "订单已过期: order_id=%d", orderID)
	}
}
