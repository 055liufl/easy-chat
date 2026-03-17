// =============================================================================
// SyncxLimit Interceptor - RPC 服务端并发限流拦截器
// =============================================================================
// 提供 gRPC 服务端并发连接数限制功能，防止服务过载。
//
// 功能特性:
//   - 并发连接数限制
//   - 超过限制时快速拒绝请求
//   - 基于信号量实现
//   - 轻量级，性能开销小
//
// 使用场景:
//   - 限制服务端最大并发连接数
//   - 防止服务过载
//   - 保护下游资源（如数据库连接池）
//   - 高并发场景下的流量控制
//
// 设计思路:
//   - 使用 go-zero 的 syncx.Limit 实现信号量
//   - 请求到来时尝试获取信号量
//   - 如果获取成功，执行请求并释放信号量
//   - 如果获取失败，直接拒绝请求
//
// 项目中的应用:
//   - 所有 RPC 服务端
//   - 防止服务过载
//   - 保护系统稳定性
//
// 与降载的区别:
//   - 限流：硬性限制并发数，超过直接拒绝
//   - 降载：根据响应时间动态调整通过率
//
// 第三方依赖:
//   - github.com/zeromicro/go-zero/core/syncx: 信号量实现
//
// =============================================================================
package rpcserver

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/syncx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

// SyncxLimitInterceptor 创建并发限流拦截器
// 限制 RPC 服务端的最大并发连接数
//
// 参数:
//   - maxCount: 最大并发连接数
//
// 返回:
//   - grpc.UnaryServerInterceptor: gRPC 服务端拦截器
//
// 工作流程:
//   1. 请求到来时，尝试获取信号量（TryBorrow）
//   2. 如果获取成功，执行请求并在完成后释放信号量（Return）
//   3. 如果获取失败，记录日志并返回 Unavailable 错误
//
// 使用场景:
//   - 所有 RPC 服务端
//
// 示例:
//   server := grpc.NewServer(
//       grpc.UnaryInterceptor(rpcserver.SyncxLimitInterceptor(1000)),
//   )
//
// 注意:
//   - maxCount 应根据服务器资源合理设置
//   - 设置过小会导致请求被拒绝
//   - 设置过大会导致服务过载
func SyncxLimitInterceptor(maxCount int) grpc.UnaryServerInterceptor {
	// 创建信号量，限制最大并发数
	l := syncx.NewLimit(maxCount)
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		// 尝试获取信号量
		if l.TryBorrow() {
			defer func() {
				// 释放信号量
				if err := l.Return(); err != nil {
					logx.Error(err)
				}
			}()
			// 执行请求
			return handler(ctx, req)
		} else {
			// 超过并发限制，拒绝请求
			logx.Errorf("concurrent connections over %d, rejected with code %d",
				maxCount, http.StatusServiceUnavailable)
			return nil, status.Error(codes.Unavailable, "concurrent connections over limit")
		}
	}
}
