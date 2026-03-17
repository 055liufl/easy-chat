// =============================================================================
// Job Option - 任务重试选项配置
// =============================================================================
// 提供任务重试的选项配置，支持自定义重试策略。
//
// 功能特性:
//   - 自定义重试超时时间
//   - 自定义重试次数
//   - 自定义重试判断函数
//   - 自定义重试间隔函数
//
// 使用场景:
//   - 配置任务重试策略
//   - 根据不同错误类型决定是否重试
//   - 实现指数退避重试
//   - 实现自适应重试间隔
//
// 设计思路:
//   - 使用选项模式（Functional Options Pattern）
//   - 提供灵活的配置方式
//   - 支持链式调用
//
// 项目中的应用:
//   - 所有需要重试的任务
//   - 自定义重试策略
//
// =============================================================================
package job

import "time"

type (
	// RetryOptions 重试选项函数类型
	// 用于配置重试选项
	RetryOptions func(opts *retryOptions)

	// retryOptions 重试选项结构
	// 包含所有重试相关的配置
	retryOptions struct {
		timeout     time.Duration   // 单次任务执行超时时间
		retryNums   int             // 最大重试次数
		isRetryFunc IsRetryFunc     // 判断是否需要重试的函数
		retryJetLag RetryJetLagFunc // 计算重试间隔的函数
	}
)

// newOptions 创建默认重试选项
// 使用默认配置初始化重试选项，然后应用自定义选项
//
// 参数:
//   - opts: 自定义选项列表
//
// 返回:
//   - *retryOptions: 重试选项实例
func newOptions(opts ...RetryOptions) *retryOptions {
	opt := &retryOptions{
		timeout:     DefaultRetryTimeout,
		retryNums:   DefaultRetryNums,
		isRetryFunc: RetryAlways,
		retryJetLag: RetryJetLagAlways,
	}

	for _, options := range opts {
		options(opt)
	}
	return opt
}

// WithRetryTimeout 设置重试超时时间
// 配置单次任务执行的最大超时时间
//
// 参数:
//   - timeout: 超时时间
//
// 返回:
//   - RetryOptions: 选项函数
//
// 示例:
//   WithRetry(ctx, handler, WithRetryTimeout(5*time.Second))
func WithRetryTimeout(timeout time.Duration) RetryOptions {
	return func(opts *retryOptions) {
		if timeout > 0 {
			opts.timeout = timeout
		}
	}
}

// WithRetryNums 设置重试次数
// 配置任务失败后的最大重试次数
//
// 参数:
//   - nums: 重试次数
//
// 返回:
//   - RetryOptions: 选项函数
//
// 示例:
//   WithRetry(ctx, handler, WithRetryNums(10))
func WithRetryNums(nums int) RetryOptions {
	return func(opts *retryOptions) {
		opts.retryNums = 1

		if nums > 1 {
			opts.retryNums = nums
		}
	}
}

// WithIsRetryFunc 设置重试判断函数
// 配置判断是否需要重试的函数
//
// 参数:
//   - retryFunc: 重试判断函数
//
// 返回:
//   - RetryOptions: 选项函数
//
// 示例:
//   WithRetry(ctx, handler, WithIsRetryFunc(func(ctx context.Context, retryCount int, err error) bool {
//       // 只有网络错误才重试
//       return isNetworkError(err)
//   }))
func WithIsRetryFunc(retryFunc IsRetryFunc) RetryOptions {
	return func(opts *retryOptions) {
		if retryFunc != nil {
			opts.isRetryFunc = retryFunc
		}
	}
}

// WithRetryJetLagFunc 设置重试间隔函数
// 配置计算重试间隔的函数
//
// 参数:
//   - retryJetLagFunc: 重试间隔函数
//
// 返回:
//   - RetryOptions: 选项函数
//
// 示例:
//   // 指数退避重试
//   WithRetry(ctx, handler, WithRetryJetLagFunc(func(ctx context.Context, retryCount int, lastTime time.Duration) time.Duration {
//       return time.Duration(math.Pow(2, float64(retryCount))) * time.Second
//   }))
func WithRetryJetLagFunc(retryJetLagFunc RetryJetLagFunc) RetryOptions {
	return func(opts *retryOptions) {
		if retryJetLagFunc != nil {
			opts.retryJetLag = retryJetLagFunc
		}
	}
}
