// =============================================================================
// Idempotence Interceptor - RPC 幂等性拦截器
// =============================================================================
// 提供 gRPC 幂等性拦截器，防止重复请求导致的重复操作。
//
// 功能特性:
//   - 客户端拦截器：为每个请求生成唯一标识
//   - 服务端拦截器：检测并拦截重复请求
//   - 基于 Redis 实现分布式幂等性控制
//   - 支持自定义幂等性方法列表
//   - 缓存请求结果，重复请求直接返回缓存
//
// 使用场景:
//   - 防止网络抖动导致的重复请求
//   - 防止用户重复点击导致的重复操作
//   - 防止消息队列重复消费
//   - 关键业务操作的幂等性保证（如创建订单、转账）
//
// 设计思路:
//   - 客户端生成请求 ID，通过 gRPC metadata 传递
//   - 服务端使用 Redis SETNX 实现分布式锁
//   - 首次请求执行业务逻辑，结果缓存到内存
//   - 重复请求直接返回缓存结果或返回错误
//
// 项目中的应用:
//   - Social 服务：群组创建（防止重复创建）
//   - 支付服务：订单创建（防止重复下单）
//   - 消息服务：消息发送（防止重复发送）
//
// 工作流程:
//   1. 客户端：生成请求 ID -> 添加到 metadata -> 发起 RPC 调用
//   2. 服务端：提取请求 ID -> Redis SETNX 获取锁 -> 执行业务逻辑 -> 缓存结果
//   3. 重复请求：提取请求 ID -> Redis SETNX 失败 -> 返回缓存结果或错误
//
// =============================================================================
package interceptor

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"imooc.com/easy-chat/pkg/xerr"
)

// Idempotent 幂等性接口
// 定义了实现幂等性控制所需的方法
type Idempotent interface {
	// Identify 获取请求的唯一标识
	// 根据上下文和方法名生成请求的唯一 ID
	//
	// 参数:
	//   - ctx: 上下文对象，包含请求 ID
	//   - method: RPC 方法名
	//
	// 返回:
	//   - string: 请求的唯一标识
	Identify(ctx context.Context, method string) string

	// IsIdempotentMethod 判断方法是否支持幂等性
	// 检查指定的 RPC 方法是否需要幂等性控制
	//
	// 参数:
	//   - fullMethod: RPC 方法全名（如 /social.social/GroupCreate）
	//
	// 返回:
	//   - bool: true 表示需要幂等性控制，false 表示不需要
	IsIdempotentMethod(fullMethod string) bool

	// TryAcquire 尝试获取幂等性锁
	// 尝试获取分布式锁，如果获取成功则执行业务逻辑
	//
	// 参数:
	//   - ctx: 上下文对象
	//   - id: 请求的唯一标识
	//
	// 返回:
	//   - resp: 缓存的响应结果（如果存在）
	//   - isAcquire: 是否成功获取锁（true 表示首次请求，false 表示重复请求）
	TryAcquire(ctx context.Context, id string) (resp interface{}, isAcquire bool)

	// SaveResp 保存请求的响应结果
	// 将请求的响应结果缓存，供重复请求使用
	//
	// 参数:
	//   - ctx: 上下文对象
	//   - id: 请求的唯一标识
	//   - resp: 响应结果
	//   - respErr: 响应错误
	//
	// 返回:
	//   - error: 保存结果时的错误
	SaveResp(ctx context.Context, id string, resp interface{}, respErr error) error
}

var (
	// TKey 请求任务标识键名
	// 用于在 Context 中存储请求的唯一 ID
	TKey = "easy-chat-idempotence-task-id"

	// DKey RPC 调度标识键名
	// 用于在 gRPC metadata 中传递请求的唯一 ID
	DKey = "easy-chat-idempotence-dispatch-key"
)

// ContextWithVal 创建包含请求 ID 的上下文
// 为上下文添加唯一的请求 ID，用于幂等性控制
//
// 参数:
//   - ctx: 原始上下文对象
//
// 返回:
//   - context.Context: 包含请求 ID 的新上下文
//
// 使用场景:
//   - HTTP 请求处理前，为请求生成唯一 ID
//   - RPC 调用前，为调用生成唯一 ID
//
// 示例:
//   ctx = interceptor.ContextWithVal(ctx)
//   // 后续的 RPC 调用会自动携带请求 ID
func ContextWithVal(ctx context.Context) context.Context {
	// 生成 UUID 作为请求 ID
	return context.WithValue(ctx, TKey, utils.NewUuid())
}

// NewIdempotenceClient 创建幂等性客户端拦截器
// 为 gRPC 客户端添加幂等性支持，自动为请求添加唯一标识
//
// 参数:
//   - idempotent: 幂等性接口实现
//
// 返回:
//   - grpc.UnaryClientInterceptor: gRPC 客户端拦截器
//
// 工作流程:
//   1. 从上下文中获取请求 ID
//   2. 生成请求的唯一标识（请求 ID + 方法名）
//   3. 将标识添加到 gRPC metadata 中
//   4. 发起 RPC 调用
//
// 使用场景:
//   - 所有需要幂等性保证的 RPC 客户端
//
// 示例:
//   client := social.NewSocialClient(conn,
//       grpc.WithUnaryInterceptor(interceptor.NewIdempotenceClient(idempotent)),
//   )
func NewIdempotenceClient(idempotent Idempotent) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// 获取请求的唯一标识
		identify := idempotent.Identify(ctx, method)

		// 将标识添加到 gRPC metadata 中
		ctx = metadata.NewOutgoingContext(ctx, map[string][]string{
			DKey: {identify},
		})

		// 发起 RPC 调用
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// NewIdempotenceServer 创建幂等性服务端拦截器
// 为 gRPC 服务端添加幂等性支持，自动检测并拦截重复请求
//
// 参数:
//   - idempotent: 幂等性接口实现
//
// 返回:
//   - grpc.UnaryServerInterceptor: gRPC 服务端拦截器
//
// 工作流程:
//   1. 从 gRPC metadata 中提取请求标识
//   2. 判断方法是否需要幂等性控制
//   3. 尝试获取分布式锁（Redis SETNX）
//   4. 如果获取成功，执行业务逻辑并缓存结果
//   5. 如果获取失败，返回缓存结果或错误
//
// 使用场景:
//   - 所有需要幂等性保证的 RPC 服务端
//
// 示例:
//   server := grpc.NewServer(
//       grpc.UnaryInterceptor(interceptor.NewIdempotenceServer(idempotent)),
//   )
func NewIdempotenceServer(idempotent Idempotent) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		// 从 metadata 中获取请求标识
		identify := metadata.ValueFromIncomingContext(ctx, DKey)
		if len(identify) == 0 || !idempotent.IsIdempotentMethod(info.FullMethod) {
			// 不需要幂等性处理，直接执行
			return handler(ctx, req)
		}

		fmt.Println("----", "请求进入幂等性处理", identify)

		// 尝试获取分布式锁
		r, isAcquire := idempotent.TryAcquire(ctx, identify[0])
		if isAcquire {
			// 首次请求，执行业务逻辑
			resp, err = handler(ctx, req)
			fmt.Println("---- 执行任务", identify)

			// 缓存请求结果
			if err := idempotent.SaveResp(ctx, identify[0], resp, err); err != nil {
				return resp, err
			}

			return resp, err
		}

		// 重复请求
		fmt.Println("----- 任务在执行", identify)

		if r != nil {
			// 任务已执行完成，返回缓存结果
			fmt.Println("--- 任务已经执行完了", identify)
			return r, nil
		}

		// 任务可能还在执行中，返回错误
		return nil, errors.WithStack(xerr.New(int(codes.DeadlineExceeded), fmt.Sprintf("存在其他任务在执行 id %v", identify[0])))
	}
}

var (
	// DefaultIdempotent 默认幂等性实现实例
	DefaultIdempotent = new(defaultIdempotent)

	// DefaultIdempotentClient 默认幂等性客户端拦截器
	DefaultIdempotentClient = NewIdempotenceClient(DefaultIdempotent)
)

// defaultIdempotent 默认幂等性实现
// 基于 Redis 和内存缓存实现的幂等性控制
type defaultIdempotent struct {
	*redis.Redis          // Redis 客户端，用于分布式锁
	*collection.Cache     // 内存缓存，用于存储请求结果
	method        map[string]bool // 需要幂等性控制的方法列表
}

// NewDefaultIdempotent 创建默认幂等性实现实例
// 基于 Redis 和内存缓存实现幂等性控制
//
// 参数:
//   - c: Redis 配置
//
// 返回:
//   - Idempotent: 幂等性接口实现
//
// 实现细节:
//   - 使用 Redis SETNX 实现分布式锁（有效期 1 小时）
//   - 使用内存缓存存储请求结果（有效期 1 小时）
//   - 默认支持 /social.social/GroupCreate 方法的幂等性
//
// 示例:
//   idempotent := interceptor.NewDefaultIdempotent(redis.RedisConf{
//       Host: "127.0.0.1:6379",
//       Type: "node",
//   })
func NewDefaultIdempotent(c redis.RedisConf) Idempotent {
	cache, err := collection.NewCache(60 * 60)
	if err != nil {
		panic(err)
	}

	return &defaultIdempotent{
		Redis: redis.MustNewRedis(c),
		Cache: cache,
		method: map[string]bool{
			"/social.social/GroupCreate": true, // 群组创建需要幂等性控制
		},
	}
}

// Identify 获取请求的唯一标识
// 实现 Idempotent 接口的 Identify 方法
//
// 参数:
//   - ctx: 上下文对象，包含请求 ID
//   - method: RPC 方法名
//
// 返回:
//   - string: 请求的唯一标识（格式: 请求ID.方法名）
func (d *defaultIdempotent) Identify(ctx context.Context, method string) string {
	id := ctx.Value(TKey)
	// 生成请求唯一标识: UUID.方法名
	rpcId := fmt.Sprintf("%v.%s", id, method)
	return rpcId
}

// IsIdempotentMethod 判断方法是否支持幂等性
// 实现 Idempotent 接口的 IsIdempotentMethod 方法
//
// 参数:
//   - fullMethod: RPC 方法全名
//
// 返回:
//   - bool: true 表示需要幂等性控制
func (d *defaultIdempotent) IsIdempotentMethod(fullMethod string) bool {
	return d.method[fullMethod]
}

// TryAcquire 尝试获取幂等性锁
// 实现 Idempotent 接口的 TryAcquire 方法
//
// 参数:
//   - ctx: 上下文对象
//   - id: 请求的唯一标识
//
// 返回:
//   - resp: 缓存的响应结果（如果存在）
//   - isAcquire: 是否成功获取锁
//
// 实现细节:
//   - 使用 Redis SETNX 尝试获取锁（有效期 1 小时）
//   - 如果获取成功，返回 (nil, true)，表示首次请求
//   - 如果获取失败，从内存缓存中获取结果，返回 (result, false)
func (d *defaultIdempotent) TryAcquire(ctx context.Context, id string) (resp interface{}, isAcquire bool) {
	// 基于 Redis SETNX 实现分布式锁
	retry, err := d.SetnxEx(id, "1", 60*60)
	if err != nil {
		return nil, false
	}

	if retry {
		// 获取锁成功，首次请求
		return nil, true
	}

	// 获取锁失败，从缓存中获取结果
	resp, _ = d.Cache.Get(id)
	return resp, false
}

// SaveResp 保存请求的响应结果
// 实现 Idempotent 接口的 SaveResp 方法
//
// 参数:
//   - ctx: 上下文对象
//   - id: 请求的唯一标识
//   - resp: 响应结果
//   - respErr: 响应错误
//
// 返回:
//   - error: 保存结果时的错误
//
// 实现细节:
//   - 将响应结果存储到内存缓存中
//   - 缓存有效期与 Redis 锁一致（1 小时）
func (d *defaultIdempotent) SaveResp(ctx context.Context, id string, resp interface{}, respErr error) error {
	d.Cache.Set(id, resp)
	return nil
}
