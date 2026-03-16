// =============================================================================
// IM Models - 会话模型层
// =============================================================================
// 提供会话的数据访问接口和实现
//
// 接口:
//   - ConversationModel: 会话模型接口，可自定义扩展
//
// 实现:
//   - customConversationModel: 自定义会话模型实现
//
// 功能:
//   - 创建会话模型实例
//   - 连接 MongoDB 数据库
//   - 提供便捷的模型创建方法
//
// =============================================================================
package immodels

import "github.com/zeromicro/go-zero/core/stores/mon"

// 确保 customConversationModel 实现了 ConversationModel 接口
var _ ConversationModel = (*customConversationModel)(nil)

type (
	// ConversationModel 会话模型接口
	// 可在此接口中添加自定义方法，并在 customConversationModel 中实现
	// 继承了 conversationModel 接口的所有基础 CRUD 方法
	ConversationModel interface {
		conversationModel
	}

	// customConversationModel 自定义会话模型
	// 嵌入 defaultConversationModel，可以重写或扩展其方法
	customConversationModel struct {
		*defaultConversationModel
	}
)

// NewConversationModel 创建会话模型实例
// 连接指定的 MongoDB 数据库和集合，返回会话模型接口
//
// 参数:
//   - url: MongoDB 连接字符串（如: mongodb://localhost:27017）
//   - db: 数据库名称
//   - collection: 集合名称
//
// 返回:
//   - ConversationModel: 会话模型接口实例
func NewConversationModel(url, db, collection string) ConversationModel {
	conn := mon.MustNewModel(url, db, collection)
	return &customConversationModel{
		defaultConversationModel: newDefaultConversationModel(conn),
	}
}

// MustConversationModel 使用默认集合名创建会话模型
// 自动使用 "conversation" 作为集合名称
//
// 参数:
//   - url: MongoDB 连接字符串
//   - db: 数据库名称
//
// 返回:
//   - ConversationModel: 会话模型接口实例
func MustConversationModel(url, db string) ConversationModel {
	return NewConversationModel(url, db, "conversation")
}
