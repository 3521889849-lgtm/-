// Package ws 提供 WebSocket 实时通信功能
// 用于客服与用户之间的实时消息传递
// 主要功能：
// 1. 管理 WebSocket 连接的注册和注销
// 2. 广播和单播消息
// 3. 统计在线客服人数
package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"example_shop/gateway/rpc"                         // RPC客户端，用于消息持久化
	"example_shop/pkg/logger"                          // 日志工具
	"example_shop/service/customer/kitex_gen/customer" // RPC接口定义

	"go.uber.org/zap"
)

// Hub 维护活跃的WebSocket连接并处理消息路由
// 功能说明：
// 1. 管理客户端连接的注册和注销
// 2. 支持同账号互踢（新连接替换旧连接）
// 3. 广播消息给所有客户端
// 4. 单播消息给指定用户
type Hub struct {
	// clients 注册的客户端映射 map[UserID]*Client
	// 当前实现同账号互踢，新连接会替换旧连接
	clients map[int64]*Client

	// staffs 在线客服列表（仅存储UserID）
	staffs map[int64]bool

	// broadcast 广播消息通道（发送给所有连接的客户端）
	broadcast chan []byte

	// register 客户端注册请求通道
	register chan *Client

	// unregister 客户端注销请求通道
	unregister chan *Client

	// rpcClient RPC客户端，用于将消息持久化到数据库
	rpcClient *rpc.CustomerClient

	// mu 读写锁，保护客户端映射的并发访问
	mu sync.RWMutex
}

// NewHub 创建新的Hub实例
// 参数:
//   - rpcClient: RPC客户端，用于消息持久化
//
// 返回: Hub实例指针
func NewHub(rpcClient *rpc.CustomerClient) *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[int64]*Client),
		staffs:     make(map[int64]bool),
		rpcClient:  rpcClient,
	}
}

// Run 启动Hub的主循环
// 处理客户端注册、注销和广播消息
// 应在单独的goroutine中运行
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			// 处理客户端注册
			h.mu.Lock()
			// 如果已有连接，关闭旧连接（同账号互踢）
			if old, ok := h.clients[client.UserID]; ok {
				// 使用 closeOnce 确保只关闭一次，防止 panic
				old.closeOnce.Do(func() {
					close(old.send)
				})
				delete(h.clients, client.UserID)
			}
			// 注册新连接
			h.clients[client.UserID] = client
			// 如果是客服或管理员，加入在线客服列表
			if client.Role == "customer_service" || client.Role == "admin" {
				h.staffs[client.UserID] = true
			}
			h.mu.Unlock()
			logger.Info("Client registered", zap.Int64("user_id", client.UserID), zap.String("role", client.Role))

		case client := <-h.unregister:
			// 处理客户端注销
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				// 使用 closeOnce 确保只关闭一次，防止 panic
				client.closeOnce.Do(func() {
					close(client.send)
				})
				delete(h.staffs, client.UserID)
			}
			h.mu.Unlock()
			logger.Info("Client unregistered", zap.Int64("user_id", client.UserID))

		case message := <-h.broadcast:
			// 广播消息给所有连接的客户端
			h.mu.Lock() // 改为写锁，因为可能需要 delete
			for _, client := range h.clients {
				select {
				case client.send <- message:
				default:
					// 缓冲区已满，关闭连接
					// 使用 closeOnce 确保只关闭一次，防止 panic
					client.closeOnce.Do(func() {
						close(client.send)
					})
					delete(h.clients, client.UserID)
				}
			}
			h.mu.Unlock()
		}
	}
}

// UnicastRaw 发送原始字节消息给指定用户
// 参数:
//   - targetID: 目标用户ID
//   - message: 要发送的字节消息
func (h *Hub) UnicastRaw(targetID int64, message []byte) {
	h.mu.RLock()
	client, ok := h.clients[targetID]
	h.mu.RUnlock()

	if ok {
		select {
		case client.send <- message:
			// 发送成功
		default:
			// 发送缓冲区已满，记录警告并注销该客户端
			logger.Warn("Failed to send message to client (buffer full)", zap.Int64("user_id", targetID))
			h.unregister <- client
		}
	}
}

// UnicastJSON 发送JSON格式消息给指定用户
// 自动将对象序列化为JSON后发送
// 参数:
//   - targetID: 目标用户ID
//   - v: 要发送的对象
func (h *Hub) UnicastJSON(targetID int64, v interface{}) {
	msg, err := json.Marshal(v)
	if err != nil {
		logger.Error("JSON marshal error", zap.Error(err))
		return
	}
	h.UnicastRaw(targetID, msg)
}

// HandleMessage 处理从客户端接收到的业务消息
// 当前支持的消息类型：
// - "chat": 聊天消息，会持久化到数据库并转发给目标用户
// 参数:
//   - client: 发送消息的客户端
//   - msg: 原始消息字节
func (h *Hub) HandleMessage(client *Client, msg []byte) {
	// 解析消息结构
	var input struct {
		Type    string          `json:"type"`    // 消息类型
		Payload json.RawMessage `json:"payload"` // 消息负载
	}
	if err := json.Unmarshal(msg, &input); err != nil {
		logger.Warn("Invalid message format", zap.Error(err))
		return
	}

	// 处理聊天消息
	if input.Type == "chat" {
		var chatMsg struct {
			ConversationID string `json:"conversation_id"` // 会话i
			Content        string `json:"content"`         // 消息内容
			MsgType        int32  `json:"msg_type"`        // 消息类型
			ToUserID       int64  `json:"to_user_id"`      // 目标用户ID
		}
		if err := json.Unmarshal(input.Payload, &chatMsg); err != nil {
			logger.Warn("Invalid chat payload", zap.Error(err))
			return
		}

		// 构造上下文，包含 TraceID 用于链路追踪
		ctx := context.Background()
		traceID := logger.NewTraceID()
		ctx = logger.WithTraceID(ctx, traceID)

		logger.InfoWithTrace(ctx, "Received chat message",
			zap.Int64("from", client.UserID),
			zap.Int64("to", chatMsg.ToUserID),
			zap.String("content", chatMsg.Content))

		// 确定发送者类型：1=用户, 2=客服
		senderType := int8(1) // 默认为用户
		if client.Role == "customer_service" || client.Role == "admin" {
			senderType = 2 // 客服
		}

		// 调用 RPC 将消息持久化到数据库
		req := &customer.SendConversationMessageReq{
			ConvId:     chatMsg.ConversationID,
			SenderType: senderType,
			SenderId:   fmt.Sprintf("%d", client.UserID),
			MsgContent: chatMsg.Content,
		}

		_, err := h.rpcClient.SendConversationMessage(ctx, req)
		if err != nil {
			logger.ErrorWithTrace(ctx, "Failed to save message", zap.Error(err))
			return
		}

		// 构造响应消息
		response := map[string]interface{}{
			"type": "chat",
			"payload": map[string]interface{}{
				"conversation_id": chatMsg.ConversationID,
				"from_user_id":    client.UserID,
				"content":         chatMsg.Content,
				"msg_type":        chatMsg.MsgType,
				"create_time":     time.Now().Format(time.RFC3339),
			},
		}

		// 发送ACK给发送者
		h.UnicastJSON(client.UserID, response)

		// 转发消息给目标用户
		if chatMsg.ToUserID != 0 {
			h.UnicastJSON(chatMsg.ToUserID, response)
		}
	}
}

// GetOnlineStats 获取在线统计信息
// 返回当前在线客服数量、客服ID列表和总连接数
// 返回: 包含统计信息的map
func (h *Hub) GetOnlineStats() map[string]interface{} {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// 收集在线客服ID列表
	staffList := make([]int64, 0, len(h.staffs))
	for uid := range h.staffs {
		staffList = append(staffList, uid)
	}

	return map[string]interface{}{
		"online_staff_count": len(h.staffs),   // 在线客服数量
		"online_staff_ids":   staffList,       // 在线客服ID列表
		"total_connections":  len(h.clients),  // 总连接数
	}
}
