// =============================================================================
// IM Models - 用户会话列表模型层
// =============================================================================
// 提供用户会话列表的数据访问接口和实现
//
// 接口:
//   - ConversationsModel: 用户会话列表模型接口，可自定义扩展
//
// 实现:
//   - customConversationsModel: 自定义用户会话列表模型实现
//
// 功能:
//   - 创建用户会话列表模型实例
//   - 连接 MongoDB 数据库
//   - 提供便捷的模型创建方法
//
// =============================================================================
package immodels

import "github.com/zeromicro/go-zero/core/stores/mon"

// 确保 customConversationsModel 实现了 ConversationsModel 接口
var _ ConversationsModel = (*customConversationsModel)(nil)

type (
	// ConversationsModel 用户会话列表模型接口
	// 可在此接口中添加自定义方法，并在 customConversationsModel 中实现
	// 继承了 conversationsModel 接口的所有基础 CRUD 方法
	ConversationsModel interface {
		conversationsModel
	}

	// customConversationsModel 自定义用户会话列表模型
	// 嵌入 defaultConversationsModel，可以重写或扩展其方法
	customConversationsModel struct {
		*defaultConversationsModel
	}
)

// NewConversationsModel 创建用户会话列表模型实例
// 连接指定的 MongoDB 数据库和集合，返回用户会话列表模型接口
//
// 参数:
//   - url: MongoDB 连接字符串（如: mongodb://localhost:27017）
//   - db: 数据库名称
//   - collection: 集合名称
//
// 返回:
//   - ConversationsModel: 用户会话列表模型接口实例
func NewConversationsModel(url, db, collection string) ConversationsModel {
	conn := mon.MustNewModel(url, db, collection)
	return &customConversationsModel{
		defaultConversationsModel: newDefaultConversationsModel(conn),
	}
}

// MustConversationsModel 使用默认集合名创建用户会话列表模型
// 自动使用 "conversations" 作为集合名称
//
// 参数:
//   - url: MongoDB 连接字符串
//   - db: 数据库名称
//
// 返回:
//   - ConversationsModel: 用户会话列表模型接口实例
func MustConversationsModel(url, db string) ConversationsModel {
	return NewConversationsModel(url, db, "conversations")
}
