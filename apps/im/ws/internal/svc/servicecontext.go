// =============================================================================
// 服务上下文 - 依赖注入容器
// =============================================================================
// 提供服务运行所需的所有依赖组件，包括:
//   - 配置信息
//   - 数据库模型（MongoDB）
//   - 消息队列客户端（Kafka）
//
// 依赖组件:
//   - ChatLogModel: 聊天记录数据模型，用于存储和查询聊天消息
//   - MsgChatTransferClient: 聊天消息传输客户端，将消息推送到 Kafka
//   - MsgReadTransferClient: 已读消息传输客户端，将已读状态推送到 Kafka
//
// 使用场景:
//   - 在服务启动时初始化一次
//   - 在各个 Handler 中通过依赖注入使用
//   - 确保所有组件共享同一个实例（单例模式）
//
// =============================================================================

package svc

import (
	"imooc.com/easy-chat/apps/im/immodels"
	"imooc.com/easy-chat/apps/im/ws/internal/config"
	"imooc.com/easy-chat/apps/task/mq/mqclient"
)

// ServiceContext 服务上下文结构体
// 作为依赖注入容器，持有服务运行所需的所有依赖组件
type ServiceContext struct {
	Config config.Config // 服务配置信息

	immodels.ChatLogModel              // 聊天记录数据模型，用于 MongoDB 操作
	mqclient.MsgChatTransferClient     // 聊天消息传输客户端，用于发送消息到 Kafka
	mqclient.MsgReadTransferClient     // 已读消息传输客户端，用于发送已读状态到 Kafka
}

// NewServiceContext 创建服务上下文实例
// 初始化所有依赖组件，包括数据库连接和消息队列客户端
//
// 初始化流程:
//  1. 创建 Kafka 聊天消息传输客户端（连接到配置的 Broker 和 Topic）
//  2. 创建 Kafka 已读消息传输客户端
//  3. 创建 MongoDB 聊天记录模型（连接到配置的数据库）
//
// 参数:
//   - c: 服务配置对象，包含数据库、消息队列等配置信息
//
// 返回:
//   - *ServiceContext: 初始化完成的服务上下文实例
func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		// 初始化聊天消息传输客户端，用于将用户发送的消息推送到 Kafka
		MsgChatTransferClient: mqclient.NewMsgChatTransferClient(c.MsgChatTransfer.Addrs, c.MsgChatTransfer.Topic),
		// 初始化已读消息传输客户端，用于将消息已读状态推送到 Kafka
		MsgReadTransferClient: mqclient.NewMsgReadTransferClient(c.MsgReadTransfer.Addrs, c.MsgReadTransfer.Topic),
		// 初始化聊天记录模型，连接到 MongoDB 数据库
		ChatLogModel:          immodels.MustChatLogModel(c.Mongo.Url, c.Mongo.Db),
	}
}
