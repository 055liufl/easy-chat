// =============================================================================
// IM Models - 聊天记录数据结构定义
// =============================================================================
// 定义聊天记录的数据结构和常量
//
// 数据结构:
//   - ChatLog: 聊天记录实体，存储单条聊天消息的完整信息
//
// 常量:
//   - DefaultChatLogLimit: 默认查询聊天记录的数量限制
//
// =============================================================================
package immodels

import (
	"imooc.com/easy-chat/pkg/constants"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DefaultChatLogLimit 默认聊天记录查询限制
// 用于限制单次查询返回的聊天记录数量，防止一次性加载过多数据
var DefaultChatLogLimit int64 = 100

// ChatLog 聊天记录实体
// 存储单条聊天消息的完整信息，包括发送者、接收者、消息内容、时间等
type ChatLog struct {
	// ID MongoDB 文档唯一标识符
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	// ConversationId 会话 ID，标识消息所属的会话
	ConversationId string `bson:"conversationId"`

	// SendId 发送者用户 ID
	SendId string `bson:"sendId"`

	// RecvId 接收者用户 ID（单聊时使用）
	RecvId string `bson:"recvId"`

	// MsgFrom 消息来源标识
	MsgFrom int `bson:"msgFrom"`

	// ChatType 聊天类型（单聊、群聊等）
	ChatType constants.ChatType `bson:"chatType"`

	// MsgType 消息类型（文本、图片、语音等）
	MsgType constants.MType `bson:"msgType"`

	// MsgContent 消息内容
	MsgContent string `bson:"msgContent"`

	// SendTime 消息发送时间戳（毫秒）
	SendTime int64 `bson:"sendTime"`

	// Status 消息状态（已发送、已送达、已读等）
	Status int `bson:"status"`

	// ReadRecords 已读记录，存储已读用户信息的序列化数据
	ReadRecords []byte `bson:"readRecords"`

	// UpdateAt 记录更新时间
	UpdateAt time.Time `bson:"updateAt,omitempty" json:"updateAt,omitempty"`

	// CreateAt 记录创建时间
	CreateAt time.Time `bson:"createAt,omitempty" json:"createAt,omitempty"`
}
