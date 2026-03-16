// =============================================================================
// 会话业务逻辑 - 聊天消息业务处理
// =============================================================================
// 提供会话相关的业务逻辑处理，包括:
//   - 单聊消息处理
//   - 群聊消息处理
//   - 消息持久化存储
//   - 会话 ID 生成
//
// 业务流程:
//   1. 接收聊天消息数据
//   2. 生成或验证会话 ID
//   3. 将消息持久化到 MongoDB
//   4. 返回处理结果
//
// 注意:
//   当前代码中包含 time.Sleep(time.Minute)，可能是测试代码，生产环境应移除
//
// =============================================================================

package logic

import (
	"context"
	"imooc.com/easy-chat/apps/im/immodels"
	"imooc.com/easy-chat/apps/im/ws/internal/svc"
	"imooc.com/easy-chat/apps/im/ws/websocket"
	"imooc.com/easy-chat/apps/im/ws/ws"
	"imooc.com/easy-chat/pkg/wuid"
	"time"
)

// Conversation 会话业务逻辑处理器
// 负责处理会话相关的业务逻辑
type Conversation struct {
	ctx context.Context       // 上下文，用于控制超时和取消
	srv *websocket.Server     // WebSocket 服务器实例
	svc *svc.ServiceContext   // 服务上下文，包含数据库等依赖
}

// NewConversation 创建会话业务逻辑处理器实例
// 初始化会话处理器，注入依赖组件
//
// 参数:
//   - ctx: 上下文对象，用于控制超时和取消
//   - srv: WebSocket 服务器实例
//   - svc: 服务上下文，包含数据库连接等依赖
//
// 返回:
//   - *Conversation: 初始化完成的会话处理器实例
func NewConversation(ctx context.Context, srv *websocket.Server, svc *svc.ServiceContext) *Conversation {
	return &Conversation{
		ctx: ctx,
		srv: srv,
		svc: svc,
	}
}

// SingleChat 处理单聊消息
// 将单聊消息持久化到 MongoDB 数据库
//
// 业务流程:
//  1. 检查并生成会话 ID（如果未提供）
//     - 使用发送者和接收者 ID 组合生成唯一会话 ID
//  2. 构造聊天记录对象
//  3. 将聊天记录插入到 MongoDB
//
// 参数:
//   - data: 聊天消息数据，包含会话 ID、接收者、消息内容等
//   - userId: 发送者用户 ID
//
// 返回:
//   - error: 插入失败时返回错误信息
func (l *Conversation) SingleChat(data *ws.Chat, userId string) error {
	// 如果客户端未提供会话 ID，使用发送者和接收者 ID 组合生成
	if data.ConversationId == "" {
		data.ConversationId = wuid.CombineId(userId, data.RecvId)
	}

	// 注意: 此处有一个 1 分钟的延迟，可能是测试代码，生产环境应移除
	time.Sleep(time.Minute)

	// 构造聊天记录对象
	chatLog := immodels.ChatLog{
		ConversationId: data.ConversationId,
		SendId:         userId,
		RecvId:         data.RecvId,
		ChatType:       data.ChatType,
		MsgFrom:        0,                      // 消息来源: 0-用户发送
		MsgType:        data.MType,
		MsgContent:     data.Content,
		SendTime:       time.Now().UnixNano(), // 发送时间（纳秒时间戳）
	}

	// 将聊天记录插入到 MongoDB
	err := l.svc.ChatLogModel.Insert(l.ctx, &chatLog)

	return err
}
