// =============================================================================
// Task MQ 配置模块 - Kafka 消息队列消费者配置
// =============================================================================
// 本模块定义 Task 服务的配置结构，Task 服务是 Kafka 消息队列的消费者，负责：
//   - 消费聊天消息并持久化到 MongoDB
//   - 消费已读消息并更新已读状态
//   - 通过 WebSocket 推送消息给在线用户
//
// 配置来源:
//   - 本地配置文件（YAML 格式）
//   - ETCD 配置中心（支持动态更新）
//
// 业务场景:
//   - 聊天消息异步持久化：IM 服务将消息发送到 Kafka，Task 服务消费后存储到 MongoDB
//   - 消息推送：消息持久化后，通过 WebSocket 推送给接收方
//   - 已读状态处理：用户标记消息已读后，更新 MongoDB 中的已读记录并推送给发送方
//
// =============================================================================

package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

// Config Task 服务配置结构
type Config struct {
	service.ServiceConf        // 服务基础配置（名称、日志等）
	ListenOn            string // HTTP 监听地址（用于健康检查等）

	// Kafka 消费者配置
	MsgChatTransfer kq.KqConf // 聊天消息传输队列配置（消费聊天消息）
	MsgReadTransfer kq.KqConf // 已读消息传输队列配置（消费已读消息）

	// Redis 配置（用于缓存和分布式锁）
	Redisx redis.RedisConf

	// MongoDB 配置（用于消息持久化）
	Mongo struct {
		Url string // MongoDB 连接地址
		Db  string // 数据库名称
	}

	// 群聊已读消息处理配置
	MsgReadHandler struct {
		GroupMsgReadHandler          int   // 群聊已读消息处理模式：0-立即推送，1-延迟合并推送
		GroupMsgReadRecordDelayTime  int64 // 延迟推送时间阈值（秒）
		GroupMsgReadRecordDelayCount int   // 延迟推送数量阈值（条）
	}

	// Social RPC 客户端配置（用于查询群成员信息）
	SocialRpc zrpc.RpcClientConf

	// WebSocket 服务配置（用于推送消息）
	Ws struct {
		Host string // WebSocket 服务地址
	}
}
