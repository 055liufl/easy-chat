// =============================================================================
// 配置定义 - IM WebSocket 服务配置结构
// =============================================================================
// 定义 IM WebSocket 服务的所有配置项，包括:
//   - 服务基础配置（服务名、监听地址等）
//   - JWT 认证配置（Token 密钥）
//   - MongoDB 数据库配置（连接地址、数据库名）
//   - Kafka 消息队列配置（聊天消息传输、已读消息传输）
//
// 配置来源:
//   - 本地 YAML 配置文件
//   - ETCD 配置中心（支持动态更新）
//
// 配置示例:
//   Name: im-ws
//   ListenOn: :8080
//   JwtAuth:
//     AccessSecret: your-secret-key
//   Mongo:
//     Url: mongodb://localhost:27017
//     Db: easy-chat
//   MsgChatTransfer:
//     Topic: chat-transfer
//     Addrs: [localhost:9092]
//   MsgReadTransfer:
//     Topic: read-transfer
//     Addrs: [localhost:9092]
//
// =============================================================================

package config

import "github.com/zeromicro/go-zero/core/service"

// Config IM WebSocket 服务配置结构体
// 包含服务运行所需的所有配置项
type Config struct {
	service.ServiceConf // 继承 go-zero 服务基础配置（Name、Log、Telemetry 等）

	ListenOn string // WebSocket 服务监听地址，格式: host:port 或 :port

	// JwtAuth JWT 认证配置
	// 用于验证 WebSocket 连接的 Token 合法性
	JwtAuth struct {
		AccessSecret string // JWT Token 签名密钥，用于验证 Token 签名
	}

	// Mongo MongoDB 数据库配置
	// 用于存储聊天记录、会话信息等持久化数据
	Mongo struct {
		Url string // MongoDB 连接地址，格式: mongodb://host:port
		Db  string // 数据库名称
	}

	// MsgChatTransfer 聊天消息传输队列配置
	// 用于将用户发送的聊天消息推送到 Kafka，由消息处理服务消费
	MsgChatTransfer struct {
		Topic string   // Kafka Topic 名称，用于聊天消息传输
		Addrs []string // Kafka Broker 地址列表，格式: ["host1:port1", "host2:port2"]
	}

	// MsgReadTransfer 已读消息传输队列配置
	// 用于将消息已读状态推送到 Kafka，由消息处理服务更新已读状态
	MsgReadTransfer struct {
		Topic string   // Kafka Topic 名称，用于已读消息传输
		Addrs []string // Kafka Broker 地址列表，格式: ["host1:port1", "host2:port2"]
	}
}
