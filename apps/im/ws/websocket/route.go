// =============================================================================
// WebSocket Route - WebSocket 路由定义
// =============================================================================
// 定义 WebSocket 消息的路由结构
// 类似于 HTTP 路由，根据消息的 Method 字段将消息分发到对应的处理函数
//
// 路由机制:
//   - 每个路由包含一个 Method（方法名）和对应的 Handler（处理函数）
//   - 服务器收到消息后，根据 Message.Method 查找对应的 Handler
//   - Handler 接收服务器、连接和消息三个参数，可以进行业务处理
//
// 使用示例:
//   routes := []Route{
//       {Method: "chat.send", Handler: handleChatSend},
//       {Method: "user.online", Handler: handleUserOnline},
//   }
//   server.AddRoutes(routes)
//
// =============================================================================
package websocket

// Route WebSocket 路由
// 定义了消息方法名与处理函数的映射关系
type Route struct {
	Method  string      // 方法名，用于标识消息类型（如 "chat.send", "user.online"）
	Handler HandlerFunc // 处理函数，用于处理该类型的消息
}

// HandlerFunc 路由处理函数类型
// 定义了消息处理函数的签名
//
// 参数:
//   - srv: 服务器实例，可用于发送消息、管理连接等
//   - conn: 当前连接对象，表示消息来源
//   - msg: 接收到的消息对象，包含消息数据和元信息
type HandlerFunc func(srv *Server, conn *Conn, msg *Message)
