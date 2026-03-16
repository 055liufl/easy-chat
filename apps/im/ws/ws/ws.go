// =============================================================================
// WebSocket 消息数据结构 - IM 消息协议定义
// =============================================================================
// 定义 IM 系统中 WebSocket 通信的所有消息数据结构，包括:
//   - 消息基础结构（Msg）
//   - 聊天消息结构（Chat）
//   - 推送消息结构（Push）
//   - 已读标记结构（MarkRead）
//
// 消息类型:
//   - 文本消息、图片消息、语音消息、视频消息等
//
// 聊天类型:
//   - 单聊（SingleChatType）
//   - 群聊（GroupChatType）
//
// 使用场景:
//   - 客户端与服务端之间的消息传输
//   - 服务端内部消息流转
//   - Kafka 消息队列数据格式
//
// =============================================================================

package ws

import "imooc.com/easy-chat/pkg/constants"

type (
	// Msg 消息基础结构
	// 包含消息的核心信息，如消息 ID、内容、类型、已读记录等
	Msg struct {
		MsgId           string            `mapstructure:"msgId"`       // 消息唯一标识，用于去重和已读标记
		ReadRecords     map[string]string `mapstructure:"readRecords"` // 已读记录，key: 用户ID, value: 已读时间
		constants.MType `mapstructure:"mType"`                         // 消息类型: 1-文本, 2-图片, 3-语音, 4-视频等
		Content         string            `mapstructure:"content"`     // 消息内容（文本内容或媒体文件 URL）
	}

	// Chat 聊天消息结构
	// 用于客户端发送聊天消息和服务端推送聊天消息
	Chat struct {
		ConversationId     string `mapstructure:"conversationId"` // 会话 ID，单聊为两个用户 ID 组合，群聊为群组 ID
		constants.ChatType `mapstructure:"chatType"`           // 聊天类型: 1-单聊, 2-群聊
		SendId             string `mapstructure:"sendId"`         // 发送者用户 ID
		RecvId             string `mapstructure:"recvId"`         // 接收者 ID（单聊为用户 ID，群聊为群组 ID）
		SendTime           int64  `mapstructure:"sendTime"`       // 发送时间（纳秒时间戳）
		Msg                `mapstructure:"msg"`                 // 消息内容（嵌入 Msg 结构）
	}

	// Push 推送消息结构
	// 用于服务端向客户端推送消息（从 Kafka 消费后推送）
	Push struct {
		ConversationId     string `mapstructure:"conversationId"` // 会话 ID
		constants.ChatType `mapstructure:"chatType"`           // 聊天类型: 1-单聊, 2-群聊
		SendId             string   `mapstructure:"sendId"`       // 发送者用户 ID
		RecvId             string   `mapstructure:"recvId"`       // 接收者 ID（单聊使用）
		RecvIds            []string `mapstructure:"recvIds"`      // 接收者 ID 列表（群聊使用）
		SendTime           int64    `mapstructure:"sendTime"`     // 发送时间（纳秒时间戳）

		MsgId       string                `mapstructure:"msgId"`       // 消息唯一标识
		ReadRecords map[string]string     `mapstructure:"readRecords"` // 已读记录
		ContentType constants.ContentType `mapstructure:"contentType"` // 内容类型（用于扩展）

		constants.MType `mapstructure:"mType"`   // 消息类型
		Content         string `mapstructure:"content"` // 消息内容
	}

	// MarkRead 消息已读标记结构
	// 用于客户端标记消息为已读
	MarkRead struct {
		constants.ChatType `mapstructure:"chatType"`       // 聊天类型: 1-单聊, 2-群聊
		RecvId             string   `mapstructure:"recvId"`         // 接收者 ID
		ConversationId     string   `mapstructure:"conversationId"` // 会话 ID
		MsgIds             []string `mapstructure:"msgIds"`         // 需要标记为已读的消息 ID 列表
	}
)
