// =============================================================================
// 用户处理器 - 用户在线状态管理
// =============================================================================
// 处理用户相关的 WebSocket 消息，包括:
//   - 用户上线通知
//   - 在线用户列表查询
//   - 用户状态广播
//
// 业务场景:
//   - 用户建立 WebSocket 连接后，发送上线消息
//   - 服务端返回当前所有在线用户列表
//   - 可用于实现在线状态显示、好友在线提醒等功能
//
// =============================================================================

package user

import (
	"imooc.com/easy-chat/apps/im/ws/internal/svc"
	"imooc.com/easy-chat/apps/im/ws/websocket"
)

// OnLine 处理用户上线请求
// 返回当前所有在线用户的 ID 列表
//
// 业务流程:
//  1. 获取服务器上所有在线用户的 ID 列表
//  2. 获取当前连接用户的 ID
//  3. 将在线用户列表发送给当前用户
//
// 消息格式:
//   请求: {"method": "user.online", "data": {}}
//   响应: {"sendId": "当前用户ID", "data": ["用户ID1", "用户ID2", ...]}
//
// 参数:
//   - svc: 服务上下文，包含业务处理所需的依赖组件
//
// 返回:
//   - websocket.HandlerFunc: WebSocket 消息处理函数
func OnLine(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		// 获取服务器上所有在线用户的 ID 列表
		uids := srv.GetUsers()

		// 获取当前连接用户的 ID（返回数组，取第一个元素）
		u := srv.GetUsers(conn)

		// 将在线用户列表发送给当前用户
		err := srv.Send(websocket.NewMessage(u[0], uids), conn)
		srv.Info("err ", err)
	}
}
