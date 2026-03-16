// =============================================================================
// WebSocket Server Options - WebSocket 服务器配置选项
// =============================================================================
// 提供服务器的配置选项，使用函数式选项模式（Functional Options Pattern）
//
// 可配置项:
//   - Authentication: 身份认证实现
//   - Patten: WebSocket 路径（如 "/ws"）
//   - Ack: 消息确认模式（NoAck/OnlyAck/RigorAck）
//   - MaxConnectionIdle: 最大空闲时长
//   - Concurrency: 并发处理任务数
//
// 使用示例:
//   server := NewServer("0.0.0.0:8080",
//       WithServerAuthentication(myAuth),
//       WithServerAck(RigorAck),
//       WithServerMaxConnectionIdle(5*time.Minute),
//   )
//
// =============================================================================
package websocket

import "time"

// ServerOptions 服务器配置选项函数类型
// 使用函数式选项模式，允许灵活配置服务器参数
type ServerOptions func(opt *serverOption)

// serverOption 服务器配置选项结构
// 包含服务器运行所需的所有配置参数
type serverOption struct {
	Authentication // 身份认证实现

	ack        AckType       // 消息确认模式
	ackTimeout time.Duration // ACK 超时时间（用于 RigorAck 模式）

	patten string // WebSocket 路径模式（如 "/ws"）

	maxConnectionIdle time.Duration // 最大空闲连接时长，超过后自动关闭

	concurrency int // 并发处理任务数
}

// newServerOptions 创建服务器配置选项
// 合并默认配置和用户提供的配置
//
// 参数:
//   - opts: 可变参数，用户提供的配置选项
//
// 返回:
//   - serverOption: 合并后的配置选项
func newServerOptions(opts ...ServerOptions) serverOption {
	// 设置默认配置
	o := serverOption{
		Authentication:    new(authentication),
		maxConnectionIdle: defaultMaxConnectionIdle,
		ackTimeout:        defaultAckTimeout,
		patten:            "/ws",
		concurrency:       defaultConcurrency,
	}

	// 应用用户提供的配置
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

// WithServerAuthentication 设置身份认证实现
//
// 参数:
//   - auth: 身份认证接口实现
//
// 返回:
//   - ServerOptions: 配置选项函数
func WithServerAuthentication(auth Authentication) ServerOptions {
	return func(opt *serverOption) {
		opt.Authentication = auth
	}
}

// WithServerPatten 设置 WebSocket 路径模式
//
// 参数:
//   - patten: WebSocket 路径（如 "/ws", "/chat"）
//
// 返回:
//   - ServerOptions: 配置选项函数
func WithServerPatten(patten string) ServerOptions {
	return func(opt *serverOption) {
		opt.patten = patten
	}
}

// WithServerAck 设置消息确认模式
//
// 参数:
//   - ack: 确认模式（NoAck/OnlyAck/RigorAck）
//
// 返回:
//   - ServerOptions: 配置选项函数
func WithServerAck(ack AckType) ServerOptions {
	return func(opt *serverOption) {
		opt.ack = ack
	}
}

// WithServerMaxConnectionIdle 设置最大空闲连接时长
// 超过此时长没有消息收发的连接会被自动关闭
//
// 参数:
//   - maxConnectionIdle: 最大空闲时长，必须大于 0
//
// 返回:
//   - ServerOptions: 配置选项函数
func WithServerMaxConnectionIdle(maxConnectionIdle time.Duration) ServerOptions {
	return func(opt *serverOption) {
		if maxConnectionIdle > 0 {
			opt.maxConnectionIdle = maxConnectionIdle
		}
	}
}
