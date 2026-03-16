// =============================================================================
// IM Models - 聊天记录模型层
// =============================================================================
// 提供聊天记录的数据访问接口和实现
//
// 接口:
//   - ChatLogModel: 聊天记录模型接口，可自定义扩展
//
// 实现:
//   - customChatLogModel: 自定义聊天记录模型实现
//
// 功能:
//   - 创建聊天记录模型实例
//   - 连接 MongoDB 数据库
//   - 提供便捷的模型创建方法
//
// =============================================================================
package immodels

import "github.com/zeromicro/go-zero/core/stores/mon"

// 确保 customChatLogModel 实现了 ChatLogModel 接口
var _ ChatLogModel = (*customChatLogModel)(nil)

type (
	// ChatLogModel 聊天记录模型接口
	// 可在此接口中添加自定义方法，并在 customChatLogModel 中实现
	// 继承了 chatLogModel 接口的所有基础 CRUD 方法
	ChatLogModel interface {
		chatLogModel
	}

	// customChatLogModel 自定义聊天记录模型
	// 嵌入 defaultChatLogModel，可以重写或扩展其方法
	customChatLogModel struct {
		*defaultChatLogModel
	}
)

// NewChatLogModel 创建聊天记录模型实例
// 连接指定的 MongoDB 数据库和集合，返回聊天记录模型接口
//
// 参数:
//   - url: MongoDB 连接字符串（如: mongodb://localhost:27017）
//   - db: 数据库名称
//   - collection: 集合名称
//
// 返回:
//   - ChatLogModel: 聊天记录模型接口实例
func NewChatLogModel(url, db, collection string) ChatLogModel {
	conn := mon.MustNewModel(url, db, collection)
	return &customChatLogModel{
		defaultChatLogModel: newDefaultChatLogModel(conn),
	}
}

// MustChatLogModel 使用默认集合名创建聊天记录模型
// 自动使用 "chat_log" 作为集合名称
//
// 参数:
//   - url: MongoDB 连接字符串
//   - db: 数据库名称
//
// 返回:
//   - ChatLogModel: 聊天记录模型接口实例
func MustChatLogModel(url, db string) ChatLogModel {
	return NewChatLogModel(url, db, "chat_log")
}
