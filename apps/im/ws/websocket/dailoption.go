// =============================================================================
// WebSocket Client Options - WebSocket 客户端配置选项
// =============================================================================
// 提供客户端的配置选项，使用函数式选项模式（Functional Options Pattern）
//
// 可配置项:
//   - Pattern: WebSocket 连接路径（如 "/ws"）
//   - Header: HTTP 请求头（可用于传递认证信息等）
//
// 使用示例:
//   client := NewClient("localhost:8080",
//       WithClientPatten("/chat"),
//       WithClientHeader(http.Header{"Authorization": []string{"Bearer token"}}),
//   )
//
// =============================================================================
package websocket

import "net/http"

// DailOptions 客户端配置选项函数类型
// 使用函数式选项模式，允许灵活配置客户端参数
type DailOptions func(option *dailOption)

// dailOption 客户端配置选项结构
// 包含客户端连接所需的配置参数
type dailOption struct {
	pattern string      // WebSocket 连接路径（如 "/ws"）
	header  http.Header // HTTP 请求头，可用于传递认证信息、自定义头等
}

// newDailOptions 创建客户端配置选项
// 合并默认配置和用户提供的配置
//
// 参数:
//   - opts: 可变参数，用户提供的配置选项
//
// 返回:
//   - dailOption: 合并后的配置选项
func newDailOptions(opts ...DailOptions) dailOption {
	// 设置默认配置
	o := dailOption{
		pattern: "/ws",
		header:  nil,
	}

	// 应用用户提供的配置
	for _, opt := range opts {
		opt(&o)
	}

	return o
}

// WithClientPatten 设置 WebSocket 连接路径
//
// 参数:
//   - pattern: WebSocket 路径（如 "/ws", "/chat"）
//
// 返回:
//   - DailOptions: 配置选项函数
func WithClientPatten(pattern string) DailOptions {
	return func(opt *dailOption) {
		opt.pattern = pattern
	}
}

// WithClientHeader 设置 HTTP 请求头
// 可用于传递认证信息（如 Authorization）、自定义头等
//
// 参数:
//   - header: HTTP 请求头
//
// 返回:
//   - DailOptions: 配置选项函数
func WithClientHeader(header http.Header) DailOptions {
	return func(opt *dailOption) {
		opt.header = header
	}
}


