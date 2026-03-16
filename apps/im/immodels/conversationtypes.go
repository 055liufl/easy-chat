// =============================================================================
// IM Models - 会话数据结构定义
// =============================================================================
// 定义单个会话的数据结构
//
// 数据结构:
//   - Conversation: 会话实体，存储会话的基本信息和最新消息
//
// 用途:
//   - 用于会话列表展示
//   - 存储会话的统计信息（消息总数、序列号等）
//   - 关联最新的一条聊天消息
//
// =============================================================================
package immodels

import (
	"imooc.com/easy-chat/pkg/constants"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Conversation 会话实体
// 表示一个聊天会话的完整信息，包括会话类型、消息统计、最新消息等
type Conversation struct {
	// ID MongoDB 文档唯一标识符
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	// ConversationId 会话唯一标识符
	// 格式通常为: 单聊 "user1_user2"，群聊 "group_groupId"
	ConversationId string `bson:"conversationId,omitempty"`

	// ChatType 聊天类型（单聊、群聊等）
	ChatType constants.ChatType `bson:"chatType,omitempty"`

	// IsShow 是否在会话列表中显示
	// true: 显示在用户的会话列表中
	// false: 隐藏（用户删除会话但不删除消息记录）
	IsShow bool `bson:"isShow,omitempty"`

	// Total 会话中的消息总数
	Total int `bson:"total,omitempty"`

	// Seq 会话序列号，用于消息排序和同步
	Seq int64 `bson:"seq"`

	// Msg 会话中的最新一条消息
	// 用于在会话列表中显示最后一条消息的预览
	Msg *ChatLog `bson:"msg,omitempty"`

	// UpdateAt 会话更新时间
	UpdateAt time.Time `bson:"updateAt,omitempty" json:"updateAt,omitempty"`

	// CreateAt 会话创建时间
	CreateAt time.Time `bson:"createAt,omitempty" json:"createAt,omitempty"`
}
