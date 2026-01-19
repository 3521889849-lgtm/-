// Package util 提供通用工具函数
//
// 本文件提供操作日志记录的辅助函数
package util

import (
	"context"
	"example_shop/api/client"
	"example_shop/kitex_gen/hotel"
	"strings"
)

// LogOperation 记录操作日志的辅助函数
//
// 功能：
//   - 异步记录操作日志，不阻塞主业务流程
//   - 自动处理IP地址格式
//   - 支持关联ID（可选）
//   - 失败不影响主流程（日志记录失败时静默处理）
//
// 参数：
//   - ctx: 上下文，用于传递请求信息和超时控制
//   - operatorID: 操作人ID，关联用户账号表
//   - module: 操作模块，如："ROOM"、"ORDER"、"MEMBER"
//   - operationType: 操作类型，如："CREATE"、"UPDATE"、"DELETE"
//   - content: 操作内容，详细描述这次操作（建议使用JSON格式）
//   - clientIP: 客户端IP地址，从请求中获取
//   - relatedID: 关联ID（可选），如：房间ID、订单ID等
//   - isSuccess: 操作是否成功，true表示成功，false表示失败
//
// 使用示例：
//
//	// 记录房间创建操作
//	content := fmt.Sprintf(`{"action":"CREATE","room_no":"%s","price":%f}`, roomNo, price)
//	util.LogOperation(ctx, userID, "ROOM", "CREATE", content, c.ClientIP(), &roomID, true)
//
//	// 记录操作失败
//	util.LogOperation(ctx, userID, "ORDER", "UPDATE", "修改订单失败：权限不足", c.ClientIP(), &orderID, false)
func LogOperation(ctx context.Context, operatorID uint64, module, operationType, content, clientIP string, relatedID *uint64, isSuccess bool) {
	// 使用goroutine异步记录日志
	// 原因：日志记录不是核心业务逻辑，不应阻塞主流程
	// 即使日志记录失败，也不影响业务操作的成功
	go func() {
		// 构建操作日志请求
		req := &hotel.CreateOperationLogReq{
			OperatorId:    int64(operatorID),     // 操作人ID
			Module:        module,                // 操作模块
			OperationType: operationType,         // 操作类型
			Content:       content,               // 操作内容
			OperationIp:   getClientIP(clientIP), // 客户端IP（格式化后）
			IsSuccess:     isSuccess,             // 是否成功
		}

		// 设置关联ID（可选）
		if relatedID != nil {
			id := int64(*relatedID)
			req.RelatedId = &id
		}

		// 调用RPC接口记录日志
		// 使用 _ 忽略返回值和错误
		// 原因：日志记录是辅助功能，失败时不需要处理
		_, _ = client.HotelClient.CreateOperationLog(ctx, req)
	}()
}

// getClientIP 从IP字符串中提取客户端真实IP
//
// 功能：
//   - 处理单个IP地址
//   - 处理多个IP地址（通过代理或负载均衡时）
//   - 处理空字符串
//
// 参数：
//   - ip: 原始IP字符串，可能包含多个IP（逗号分隔）
//
// 返回：
//   - string: 格式化后的客户端IP
//
// 处理逻辑：
//   - 如果IP为空，返回"unknown"
//   - 如果IP包含多个地址（逗号分隔），取第一个
//   - 自动去除首尾空格
//
// 示例：
//   - "192.168.1.100" → "192.168.1.100"
//   - "192.168.1.100, 10.0.0.1" → "192.168.1.100"（通过代理）
//   - "" → "unknown"
func getClientIP(ip string) string {
	// 处理空IP
	if ip == "" {
		return "unknown"
	}

	// 处理多个IP的情况
	// 通过代理或负载均衡时，X-Forwarded-For可能包含多个IP
	// 格式：客户端IP, 代理1 IP, 代理2 IP
	// 取第一个IP即为客户端真实IP
	ips := strings.Split(ip, ",")
	if len(ips) > 0 {
		// 去除首尾空格并返回
		return strings.TrimSpace(ips[0])
	}

	return ip
}
