// =============================================================================
// WebSocket Connection - WebSocket 连接封装
// =============================================================================
// 对 gorilla/websocket 的连接进行封装，提供额外的功能：
//   - 空闲连接检测和自动关闭
//   - 消息确认队列管理
//   - 线程安全的读写操作
//   - 用户 ID 绑定
//
// 连接生命周期:
//  1. 创建连接（NewConn）
//  2. 启动 keepalive 协程监控空闲状态
//  3. 处理消息读写
//  4. 超时或主动关闭连接
//
// 空闲检测机制:
//   - 每次读写操作都会更新 idle 时间
//   - keepalive 协程定期检查空闲时长
//   - 超过 maxConnectionIdle 时自动关闭连接
//
// =============================================================================
package websocket

import (
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

// Conn WebSocket 连接封装
// 在原生 WebSocket 连接基础上增加了空闲检测、消息队列等功能
type Conn struct {
	idleMu sync.Mutex // 空闲时间锁，保护 idle 字段的并发访问

	Uid string // 用户 ID，标识该连接属于哪个用户

	*websocket.Conn // 底层 WebSocket 连接
	s               *Server // 所属的服务器实例

	idle              time.Time     // 最后一次活动时间，用于空闲检测
	maxConnectionIdle time.Duration // 最大空闲时长，超过后自动关闭连接

	messageMu      sync.Mutex         // 消息队列锁，保护 readMessage 和 readMessageSeq
	readMessage    []*Message         // 待确认消息队列（按接收顺序）
	readMessageSeq map[string]*Message // 消息 ID 到消息的映射，用于快速查找

	message chan *Message // 消息通道，用于传递已确认的消息到业务处理

	done chan struct{} // 关闭信号通道，用于通知协程退出
}

// NewConn 创建新的 WebSocket 连接
// 将 HTTP 连接升级为 WebSocket 连接，并初始化相关字段
//
// 参数:
//   - s: 服务器实例
//   - w: HTTP 响应写入器
//   - r: HTTP 请求对象
//
// 返回:
//   - *Conn: 初始化的连接对象，如果升级失败返回 nil
func NewConn(s *Server, w http.ResponseWriter, r *http.Request) *Conn {
	// 将 HTTP 连接升级为 WebSocket
	c, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Errorf("upgrade err %v", err)
		return nil
	}

	conn := &Conn{
		Conn:              c,
		s:                 s,
		idle:              time.Now(),
		maxConnectionIdle: s.opt.maxConnectionIdle,
		readMessage:       make([]*Message, 0, 2),
		readMessageSeq:    make(map[string]*Message, 2),
		message:           make(chan *Message, 1),
		done:              make(chan struct{}),
	}

	// 启动空闲连接检测协程
	go conn.keepalive()
	return conn
}

// appendMsgMq 将消息添加到待确认队列
// 处理消息去重和 ACK 序号验证
//
// 参数:
//   - msg: 要添加的消息
func (c *Conn) appendMsgMq(msg *Message) {
	c.messageMu.Lock()
	defer c.messageMu.Unlock()

	// 检查消息是否已在队列中
	if m, ok := c.readMessageSeq[msg.Id]; ok {
		// 消息已存在，检查是否需要更新 ACK 序号
		if len(c.readMessage) == 0 {
			// 队列为空，说明消息已被处理
			return
		}

		// 检查 ACK 序号是否有更新
		if m.AckSeq >= msg.AckSeq {
			// ACK 序号没有增加，可能是重复消息
			return
		}

		// 更新 ACK 序号
		c.readMessageSeq[msg.Id] = msg
		return
	}

	// 新消息，但如果是 ACK 帧则忽略（避免客户端重复发送 ACK）
	if msg.FrameType == FrameAck {
		return
	}

	// 将新消息添加到队列
	c.readMessage = append(c.readMessage, msg)
	c.readMessageSeq[msg.Id] = msg

}

// ReadMessage 读取消息
// 封装了原生的 ReadMessage，增加了空闲时间更新
//
// 返回:
//   - messageType: 消息类型
//   - p: 消息内容
//   - err: 读取错误
func (c *Conn) ReadMessage() (messageType int, p []byte, err error) {
	messageType, p, err = c.Conn.ReadMessage()

	// 更新空闲时间（读取消息表示连接活跃）
	c.idleMu.Lock()
	defer c.idleMu.Unlock()
	c.idle = time.Time{} // 设置为零值表示连接活跃
	return
}

// WriteMessage 写入消息
// 封装了原生的 WriteMessage，增加了空闲时间更新
// 注意：此方法不是并发安全的，需要外部加锁
//
// 参数:
//   - messageType: 消息类型
//   - data: 消息内容
//
// 返回:
//   - error: 写入错误
func (c *Conn) WriteMessage(messageType int, data []byte) error {
	c.idleMu.Lock()
	defer c.idleMu.Unlock()

	// 写入消息（注意：此方法不是并发安全的）
	err := c.Conn.WriteMessage(messageType, data)
	// 更新空闲时间
	c.idle = time.Now()
	return err
}

// Close 关闭连接
// 发送关闭信号并关闭底层 WebSocket 连接
//
// 返回:
//   - error: 关闭错误
func (c *Conn) Close() error {
	// 关闭 done 通道，通知所有协程退出
	select {
	case <-c.done:
		// 已经关闭
	default:
		close(c.done)
	}

	// 关闭底层连接
	return c.Conn.Close()
}

// keepalive 空闲连接检测
// 定期检查连接的空闲时长，超过阈值时自动关闭连接
func (c *Conn) keepalive() {
	idleTimer := time.NewTimer(c.maxConnectionIdle)
	defer func() {
		idleTimer.Stop()
	}()

	for {
		select {
		case <-idleTimer.C:
			// 定时器触发，检查空闲时长
			c.idleMu.Lock()
			idle := c.idle
			if idle.IsZero() {
				// 连接活跃（idle 为零值表示最近有读写操作）
				c.idleMu.Unlock()
				idleTimer.Reset(c.maxConnectionIdle)
				continue
			}
			// 计算剩余空闲时间
			val := c.maxConnectionIdle - time.Since(idle)
			c.idleMu.Unlock()

			if val <= 0 {
				// 空闲时间超过阈值，关闭连接
				c.s.Close(c)
				return
			}
			// 重置定时器，等待剩余时间
			idleTimer.Reset(val)
		case <-c.done:
			// 连接已关闭，退出协程
			return
		}
	}
}
