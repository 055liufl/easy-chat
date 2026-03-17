// =============================================================================
// Job Retry - 任务重试机制
// =============================================================================
// 提供通用的任务重试机制，支持自定义重试策略。
//
// 功能特性:
//   - 自动重试失败的任务
//   - 支持超时控制
//   - 支持自定义重试次数
//   - 支持自定义重试间隔
//   - 支持自定义重试条件
//
// 使用场景:
//   - 网络请求重试
//   - 数据库操作重试
//   - 外部服务调用重试
//   - 消息发送重试
//
// 设计思路:
//   - 使用 goroutine 异步执行任务
//   - 使用 channel 接收任务结果
//   - 使用 context 控制超时
//   - 支持自定义重试策略
//
// 项目中的应用:
//   - RPC 调用重试
//   - 消息发送重试
//   - 数据同步重试
//
// =============================================================================
package job

import (
	"context"
	"errors"
	"time"
)

// ErrJobTimeout 任务超时错误
var ErrJobTimeout = errors.New("任务超时")

// RetryJetLagFunc 重试间隔计算函数类型
// 根据重试次数和上次间隔时间，计算下次重试的间隔时间
//
// 参数:
//   - ctx: 上下文对象
//   - retryCount: 当前重试次数（从 0 开始）
//   - lastTime: 上次重试的间隔时间
//
// 返回:
//   - time.Duration: 下次重试的间隔时间
type RetryJetLagFunc func(ctx context.Context, retryCount int, lastTime time.Duration) time.Duration

// RetryJetLagAlways 固定间隔重试策略
// 每次重试使用固定的间隔时间（DefaultRetryJetLag）
//
// 参数:
//   - ctx: 上下文对象
//   - retryCount: 当前重试次数
//   - lastTime: 上次重试的间隔时间
//
// 返回:
//   - time.Duration: 固定的重试间隔时间
func RetryJetLagAlways(ctx context.Context, retryCount int, lastTime time.Duration) time.Duration {
	return DefaultRetryJetLag
}

// IsRetryFunc 重试判断函数类型
// 根据重试次数和错误信息，判断是否需要继续重试
//
// 参数:
//   - ctx: 上下文对象
//   - retryCount: 当前重试次数（从 0 开始）
//   - err: 任务执行错误
//
// 返回:
//   - bool: true 表示需要重试，false 表示不重试
type IsRetryFunc func(ctx context.Context, retryCount int, err error) bool

// RetryAlways 总是重试策略
// 无论什么错误都进行重试
//
// 参数:
//   - ctx: 上下文对象
//   - retryCount: 当前重试次数
//   - err: 任务执行错误
//
// 返回:
//   - bool: 总是返回 true
func RetryAlways(ctx context.Context, retryCount int, err error) bool {
	return true
}

// WithRetry 执行带重试的任务
// 自动重试失败的任务，支持超时控制和自定义重试策略
//
// 参数:
//   - ctx: 上下文对象，用于超时控制
//   - handler: 任务处理函数，返回 nil 表示成功，返回 error 表示失败
//   - opts: 重试选项（可选）
//
// 返回:
//   - error: 任务执行错误，如果所有重试都失败则返回最后一次的错误
//
// 工作流程:
//   1. 检查上下文是否已设置超时，如果没有则使用默认超时
//   2. 循环执行任务，最多重试 retryNums 次
//   3. 每次执行任务使用 goroutine 异步执行
//   4. 使用 select 等待任务完成或超时
//   5. 如果任务成功，直接返回
//   6. 如果任务失败，判断是否需要重试
//   7. 如果需要重试，等待重试间隔后继续
//   8. 如果超时，返回超时错误
//
// 使用场景:
//   - 网络请求重试
//   - 数据库操作重试
//   - 外部服务调用重试
//
// 示例:
//   // 使用默认配置
//   err := job.WithRetry(ctx, func(ctx context.Context) error {
//       return sendMessage(ctx, msg)
//   })
//
//   // 自定义配置
//   err := job.WithRetry(ctx, func(ctx context.Context) error {
//       return sendMessage(ctx, msg)
//   },
//       job.WithRetryNums(10),
//       job.WithRetryTimeout(5*time.Second),
//       job.WithIsRetryFunc(func(ctx context.Context, retryCount int, err error) bool {
//           // 只有网络错误才重试
//           return isNetworkError(err)
//       }),
//   )
//
// 注意:
//   - 如果上下文已设置超时，会使用上下文的超时时间
//   - 重试间隔会累加到总超时时间中
//   - 任务处理函数应该是幂等的，避免重复执行导致问题
func WithRetry(ctx context.Context, handler func(ctx context.Context) error, opts ...RetryOptions) error {
	opt := newOptions(opts...)

	// 判断程序本身是否设置了超时
	_, ok := ctx.Deadline()
	if !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, opt.timeout)
		defer cancel()
	}

	var (
		herr        error
		retryJetLag time.Duration
		ch          = make(chan error, 1)
	)

	for i := 0; i < opt.retryNums; i++ {
		go func() {
			ch <- handler(ctx)
		}()

		select {
		case herr = <-ch:
			if herr == nil {
				return nil
			}

			if !opt.isRetryFunc(ctx, i, herr) {
				return herr
			}

			retryJetLag = opt.retryJetLag(ctx, i, retryJetLag)
			time.Sleep(retryJetLag)
		case <-ctx.Done():
			return ErrJobTimeout
		}
	}

	return herr
}
