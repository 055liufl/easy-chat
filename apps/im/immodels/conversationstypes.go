// =============================================================================
// IM Models - 用户会话列表数据结构定义
// =============================================================================
// 定义用户的会话列表数据结构
//
// 数据结构:
//   - Conversations: 用户会话列表实体，存储用户的所有会话
//
// 用途:
//   - 管理单个用户的所有会话
//   - 使用 map 结构快速查找和更新特定会话
//   - 一个用户对应一个 Conversations 文档
//
// =============================================================================
package immodels

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Conversations 用户会话列表实体
// 存储单个用户的所有会话信息，使用 map 结构便于快速访问和更新
type Conversations struct {
	// ID MongoDB 文档唯一标识符
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	// UserId 用户 ID
	// 标识该会话列表属于哪个用户
	UserId string `bson:"userId"`

	// ConversationList 会话列表
	// key: conversationId（会话唯一标识符）
	// value: Conversation（会话详细信息）
	// 使用 map 结构可以快速通过 conversationId 查找和更新会话
	ConversationList map[string]*Conversation `bson:"conversationList"`

	// UpdateAt 记录更新时间
	UpdateAt time.Time `bson:"updateAt,omitempty" json:"updateAt,omitempty"`

	// CreateAt 记录创建时间
	CreateAt time.Time `bson:"createAt,omitempty" json:"createAt,omitempty"`
}
