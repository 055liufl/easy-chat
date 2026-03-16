// =============================================================================
// 服务上下文 - Social API 服务依赖管理
// =============================================================================
// 管理 Social API 服务的所有依赖，包括:
//   - 配置信息
//   - 中间件（幂等性、限流）
//   - Redis 客户端
//   - RPC 客户端（Social、User、IM）
//
// 依赖说明:
//   - Social RPC: 社交关系服务，提供好友、群组管理功能
//   - User RPC: 用户信息服务，提供用户详细信息查询
//   - IM RPC: 即时通讯服务，提供会话管理和在线状态查询
//   - Redis: 缓存服务，用于幂等性控制、限流和在线状态缓存
//
// =============================================================================
package svc

import (
	"github.com/zeromicro/go-zero/core/load"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"imooc.com/easy-chat/apps/im/rpc/imclient"
	"imooc.com/easy-chat/apps/social/api/internal/config"
	"imooc.com/easy-chat/apps/social/rpc/socialclient"
	"imooc.com/easy-chat/apps/user/rpc/userclient"
	"imooc.com/easy-chat/pkg/interceptor"
	"imooc.com/easy-chat/pkg/interceptor/rpcclient"
	"imooc.com/easy-chat/pkg/middleware"
	"time"
)

// retryPolicy Social RPC 重试策略配置
// 定义 RPC 调用失败时的重试行为
//   - maxAttempts: 最大重试次数 5 次
//   - initialBackoff: 初始退避时间 1ms
//   - maxBackoff: 最大退避时间 2ms
//   - backoffMultiplier: 退避时间倍数 1.0（固定退避）
//   - retryableStatusCodes: 可重试的错误码（UNKNOWN、DEADLINE_EXCEEDED）
var retryPolicy = `{
	"methodConfig" : [{
		"name": [{
			"service": "social.social"
		}],
		"waitForReady": true,
		"retryPolicy": {
			"maxAttempts": 5,
			"initialBackoff": "0.001s",
			"maxBackoff": "0.002s",
			"backoffMultiplier": 1.0,
			"retryableStatusCodes": ["UNKNOWN", "DEADLINE_EXCEEDED"]
		}
	}]
}`

// ServiceContext 服务上下文
// 包含 Social API 服务运行所需的所有依赖
type ServiceContext struct {
	Config                config.Config      // 配置信息
	IdempotenceMiddleware rest.Middleware    // 幂等性中间件
	LimitMiddleware       rest.Middleware    // 限流中间件
	*redis.Redis                             // Redis 客户端
	socialclient.Social                      // Social RPC 客户端
	userclient.User                          // User RPC 客户端
	imclient.Im                              // IM RPC 客户端
}

// NewServiceContext 创建服务上下文实例
// 初始化所有依赖，包括 Redis、RPC 客户端和中间件
//
// 参数:
//   - c: 配置信息
//
// 返回:
//   - *ServiceContext: 服务上下文实例
func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		// 初始化 Redis 客户端，用于缓存和分布式锁
		Redis: redis.MustNewRedis(c.Redisx),
		// 初始化幂等性中间件，防止重复请求
		IdempotenceMiddleware: middleware.NewIdempotenceMiddleware().Handler,
		// 初始化限流中间件，每秒最多 1 个请求，令牌桶容量 100
		LimitMiddleware: middleware.NewLimitMiddleware(c.Redisx).TokenLimitHandler(1, 100),
		// 初始化 Social RPC 客户端，配置幂等性拦截器
		Social: socialclient.NewSocial(zrpc.MustNewClient(c.SocialRpc,
			//zrpc.WithDialOption(grpc.WithDefaultServiceConfig(retryPolicy)), // 可选：启用重试策略
			zrpc.WithUnaryClientInterceptor(interceptor.DefaultIdempotentClient), // 幂等性拦截器
		)),

		// 初始化 User RPC 客户端，配置自适应降载保护
		User: userclient.NewUser(zrpc.MustNewClient(c.UserRpc,
			zrpc.WithUnaryClientInterceptor(rpcclient.NewSheddingClient("user-rpc",
				load.WithBuckets(10),                      // 滑动窗口桶数量
				load.WithCpuThreshold(1),                  // CPU 阈值（100%）
				load.WithWindow(time.Millisecond*100000),  // 滑动窗口时间
			)),
		)),
		// 初始化 IM RPC 客户端，用于会话管理和在线状态查询
		Im: imclient.NewIm(zrpc.MustNewClient(c.ImRpc)),
	}
}
