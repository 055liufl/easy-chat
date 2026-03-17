// =============================================================================
// Kafka 消息数据结构定义
// =============================================================================
// 本模块定义 Kafka 消息队列中传输的数据结构，包括：
//   - MsgChatTransfer: 聊天消息传输数据
//   - MsgMarkRead: 已读消息标记数据
//
// 业务场景:
//   - IM 服务将消息发送到 Kafka
//   - Task 服务从 Kafka 消费消息
//   - 通过 JSON 序列化/反序列化传输
//
// =============================================================================

package mq

import "imooc.com/easy-chat/pkg/constants"

// MsgChatTransfer 聊天消息传输数据
// 用于 Kafka Topic: MsgChatTransfer
type MsgChatTransfer struct {
	ConversationId     string `json:"conversationId"`     // 会话 ID
	constants.ChatType `json:"chatType"`               // 聊天类型：1-单聊，2-群聊
	SendId             string   `json:"sendId"`          // 发送者用户 ID
	RecvId             string   `json:"recvId"`          // 接收者 ID（单聊为用户 ID，群聊为群 ID）
	RecvIds            []string `json:"recvIds"`         // 接收者 ID 列表（群聊时使用）
	SendTime           int64    `json:"sendTime"`        // 发送时间（Unix 时间戳，毫秒）

	constants.MType `json:"mType"`   // 消息类型：1-文本，2-图片，3-语音，4-视频等
	Content         string `json:"content"` // 消息内容（文本或 JSON 格式）
}

// MsgMarkRead 已读消息标记数据
// 用于 Kafka Topic: MsgReadTransfer
type MsgMarkRead struct {
	constants.ChatType `json:"chatType"`       // 聊天类型：1-单聊，2-群聊
	ConversationId     string   `json:"conversationId"` // 会话 ID
	SendId             string   `json:"sendId"`         // 标记已读的用户 ID
	RecvId             string   `json:"recvId"`         // 接收者 ID（单聊为对方用户 ID，群聊为群 ID）
	MsgIds             []string `json:"msgIds"`         // 已读消息 ID 列表
}
