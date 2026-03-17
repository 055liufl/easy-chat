// =============================================================================
// Idempotence Middleware - HTTP 幂等性中间件
// =============================================================================
// 提供 HTTP 请求的幂等性支持，为每个请求生成唯一标识。
//
// 功能特性:
//   - 为每个 HTTP 请求生成唯一 ID
//   - 将请求 ID 存入 Context
//   - 配合 RPC 幂等性拦截器使用
//
// 使用场景:
//   - HTTP API 请求幂等性控制
//   - 防止用户重复点击
//   - 防止网络抖动导致的重复请求
//
// 设计思路:
//   - 在 HTTP 请求处理前生成请求 ID
//   - 将请求 ID 存入 Context
//   - 后续的 RPC 调用会自动携带请求 ID
//
// 项目中的应用:
//   - 所有 HTTP API
//   - 配合 RPC 幂等性拦截器实现端到端幂等性
//
// 工作流程:
//   1. HTTP 请求到达
//   2. 中间件生成请求 ID
//   3. 将请求 ID 存入 Context
//   4. 后续的 RPC 调用自动携带请求 ID
//   5. RPC 服务端根据请求 ID 进行幂等性控制
//
// =============================================================================
package middleware

import (
	"imooc.com/easy-chat/pkg/interceptor"
	"net/http"
)

// IdempotenceMiddleware HTTP 幂等性中间件
type IdempotenceMiddleware struct {
}

// NewIdempotenceMiddleware 创建幂等性中间件实例
//
// 返回:
//   - *IdempotenceMiddleware: 中间件实例
//
// 示例:
//   server := rest.MustNewServer(c.RestConf)
//   server.Use(middleware.NewIdempotenceMiddleware().Handler)
func NewIdempotenceMiddleware() *IdempotenceMiddleware {
	return &IdempotenceMiddleware{}
}

// Handler 幂等性中间件处理函数
// 为每个 HTTP 请求生成唯一 ID 并存入 Context
//
// 参数:
//   - next: 下一个处理函数
//
// 返回:
//   - http.HandlerFunc: 中间件处理函数
//
// 工作流程:
//   1. 生成请求唯一 ID（UUID）
//   2. 将 ID 存入 Context
//   3. 调用下一个处理函数
//
// 使用场景:
//   - 所有需要幂等性保证的 HTTP API
//
// 示例:
//   server.Use(idempotenceMiddleware.Handler)
func (m *IdempotenceMiddleware) Handler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 为请求生成唯一 ID 并存入 Context
		r = r.WithContext(interceptor.ContextWithVal(r.Context()))

		next(w, r)
	}
}
