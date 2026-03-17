// =============================================================================
// User API 服务上下文
// =============================================================================
// 定义用户 API 服务的上下文结构，包括:
//   - 配置信息（Config）
//   - Redis 客户端（用于缓存和在线状态管理）
//   - User RPC 客户端（用于调用用户 RPC 服务）
//
// 数据来源:
//   在服务启动时初始化，从配置文件加载配置并创建各种客户端连接
//
// 业务场景:
//   为所有 HTTP Handler 和 Logic 提供统一的依赖注入，包括配置、
//   缓存客户端和 RPC 客户端等
//
// 重试策略:
//   为 User RPC 客户端配置了重试策略，最多重试 5 次，初始退避 1ms，
//   最大退避 2ms，退避倍数 1.0，对 UNKNOWN 状态码进行重试
//
// =============================================================================
package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"imooc.com/easy-chat/apps/user/api/internal/config"
	"imooc.com/easy-chat/apps/user/rpc/userclient"
	// N * client =》 别名
)

// retryPolicy gRPC 重试策略配置（JSON 格式）
// 配置说明:
//   - service: 指定重试的服务名称（user.User）
//   - waitForReady: 等待服务就绪后再发送请求
//   - maxAttempts: 最大重试次数（5 次）
//   - initialBackoff: 初始退避时间（1ms）
//   - maxBackoff: 最大退避时间（2ms）
//   - backoffMultiplier: 退避倍数（1.0，即固定退避时间）
//   - retryableStatusCodes: 可重试的状态码（UNKNOWN）
var retryPolicy = `{
	"methodConfig" : [{
		"name": [{
			"service": "user.User"
		}],
		"waitForReady": true,
		"retryPolicy": {
			"maxAttempts": 5,
			"initialBackoff": "0.001s",
			"maxBackoff": "0.002s",
			"backoffMultiplier": 1.0,
			"retryableStatusCodes": ["UNKNOWN"]
		}
	}]
}`

// ServiceContext 服务上下文结构
// 包含服务运行所需的所有依赖
type ServiceContext struct {
	Config config.Config // 服务配置

	*redis.Redis      // Redis 客户端（用于缓存和在线用户状态管理）
	userclient.User   // User RPC 客户端（用于调用用户相关的 RPC 服务）
}

// NewServiceContext 创建服务上下文实例
// 初始化所有依赖的客户端连接
//
// 参数:
//   - c: 服务配置
//
// 返回:
//   - *ServiceContext: 初始化完成的服务上下文实例
//
// 初始化流程:
//   1. 创建 Redis 客户端连接
//   2. 创建 User RPC 客户端连接（配置重试策略）
//   3. 返回服务上下文实例
func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		Redis: redis.MustNewRedis(c.Redisx),
		User: userclient.NewUser(zrpc.MustNewClient(c.UserRpc, zrpc.WithDialOption(grpc.WithDefaultServiceConfig(
			retryPolicy)))),
	}
}
