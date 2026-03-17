// =============================================================================
// 聊天消息传输处理器 - Kafka 聊天消息消费者
// =============================================================================
// 本模块负责消费 Kafka 中的聊天消息，完成以下任务：
//   - 消息持久化到 MongoDB
//   - 更新会话最新消息
//   - 初始化消息已读记录
//   - 推送消息给接收方
//
// 数据来源:
//   - Kafka Topic: MsgChatTransfer（聊天消息队列）
//
// 业务场景:
//   - 用户发送消息后，IM 服务将消息发送到 Kafka
//   - Task 服务消费消息，持久化到 MongoDB
//   - 持久化成功后，通过 WebSocket 推送给接收方
//   - 单聊消息：直接推送给接收方
//   - 群聊消息：查询群成员后批量推送
//
// 处理流程:
//   1. 从 Kafka 消费消息
//   2. 生成消息 ID（MongoDB ObjectID）
//   3. 持久化消息到 ChatLog 集合
//   4. 更新会话（Conversation）最新消息
//   5. 推送消息给接收方
//
// =============================================================================

package msgTransfer

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"imooc.com/easy-chat/apps/im/immodels"
	"imooc.com/easy-chat/apps/im/ws/ws"
	"imooc.com/easy-chat/apps/task/mq/internal/svc"
	"imooc.com/easy-chat/apps/task/mq/mq"
	"imooc.com/easy-chat/pkg/bitmap"
)

// MsgChatTransfer 聊天消息传输处理器
// 实现 Kafka 消费者接口，处理聊天消息的持久化和推送
type MsgChatTransfer struct {
	*baseMsgTransfer // 继承基础消息传输功能（推送逻辑）
}

// NewMsgChatTransfer 创建聊天消息传输处理器实例
//
// 参数:
//   - svc: 服务上下文，包含 MongoDB 模型、WebSocket 客户端等
//
// 返回:
//   - *MsgChatTransfer: 聊天消息传输处理器实例
func NewMsgChatTransfer(svc *svc.ServiceContext) *MsgChatTransfer {
	return &MsgChatTransfer{
		NewBaseMsgTransfer(svc),
	}
}

// Consume Kafka 消费者接口实现
// 消费聊天消息，完成持久化和推送
//
// 业务逻辑:
//   1. 解析 Kafka 消息（JSON 格式）
//   2. 生成消息 ID（MongoDB ObjectID）
//   3. 调用 addChatLog 持久化消息
//   4. 调用 Transfer 推送消息给接收方
//
// 参数:
//   - key: Kafka 消息 key（通常为用户 ID 或会话 ID，用于分区）
//   - value: Kafka 消息 value（JSON 格式的聊天消息）
//
// 返回:
//   - error: 处理失败时返回错误，Kafka 会重试
func (m *MsgChatTransfer) Consume(key, value string) error {
	fmt.Println("key : ", key, " value : ", value)

	var (
		data  mq.MsgChatTransfer
		ctx   = context.Background()
		msgId = primitive.NewObjectID() // 生成消息 ID
	)
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return err
	}

	// 持久化消息到 MongoDB
	if err := m.addChatLog(ctx, msgId, &data); err != nil {
		return err
	}

	// 推送消息给接收方
	return m.Transfer(ctx, &ws.Push{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		RecvIds:        data.RecvIds,
		SendTime:       data.SendTime,
		MType:          data.MType,
		MsgId:          msgId.Hex(), // 返回生成的消息 ID
		Content:        data.Content,
	})
}

// addChatLog 持久化聊天消息到 MongoDB
// 完成消息记录插入和会话更新
//
// 业务逻辑:
//   1. 构建 ChatLog 文档（消息记录）
//   2. 初始化已读记录（发送者默认已读）
//   3. 插入到 ChatLog 集合
//   4. 更新 Conversation 集合的最新消息
//
// 参数:
//   - ctx: 上下文
//   - msgId: 消息 ID（MongoDB ObjectID）
//   - data: 聊天消息数据
//
// 返回:
//   - error: 持久化失败时返回错误
func (m *MsgChatTransfer) addChatLog(ctx context.Context, msgId primitive.ObjectID, data *mq.MsgChatTransfer) error {
	// 构建聊天记录文档
	chatLog := immodels.ChatLog{
		ID:             msgId,
		ConversationId: data.ConversationId,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		ChatType:       data.ChatType,
		MsgFrom:        0,                // 消息来源：0-用户消息
		MsgType:        data.MType,       // 消息类型：文本、图片、语音等
		MsgContent:     data.Content,     // 消息内容
		SendTime:       data.SendTime,    // 发送时间
	}

	// 初始化已读记录（使用 Bitmap 存储已读用户 ID）
	readRecords := bitmap.NewBitmap(0)
	readRecords.Set(chatLog.SendId) // 发送者默认已读
	chatLog.ReadRecords = readRecords.Export()

	// 插入消息记录到 ChatLog 集合
	err := m.svcCtx.ChatLogModel.Insert(ctx, &chatLog)
	if err != nil {
		return err
	}

	// 更新会话的最新消息（用于会话列表显示）
	return m.svcCtx.ConversationModel.UpdateMsg(ctx, &chatLog)
}
