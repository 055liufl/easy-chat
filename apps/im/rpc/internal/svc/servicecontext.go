// =============================================================================
// IM RPC 服务上下文
// =============================================================================
// 管理 IM RPC 服务的依赖资源，包括:
//   - 配置信息
//   - 数据库模型（聊天记录、会话列表、会话信息）
//
// 职责:
//   - 初始化所有依赖的数据模型
//   - 提供统一的服务上下文给各个业务逻辑层
//   - 管理数据库连接的生命周期
//
// =============================================================================
package svc

import (
	"imooc.com/easy-chat/apps/im/immodels"
	"imooc.com/easy-chat/apps/im/rpc/internal/config"
)

// ServiceContext IM RPC 服务上下文
// 包含服务运行所需的所有依赖资源
type ServiceContext struct {
	Config config.Config // 服务配置

	// 数据模型
	immodels.ChatLogModel       // 聊天记录模型，用于存储和查询聊天消息
	immodels.ConversationsModel // 用户会话列表模型，存储用户的所有会话
	immodels.ConversationModel  // 会话信息模型，存储会话的基本信息
}

// NewServiceContext 创建服务上下文实例
// 根据配置初始化所有依赖的数据模型
//
// 参数:
//   - c: 服务配置，包含 MongoDB 连接信息等
//
// 返回:
//   - *ServiceContext: 初始化完成的服务上下文实例
//
// 说明:
//   - 使用 Must 系列方法初始化模型，连接失败会 panic
//   - 所有模型共享同一个 MongoDB 连接配置
func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		ChatLogModel:       immodels.MustChatLogModel(c.Mongo.Url, c.Mongo.Db),
		ConversationsModel: immodels.MustConversationsModel(c.Mongo.Url, c.Mongo.Db),
		ConversationModel:  immodels.MustConversationModel(c.Mongo.Url, c.Mongo.Db),
	}
}
