// =============================================================================
// WebSocket Client - WebSocket 客户端实现
// =============================================================================
// 提供 WebSocket 客户端功能，用于连接 WebSocket 服务器并进行通信
//
// 功能特性:
//   - 连接管理（建立连接、断线重连）
//   - 消息发送（自动序列化为 JSON）
//   - 消息接收（自动反序列化 JSON）
//   - 自定义连接选项（路径、请求头等）
//
// 使用场景:
//   - 测试 WebSocket 服务器
//   - 作为客户端连接其他 WebSocket 服务
//   - 服务间通信
//
// =============================================================================
package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/url"
)

// Client WebSocket 客户端接口
// 定义了客户端的基本操作方法
type Client interface {
	// Close 关闭连接
	Close() error

	// Send 发送消息
	Send(v any) error

	// Read 读取消息
	Read(v any) error
}

// client WebSocket 客户端实现
type client struct {
	*websocket.Conn        // 底层 WebSocket 连接
	host            string // 服务器地址，格式: "host:port"

	opt dailOption // 连接选项（路径、请求头等）
}

// NewClient 创建 WebSocket 客户端实例
// 立即建立与服务器的连接，如果连接失败会 panic
//
// 参数:
//   - host: 服务器地址，格式为 "host:port"，如 "localhost:8080"
//   - opts: 可变参数，用于配置连接选项（路径、请求头等）
//
// 返回:
//   - *client: 已连接的客户端实例
//
// 注意:
//   - 如果连接失败会 panic，调用方需要处理
func NewClient(host string, opts ...DailOptions) *client {
	opt := newDailOptions(opts...)

	c := client{
		Conn: nil,
		host: host,
		opt:  opt,
	}

	// 建立连接
	conn, err := c.dail()
	if err != nil {
		panic(err)
	}

	c.Conn = conn
	return &c
}

// dail 建立 WebSocket 连接
// 使用 ws:// 协议连接到服务器
//
// 返回:
//   - *websocket.Conn: WebSocket 连接对象
//   - error: 连接失败时返回错误
func (c *client) dail() (*websocket.Conn, error) {
	// 构建 WebSocket URL
	u := url.URL{Scheme: "ws", Host: c.host, Path: c.opt.pattern}
	// 使用默认拨号器建立连接
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), c.opt.header)
	return conn, err
}

// Send 发送消息到服务器
// 将消息对象序列化为 JSON 后发送
// 如果发送失败会尝试重连后再次发送
//
// 参数:
//   - v: 要发送的消息对象（会被序列化为 JSON）
//
// 返回:
//   - error: 发送失败时返回错误
func (c *client) Send(v any) error {
	// 序列化消息为 JSON
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	// 尝试发送消息
	err = c.WriteMessage(websocket.TextMessage, data)
	if err == nil {
		return nil
	}

	// 发送失败，尝试重连
	// TODO: 可以增加重连次数限制和退避策略
	conn, err := c.dail()
	if err != nil {
		return err
	}
	c.Conn = conn

	// 重连后再次发送
	return c.WriteMessage(websocket.TextMessage, data)
}

// Read 从服务器读取消息
// 读取消息并反序列化为指定的对象类型
//
// 参数:
//   - v: 用于接收消息的对象指针（会被反序列化填充）
//
// 返回:
//   - error: 读取或反序列化失败时返回错误
func (c *client) Read(v any) error {
	// 读取消息
	_, msg, err := c.Conn.ReadMessage()
	if err != nil {
		return err
	}

	// 反序列化 JSON 消息
	return json.Unmarshal(msg, v)
}
