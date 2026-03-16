// =============================================================================
// WebSocket Server - WebSocket 服务器核心实现
// =============================================================================
// 提供基于 WebSocket 的实时通信服务，支持：
//   - 用户连接管理（连接建立、认证、关闭）
//   - 消息路由分发（根据 method 字段路由到对应处理器）
//   - 消息确认机制（NoAck、OnlyAck、RigorAck 三种模式）
//   - 心跳保活（Ping/Pong 机制）
//   - 并发任务处理（基于 TaskRunner）
//
// 消息确认模式:
//   - NoAck: 无确认模式，消息发送后不等待客户端确认
//   - OnlyAck: 简单确认模式，服务端收到消息后立即回复 ACK
//   - RigorAck: 严格确认模式，需要客户端二次确认，支持超时重发
//
// 连接管理:
//   - connToUser: 连接到用户 ID 的映射
//   - userToConn: 用户 ID 到连接的映射
//   - 同一用户多次登录时，新连接会关闭旧连接
//
// 消息处理流程:
//   1. 客户端建立 WebSocket 连接
//   2. 服务端进行身份认证
//   3. 记录连接映射关系
//   4. 启动读写协程处理消息
//   5. 根据消息类型和 method 路由到对应处理器
//
// =============================================================================
package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/threading"
	"time"

	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

// AckType 消息确认类型
// 定义了三种消息确认模式，用于保证消息的可靠传输
type AckType int

const (
	// NoAck 无确认模式
	// 消息发送后不等待客户端确认，适用于对消息可靠性要求不高的场景
	NoAck AckType = iota

	// OnlyAck 简单确认模式
	// 服务端收到消息后立即回复 ACK，然后进行业务处理
	// 适用于需要基本可靠性保证的场景
	OnlyAck

	// RigorAck 严格确认模式
	// 需要客户端二次确认，支持超时重发机制
	// 适用于对消息可靠性要求极高的场景（如支付、订单等）
	RigorAck
)

// ToString 将确认类型转换为字符串表示
// 用于日志输出和调试
//
// 返回:
//   - string: 确认类型的字符串表示
func (t AckType) ToString() string {
	switch t {
	case OnlyAck:
		return "OnlyAck"
	case RigorAck:
		return "RigorAck"
	}

	return "NoAck"
}

// Server WebSocket 服务器
// 管理所有 WebSocket 连接，提供消息路由、用户管理、消息确认等功能
type Server struct {
	sync.RWMutex // 读写锁，保护并发访问 connToUser 和 userToConn

	*threading.TaskRunner // 任务执行器，用于并发处理消息

	opt            *serverOption  // 服务器配置选项
	authentication Authentication // 身份认证接口

	routes map[string]HandlerFunc // 消息路由表，key 为 method，value 为处理函数
	addr   string                 // 服务器监听地址，格式: "host:port"
	patten string                 // WebSocket 路径模式，如 "/ws"

	connToUser map[*Conn]string // 连接到用户 ID 的映射，用于根据连接查找用户
	userToConn map[string]*Conn // 用户 ID 到连接的映射，用于根据用户查找连接

	upgrader websocket.Upgrader // WebSocket 升级器，用于将 HTTP 连接升级为 WebSocket
	logx.Logger                 // 日志记录器
}

// NewServer 创建 WebSocket 服务器实例
// 初始化服务器的所有组件，包括路由表、连接映射、升级器等
//
// 参数:
//   - addr: 服务器监听地址，格式为 "host:port"，如 "0.0.0.0:8080"
//   - opts: 可变参数，用于配置服务器选项（认证、ACK 模式、并发数等）
//
// 返回:
//   - *Server: 初始化完成的服务器实例
func NewServer(addr string, opts ...ServerOptions) *Server {
	// 合并所有配置选项
	opt := newServerOptions(opts...)

	return &Server{
		routes:   make(map[string]HandlerFunc), // 初始化路由表
		addr:     addr,
		patten:   opt.patten,
		opt:      &opt,
		upgrader: websocket.Upgrader{}, // 使用默认的 WebSocket 升级器

		authentication: opt.Authentication,

		connToUser: make(map[*Conn]string), // 初始化连接到用户的映射
		userToConn: make(map[string]*Conn), // 初始化用户到连接的映射

		Logger:     logx.WithContext(context.Background()),
		TaskRunner: threading.NewTaskRunner(opt.concurrency), // 创建任务执行器
	}
}

// ServerWs WebSocket 连接处理入口
// 处理 HTTP 请求升级为 WebSocket 连接的整个流程
//
// 处理流程:
//  1. 将 HTTP 连接升级为 WebSocket 连接
//  2. 进行身份认证，验证用户访问权限
//  3. 记录连接映射关系（用户 ID 与连接的双向映射）
//  4. 启动连接处理协程，开始消息收发
//
// 参数:
//   - w: HTTP 响应写入器
//   - r: HTTP 请求对象，包含升级请求和认证信息
func (s *Server) ServerWs(w http.ResponseWriter, r *http.Request) {
	// 捕获 panic，防止单个连接异常导致整个服务崩溃
	defer func() {
		if r := recover(); r != nil {
			s.Errorf("server handler ws recover err %v", r)
		}
	}()

	// 创建新的 WebSocket 连接对象
	conn := NewConn(s, w, r)
	if conn == nil {
		return
	}
	//conn, err := s.upgrader.Upgrade(w, r, nil)
	//if err != nil {
	//	s.Errorf("upgrade err %v", err)
	//	return
	//}

	// 进行身份认证
	if !s.authentication.Auth(w, r) {
		//conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprint("不具备访问权限")))
		// 认证失败，发送错误消息并关闭连接
		s.Send(&Message{FrameType: FrameData, Data: fmt.Sprint("不具备访问权限")}, conn)
		conn.Close()
		return
	}

	// 记录连接映射关系
	s.addConn(conn, r)

	// 启动协程处理该连接的所有消息
	go s.handlerConn(conn)
}

// handlerConn 处理单个 WebSocket 连接的所有消息
// 为每个连接启动读写协程，实现全双工通信
//
// 处理流程:
//  1. 获取连接对应的用户 ID
//  2. 启动写协程处理消息发送和路由分发
//  3. 如果启用了 ACK 机制，启动 ACK 读取协程
//  4. 在主循环中持续读取客户端消息
//  5. 根据 ACK 模式决定消息处理方式
//
// 参数:
//   - conn: WebSocket 连接对象
func (s *Server) handlerConn(conn *Conn) {

	// 获取连接对应的用户 ID 列表，取第一个作为当前用户
	uids := s.GetUsers(conn)
	conn.Uid = uids[0]

	// 启动写协程，处理消息发送和路由分发
	go s.handlerWrite(conn)

	// 如果启用了 ACK 机制，启动 ACK 读取协程
	if s.isAck(nil) {
		go s.readAck(conn)
	}

	// 主循环：持续读取客户端消息
	for {
		// 从 WebSocket 连接读取消息
		_, msg, err := conn.ReadMessage()
		if err != nil {
			s.Errorf("websocket conn read message err %v", err)
			s.Close(conn)
			return
		}
		// 解析 JSON 消息
		var message Message
		if err = json.Unmarshal(msg, &message); err != nil {
			s.Errorf("json unmarshal err %v, msg %v", err, string(msg))
			s.Close(conn)
			return
		}

		// 根据 ACK 模式处理消息
		if s.isAck(&message) {
			// ACK 模式：将消息放入待确认队列
			s.Infof("conn message read ack msg %v", message)
			conn.appendMsgMq(&message)
		} else {
			// 非 ACK 模式：直接发送到消息通道
			conn.message <- &message
		}
	}
}

// isAck 判断是否需要消息确认
// 根据服务器配置和消息类型决定是否启用 ACK 机制
//
// 参数:
//   - message: 消息对象，如果为 nil 则只检查服务器配置
//
// 返回:
//   - bool: true 表示需要 ACK，false 表示不需要
func (s *Server) isAck(message *Message) bool {
	if message == nil {
		// 只检查服务器是否启用了 ACK
		return s.opt.ack != NoAck
	}
	// 检查服务器配置和消息的 FrameType
	// FrameNoAck 类型的消息即使服务器启用了 ACK 也不需要确认
	return s.opt.ack != NoAck && message.FrameType != FrameNoAck
}

// readAck 读取并处理消息确认
// 从连接的待确认队列中读取消息，根据 ACK 模式进行不同的确认处理
//
// ACK 处理模式:
//  1. OnlyAck: 立即回复 ACK，然后将消息发送到业务处理通道
//  2. RigorAck: 严格两阶段确认
//     - 第一阶段：服务端发送 ACK，等待客户端确认
//     - 第二阶段：客户端确认后才将消息发送到业务处理通道
//     - 支持超时重发机制
//
// 参数:
//   - conn: WebSocket 连接对象
func (s *Server) readAck(conn *Conn) {
	for {
		select {
		case <-conn.done:
			// 连接已关闭，退出 ACK 处理循环
			s.Infof("close message ack uid %v ", conn.Uid)
			return
		default:
		}

		// 从待确认队列中读取消息
		conn.messageMu.Lock()
		if len(conn.readMessage) == 0 {
			conn.messageMu.Unlock()
			// 队列为空，短暂休眠后继续
			time.Sleep(100 * time.Microsecond)
			continue
		}

		// 获取队列中的第一条消息
		message := conn.readMessage[0]

		// 根据 ACK 模式进行不同的处理
		switch s.opt.ack {
		case OnlyAck:
			// 简单确认模式：立即回复 ACK
			s.Send(&Message{
				FrameType: FrameAck,
				Id:        message.Id,
				AckSeq:    message.AckSeq + 1,
			}, conn)
			// 从队列中移除该消息
			conn.readMessage = conn.readMessage[1:]
			conn.messageMu.Unlock()

			// 将消息发送到业务处理通道
			conn.message <- message
		case RigorAck:
			// 严格确认模式：需要客户端二次确认
			if message.AckSeq == 0 {
				// 第一次确认：发送 ACK 给客户端
				conn.readMessage[0].AckSeq++
				conn.readMessage[0].ackTime = time.Now()
				s.Send(&Message{
					FrameType: FrameAck,
					Id:        message.Id,
					AckSeq:    message.AckSeq,
				}, conn)
				s.Infof("message ack RigorAck send mid %v, seq %v , time%v", message.Id, message.AckSeq,
					message.ackTime)
				conn.messageMu.Unlock()
				continue
			}

			// 第二次确认：验证客户端是否已确认

			// 1. 检查客户端返回的确认序号
			msgSeq := conn.readMessageSeq[message.Id]
			if msgSeq.AckSeq > message.AckSeq {
				// 客户端已确认，从队列中移除并发送到业务处理通道
				conn.readMessage = conn.readMessage[1:]
				conn.messageMu.Unlock()
				conn.message <- message
				s.Infof("message ack RigorAck success mid %v", message.Id)
				continue
			}

			// 2. 客户端未确认，检查是否超时
			val := s.opt.ackTimeout - time.Since(message.ackTime)
			if !message.ackTime.IsZero() && val <= 0 {
				// 超时，放弃该消息
				delete(conn.readMessageSeq, message.Id)
				conn.readMessage = conn.readMessage[1:]
				conn.messageMu.Unlock()
				continue
			}
			// 未超时，重新发送 ACK
			conn.messageMu.Unlock()
			s.Send(&Message{
				FrameType: FrameAck,
				Id:        message.Id,
				AckSeq:    message.AckSeq,
			}, conn)
			// 等待一段时间后重试
			time.Sleep(3 * time.Second)
		}
	}
}

// handlerWrite 处理消息写入和路由分发
// 从连接的消息通道中读取消息，根据消息类型进行不同的处理
//
// 消息类型处理:
//  - FramePing: 心跳消息，回复 Pong
//  - FrameData: 业务消息，根据 method 路由到对应的处理器
//
// 参数:
//   - conn: WebSocket 连接对象
func (s *Server) handlerWrite(conn *Conn) {
	for {
		select {
		case <-conn.done:
			// 连接已关闭，退出写入循环
			return
		case message := <-conn.message:
			// 根据消息帧类型进行处理
			switch message.FrameType {
			case FramePing:
				// 心跳消息，回复 Pong 保持连接活跃
				s.Send(&Message{FrameType: FramePing}, conn)
			case FrameData:
				// 业务数据消息，根据 method 路由到对应的处理器
				if handler, ok := s.routes[message.Method]; ok {
					handler(s, conn, message)
				} else {
					// 路由不存在，返回错误消息
					s.Send(&Message{FrameType: FrameData, Data: fmt.Sprintf("不存在执行的方法 %v 请检查", message.Method)}, conn)
					//conn.WriteMessage(&Message{}, []byte(fmt.Sprintf("不存在执行的方法 %v 请检查", message.Method)))
				}
			}

			// 如果启用了 ACK，清理已确认的消息
			if s.isAck(message) {
				conn.messageMu.Lock()
				delete(conn.readMessageSeq, message.Id)
				conn.messageMu.Unlock()
			}
		}
	}
}

// addConn 添加连接到服务器的连接映射表
// 建立用户 ID 与连接的双向映射关系
// 如果用户已有连接，会先关闭旧连接（实现单点登录）
//
// 参数:
//   - conn: WebSocket 连接对象
//   - req: HTTP 请求对象，用于提取用户 ID
func (s *Server) addConn(conn *Conn, req *http.Request) {
	// 从请求中提取用户 ID
	uid := s.authentication.UserId(req)

	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	// 检查用户是否已有连接（实现单点登录）
	if c := s.userToConn[uid]; c != nil {
		// 关闭旧连接，新连接会替代旧连接
		c.Close()
	}

	// 建立双向映射
	s.connToUser[conn] = uid
	s.userToConn[uid] = conn
}

// GetConn 根据用户 ID 获取对应的连接
//
// 参数:
//   - uid: 用户 ID
//
// 返回:
//   - *Conn: 用户对应的连接，如果用户未连接则返回 nil
func (s *Server) GetConn(uid string) *Conn {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	fmt.Println(s.userToConn)
	return s.userToConn[uid]
}

// GetConns 根据用户 ID 列表批量获取连接
//
// 参数:
//   - uids: 用户 ID 列表
//
// 返回:
//   - []*Conn: 连接列表，顺序与 uids 对应
func (s *Server) GetConns(uids ...string) []*Conn {
	if len(uids) == 0 {
		return nil
	}

	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	res := make([]*Conn, 0, len(uids))
	for _, uid := range uids {
		res = append(res, s.userToConn[uid])
	}
	return res
}

// GetUsers 根据连接列表获取对应的用户 ID 列表
// 如果不传入连接参数，则返回所有在线用户的 ID
//
// 参数:
//   - conns: 连接列表，如果为空则返回所有在线用户
//
// 返回:
//   - []string: 用户 ID 列表
func (s *Server) GetUsers(conns ...*Conn) []string {

	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	var res []string
	if len(conns) == 0 {
		// 获取所有在线用户
		res = make([]string, 0, len(s.connToUser))
		for _, uid := range s.connToUser {
			res = append(res, uid)
		}
	} else {
		// 获取指定连接对应的用户
		res = make([]string, 0, len(conns))
		for _, conn := range conns {
			res = append(res, s.connToUser[conn])
		}
	}

	return res
}

// Close 关闭连接并清理映射关系
// 从服务器的连接映射表中移除该连接
//
// 参数:
//   - conn: 要关闭的连接对象
func (s *Server) Close(conn *Conn) {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	// 获取连接对应的用户 ID
	uid := s.connToUser[conn]
	if uid == "" {
		// 连接已被关闭
		return
	}

	// 清理双向映射
	delete(s.connToUser, conn)
	delete(s.userToConn, uid)

	// 关闭底层 WebSocket 连接
	conn.Close()
}

// SendByUserId 根据用户 ID 列表发送消息
// 先根据用户 ID 获取对应的连接，然后发送消息
//
// 参数:
//   - msg: 要发送的消息对象（会被序列化为 JSON）
//   - sendIds: 目标用户 ID 列表
//
// 返回:
//   - error: 发送失败时返回错误
func (s *Server) SendByUserId(msg interface{}, sendIds ...string) error {
	if len(sendIds) == 0 {
		return nil
	}

	return s.Send(msg, s.GetConns(sendIds...)...)
}

// Send 向指定的连接列表发送消息
// 将消息对象序列化为 JSON 后通过 WebSocket 发送
//
// 参数:
//   - msg: 要发送的消息对象（会被序列化为 JSON）
//   - conns: 目标连接列表
//
// 返回:
//   - error: 发送失败时返回错误
func (s *Server) Send(msg interface{}, conns ...*Conn) error {
	if len(conns) == 0 {
		return nil
	}

	// 将消息序列化为 JSON
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	// 向所有目标连接发送消息
	for _, conn := range conns {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			return err
		}
	}

	return nil
}

// AddRoutes 批量添加消息路由
// 将 method 与对应的处理函数注册到路由表
//
// 参数:
//   - rs: 路由列表，每个路由包含 method 和对应的处理函数
func (s *Server) AddRoutes(rs []Route) {
	for _, r := range rs {
		s.routes[r.Method] = r.Handler
	}
}

// Start 启动 WebSocket 服务器
// 注册 HTTP 处理函数并开始监听
func (s *Server) Start() {
	http.HandleFunc(s.patten, s.ServerWs)
	s.Info(http.ListenAndServe(s.addr, nil))
}

// Stop 停止 WebSocket 服务器
// 执行清理工作并关闭服务
func (s *Server) Stop() {
	fmt.Println("停止服务")
}
