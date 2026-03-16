// =============================================================================
// 消息推送处理器 - 服务端消息推送
// =============================================================================
// 处理服务端向客户端推送消息，支持:
//   - 单聊消息推送（点对点）
//   - 群聊消息推送（一对多）
//   - 离线消息处理
//
// 推送场景:
//   - 用户 A 发送消息后，服务端将消息推送给用户 B
//   - 群聊消息推送给群内所有成员
//   - 消息从 Kafka 消费后推送给在线用户
//
// 业务流程:
//  1. 接收推送请求（包含消息内容、接收者信息）
//  2. 根据聊天类型（单聊/群聊）选择推送策略
//  3. 检查接收者在线状态
//  4. 推送消息到在线用户的 WebSocket 连接
//  5. 离线用户的消息由其他服务处理（如离线推送、消息存储）
//
// =============================================================================

package push

import (
	"github.com/mitchellh/mapstructure"
	"imooc.com/easy-chat/apps/im/ws/internal/svc"
	"imooc.com/easy-chat/apps/im/ws/websocket"
	"imooc.com/easy-chat/apps/im/ws/ws"
	"imooc.com/easy-chat/pkg/constants"
)

// Push 处理消息推送请求
// 根据聊天类型（单聊/群聊）将消息推送给目标用户
//
// 业务流程:
//  1. 解析推送消息数据（包含发送者、接收者、消息内容等）
//  2. 根据 ChatType 判断聊天类型
//  3. 单聊: 推送给指定的单个用户
//  4. 群聊: 推送给群内所有成员
//
// 消息格式:
//   请求: {
//     "method": "push",
//     "data": {
//       "conversationId": "会话ID",
//       "chatType": 1,  // 1-单聊, 2-群聊
//       "sendId": "发送者ID",
//       "recvId": "接收者ID",  // 单聊时使用
//       "recvIds": ["ID1", "ID2"],  // 群聊时使用
//       "msgId": "消息ID",
//       "content": "消息内容",
//       "sendTime": 1234567890
//     }
//   }
//
// 参数:
//   - svc: 服务上下文，包含业务处理所需的依赖组件
//
// 返回:
//   - websocket.HandlerFunc: WebSocket 消息处理函数
func Push(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		var data ws.Push
		// 解析推送消息数据
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			srv.Send(websocket.NewErrMessage(err))
			return
		}

		// 根据聊天类型选择推送策略
		switch data.ChatType {
		case constants.SingleChatType:
			// 单聊: 推送给指定的单个用户
			single(srv, &data, data.RecvId)
		case constants.GroupChatType:
			// 群聊: 推送给群内所有成员
			group(srv, &data)
		}
	}
}

// single 单聊消息推送
// 将消息推送给指定的单个用户
//
// 业务流程:
//  1. 根据接收者 ID 获取其 WebSocket 连接
//  2. 如果用户在线，构造消息并推送
//  3. 如果用户离线，记录日志（离线消息由其他服务处理）
//
// 参数:
//   - srv: WebSocket 服务器实例
//   - data: 推送消息数据
//   - recvId: 接收者用户 ID
//
// 返回:
//   - error: 推送失败时返回错误信息
func single(srv *websocket.Server, data *ws.Push, recvId string) error {
	// 根据用户 ID 获取 WebSocket 连接
	rconn := srv.GetConn(recvId)
	if rconn == nil {
		// 用户离线，离线消息由其他服务处理（如离线推送、APNs、FCM 等）
		// todo: 目标离线
		return nil
	}

	srv.Infof("push msg %v", data)

	// 构造聊天消息并推送给接收者
	return srv.Send(websocket.NewMessage(data.SendId, &ws.Chat{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendTime:       data.SendTime,
		Msg: ws.Msg{
			ReadRecords: data.ReadRecords,
			MsgId:       data.MsgId,
			MType:       data.MType,
			Content:     data.Content,
		},
	}), rconn)
}

// group 群聊消息推送
// 将消息推送给群内所有成员
//
// 业务流程:
//  1. 遍历所有接收者 ID 列表
//  2. 为每个接收者创建独立的推送任务
//  3. 使用协程池异步执行推送任务，避免阻塞
//
// 参数:
//   - srv: WebSocket 服务器实例
//   - data: 推送消息数据，包含接收者 ID 列表
//
// 返回:
//   - error: 推送失败时返回错误信息
func group(srv *websocket.Server, data *ws.Push) error {
	// 遍历所有接收者
	for _, id := range data.RecvIds {
		func(id string) {
			// 使用协程池异步推送，避免阻塞主流程
			srv.Schedule(func() {
				single(srv, data, id)
			})
		}(id)
	}
	return nil
}
