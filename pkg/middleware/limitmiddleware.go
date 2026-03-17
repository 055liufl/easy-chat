// =============================================================================
// Limit Middleware - HTTP 限流中间件
// =============================================================================
// 提供基于令牌桶算法的 HTTP 请求限流功能。
//
// 功能特性:
//   - 令牌桶限流算法
//   - 基于 Redis 实现分布式限流
//   - 支持自定义速率和突发流量
//   - 超过限制时返回 401 状态码
//
// 使用场景:
//   - API 接口限流
//   - 防止恶意攻击
//   - 保护服务稳定性
//   - 流量控制
//
// 设计思路:
//   - 使用令牌桶算法，支持突发流量
//   - 基于 Redis 实现分布式限流
//   - 超过限制时快速拒绝请求
//
// 项目中的应用:
//   - 所有 HTTP API
//   - 公开接口限流
//   - 防止刷接口
//
// 令牌桶算法:
//   - rate: 令牌生成速率（每秒生成的令牌数）
//   - burst: 桶容量（最多存储的令牌数）
//   - 请求到来时，尝试从桶中获取令牌
//   - 如果获取成功，允许请求通过
//   - 如果获取失败，拒绝请求
//
// 第三方依赖:
//   - github.com/zeromicro/go-zero/core/limit: 令牌桶限流器
//
// =============================================================================
package middleware

import (
	"github.com/zeromicro/go-zero/core/limit"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"net/http"
)

// LimitMiddleware HTTP 限流中间件
type LimitMiddleware struct {
	redisCfg redis.RedisConf // Redis 配置
	*limit.TokenLimiter       // 令牌桶限流器
}

// NewLimitMiddleware 创建限流中间件实例
//
// 参数:
//   - cfg: Redis 配置
//
// 返回:
//   - *LimitMiddleware: 中间件实例
//
// 示例:
//   limiter := middleware.NewLimitMiddleware(redis.RedisConf{
//       Host: "127.0.0.1:6379",
//       Type: "node",
//   })
func NewLimitMiddleware(cfg redis.RedisConf) *LimitMiddleware {
	return &LimitMiddleware{redisCfg: cfg}
}

func (m *LimitMiddleware) TokenLimitHandler(rate, burst int) rest.Middleware {
	m.TokenLimiter = limit.NewTokenLimiter(rate, burst, redis.MustNewRedis(m.redisCfg), "REDIS_TOKEN_LIMIT_KEY")

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if m.TokenLimiter.AllowCtx(r.Context()) {
				next(w, r)
				return
			}

			w.WriteHeader(http.StatusUnauthorized)
		}
	}
}
