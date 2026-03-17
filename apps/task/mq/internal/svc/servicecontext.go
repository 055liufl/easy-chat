// =============================================================================
// 服务上下文 - 依赖注入容器
// =============================================================================
// 本模块负责初始化和管理 Task 服务的所有依赖，包括：
//   - 配置信息
//   - Redis 客户端
//   - MongoDB 模型（ChatLog、Conversation）
//   - Social RPC 客户端（查询群成员）
//   - WebSocket 客户端（推送消息）
//
// 业务场景:
//   - 服务启动时初始化所有依赖
//   - 通过依赖注入方式提供给各个处理器使用
//   - 统一管理资源生命周期
//
// =============================================================================

package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"imooc.com/easy-chat/apps/im/immodels"
	"imooc.com/easy-chat/apps/im/ws/websocket"
	"imooc.com/easy-chat/apps/social/rpc/socialclient"
	"imooc.com/easy-chat/apps/task/mq/internal/config"
	"imooc.com/easy-chat/pkg/constants"
	"net/http"
)

// ServiceContext 服务上下文
// 包含 Task 服务运行所需的所有依赖
type ServiceContext struct {
	config.Config // 配置信息（嵌入）

	WsClient websocket.Client // WebSocket 客户端，用于推送消息给在线用户
	*redis.Redis              // Redis 客户端，用于缓存和分布式锁

	socialclient.Social       // Social RPC 客户端，用于查询群成员信息
	immodels.ChatLogModel     // ChatLog MongoDB 模型，用于消息持久化
	immodels.ConversationModel // Conversation MongoDB 模型，用于会话管理
}

// NewServiceContext 创建服务上下文实例
// 初始化所有依赖并返回服务上下文
//
// 初始化流程:
//   1. 初始化 Redis 客户端
//   2. 初始化 MongoDB 模型（ChatLog、Conversation）
//   3. 初始化 Social RPC 客户端
//   4. 获取系统 Token（用于 WebSocket 认证）
//   5. 初始化 WebSocket 客户端
//
// 参数:
//   - c: 配置信息
//
// 返回:
//   - *ServiceContext: 服务上下文实例
func NewServiceContext(c config.Config) *ServiceContext {
	svc := &ServiceContext{
		Config:            c,
		Redis:             redis.MustNewRedis(c.Redisx),
		ChatLogModel:      immodels.MustChatLogModel(c.Mongo.Url, c.Mongo.Db),
		ConversationModel: immodels.MustConversationModel(c.Mongo.Url, c.Mongo.Db),

		Social: socialclient.NewSocial(zrpc.MustNewClient(c.SocialRpc)),
	}

	// 获取系统 Token（用于 WebSocket 认证）
	token, err := svc.GetSystemToken()
	if err != nil {
		panic(err)
	}

	// 初始化 WebSocket 客户端（带认证 Header）
	header := http.Header{}
	header.Set("Authorization", token)
	svc.WsClient = websocket.NewClient(c.Ws.Host, websocket.WithClientHeader(header))
	return svc
}

// GetSystemToken 获取系统 Token
// 从 Redis 中获取系统根用户的 Token，用于 WebSocket 认证
//
// 返回:
//   - string: 系统 Token
//   - error: 获取失败时返回错误
func (svc *ServiceContext) GetSystemToken() (string, error) {
	return svc.Redis.Get(constants.REDIS_SYSTEM_ROOT_TOKEN)
}
