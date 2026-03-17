// =============================================================================
// Shedding Interceptor - RPC 客户端自适应降载拦截器
// =============================================================================
// 提供 gRPC 客户端自适应降载功能，防止服务过载导致雪崩。
//
// 功能特性:
//   - 自适应降载：根据服务响应时间动态调整请求通过率
//   - 过载保护：当服务过载时，主动拒绝部分请求
//   - 统计信息：记录总请求数和丢弃请求数
//   - 快速失败：过载时立即返回错误，避免资源浪费
//
// 使用场景:
//   - 微服务调用：防止下游服务过载
//   - 高并发场景：保护服务稳定性
//   - 流量突增：自动降载保护系统
//
// 设计思路:
//   - 基于 go-zero 的自适应降载算法
//   - 使用滑动窗口统计请求响应时间
//   - 根据响应时间动态计算通过率
//   - 过载时拒绝部分请求，保护服务
//
// 项目中的应用:
//   - 所有 RPC 客户端调用
//   - 防止服务雪崩
//   - 提高系统稳定性
//
// 工作原理:
//   1. 每个请求到来时，先尝试获取通行证（Allow）
//   2. 如果获取失败，说明服务过载，直接拒绝请求
//   3. 如果获取成功，执行请求并记录响应时间
//   4. 根据响应结果调用 Pass() 或 Fail()
//   5. 降载算法根据历史数据动态调整通过率
//
// 第三方依赖:
//   - github.com/zeromicro/go-zero/core/load: 自适应降载算法
//
// =============================================================================
package rpcclient

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/load"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
)

var (
	sheddingStat *load.SheddingStat // 降载统计信息
	shedder      load.Shedder       // 降载器
	lock         sync.Mutex         // 互斥锁，保护初始化
)

// NewSheddingClient 创建自适应降载客户端拦截器
// 为 gRPC 客户端添加自适应降载功能
//
// 参数:
//   - sname: 服务名称，用于统计信息标识
//   - opts: 降载器选项（可选）
//
// 返回:
//   - grpc.UnaryClientInterceptor: gRPC 客户端拦截器
//
// 工作流程:
//   1. 请求到来时，增加总请求计数
//   2. 尝试获取通行证（Allow）
//   3. 如果获取失败，增加丢弃计数并返回错误
//   4. 如果获取成功，执行 RPC 调用
//   5. 根据响应结果调用 Pass() 或 Fail()
//
// 使用场景:
//   - 所有 RPC 客户端调用
//
// 示例:
//   conn, err := grpc.Dial(
//       target,
//       grpc.WithUnaryInterceptor(rpcclient.NewSheddingClient("user-service")),
//   )
func NewSheddingClient(sname string, opts ...load.ShedderOption) grpc.UnaryClientInterceptor {
	ensureShedding(sname, opts...)

	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		// 增加总请求计数
		sheddingStat.IncrementTotal()
		var promise load.Promise
		// 尝试获取通行证
		promise, err = shedder.Allow()
		if err != nil {
			// 服务过载，拒绝请求
			sheddingStat.IncrementDrop()
			fmt.Println("---- sheddingStat.IncrementDrop() --------- ")
			return
		}
		fmt.Println("---- shedder.Allow() --------- ", err)
		defer func() {
			// 根据响应结果更新降载算法
			if Acceptable(err) {
				fmt.Println("---- NewSheddingClient --- acceptable promise.Pass() ", err)
				promise.Pass()
			} else {
				promise.Fail()
				fmt.Println("---- NewSheddingClient --- acceptable promise.Fail()")
			}
		}()

		// 执行 RPC 调用
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// ensureShedding 确保降载器已初始化
// 使用单例模式初始化降载器和统计信息
//
// 参数:
//   - sname: 服务名称
//   - opts: 降载器选项
func ensureShedding(sname string, opts ...load.ShedderOption) {
	lock.Lock()
	if sheddingStat == nil {
		sheddingStat = load.NewSheddingStat(sname)
	}

	if shedder == nil {
		shedder = load.NewAdaptiveShedder(opts...)
	}
	lock.Unlock()
}

// Acceptable 判断错误是否可接受
// 根据 gRPC 错误码判断请求是否成功
//
// 参数:
//   - err: gRPC 错误
//
// 返回:
//   - bool: true 表示请求成功或可接受的错误，false 表示请求失败
//
// 不可接受的错误码:
//   - DeadlineExceeded: 请求超时
//   - Internal: 内部错误
//   - Unavailable: 服务不可用
//   - DataLoss: 数据丢失
//   - Unimplemented: 方法未实现
func Acceptable(err error) bool {
	switch status.Code(err) {
	case codes.DeadlineExceeded, codes.Internal, codes.Unavailable, codes.DataLoss, codes.Unimplemented:
		return false
	default:
		return true
	}
}
