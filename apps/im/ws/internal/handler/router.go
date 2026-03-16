// =============================================================================
// 路由注册器 - WebSocket 消息路由配置
// =============================================================================
// 负责注册所有 WebSocket 消息处理路由，包括:
//   - 用户在线状态管理
//   - 会话聊天消息处理
//   - 消息已读标记处理
//   - 消息推送处理
//
// 路由格式:
//   客户端发送的消息需要包含 method 字段，格式为 "模块.操作"
//   例如: {"method": "user.online", "data": {...}}
//
// 已注册路由:
//   - user.online: 用户上线通知，返回在线用户列表
//   - conversation.chat: 发送聊天消息（单聊/群聊）
//   - conversation.markChat: 标记消息为已读
//   - push: 服务端推送消息到客户端（内部使用）
//
// =============================================================================

package handler

import (
	"imooc.com/easy-chat/apps/im/ws/internal/handler/conversation"
	"imooc.com/easy-chat/apps/im/ws/internal/handler/push"
	"imooc.com/easy-chat/apps/im/ws/internal/handler/user"
	"imooc.com/easy-chat/apps/im/ws/internal/svc"
	"imooc.com/easy-chat/apps/im/ws/websocket"
)

// RegisterHandlers 注册所有 WebSocket 消息处理路由
// 将业务处理函数绑定到对应的消息方法上
//
// 路由说明:
//  - user.online: 处理用户上线请求，返回当前在线用户列表
//  - conversation.chat: 处理聊天消息发送，支持单聊和群聊
//  - conversation.markChat: 处理消息已读标记，更新消息已读状态
//  - push: 处理服务端消息推送，将消息推送给目标用户
//
// 参数:
//   - srv: WebSocket 服务器实例，用于注册路由
//   - svc: 服务上下文，包含业务处理所需的依赖组件
func RegisterHandlers(srv *websocket.Server, svc *svc.ServiceContext) {
	srv.AddRoutes([]websocket.Route{
		{
			Method:  "user.online",           // 用户上线路由
			Handler: user.OnLine(svc),        // 处理用户上线，返回在线用户列表
		},
		{
			Method:  "conversation.chat",     // 聊天消息路由
			Handler: conversation.Chat(svc),  // 处理聊天消息发送（单聊/群聊）
		},
		{
			Method:  "conversation.markChat",    // 消息已读标记路由
			Handler: conversation.MarkRead(svc), // 处理消息已读状态更新
		},
		{
			Method:  "push",                  // 消息推送路由
			Handler: push.Push(svc),          // 处理服务端消息推送
		},
	})
}
