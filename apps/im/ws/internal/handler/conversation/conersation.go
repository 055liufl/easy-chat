// =============================================================================
// 会话处理器 - 聊天消息和已读状态处理
// =============================================================================
// 处理会话相关的 WebSocket 消息，包括:
//   - 聊天消息发送（单聊/群聊）
//   - 消息已读标记
//   - 会话 ID 自动生成
//
// 业务流程:
//   聊天消息:
//     1. 客户端发送聊天消息
//     2. 服务端生成会话 ID（如果未提供）
//     3. 将消息推送到 Kafka 消息队列
//     4. 消息处理服务消费 Kafka 消息，进行持久化和推送
//
//   已读标记:
//     1. 客户端发送已读标记请求
//     2. 服务端将已读信息推送到 Kafka 消息队列
//     3. 消息处理服务更新消息已读状态
//
// =============================================================================

package conversation

import (
	"github.com/mitchellh/mapstructure"
	"imooc.com/easy-chat/apps/im/ws/internal/svc"
	"imooc.com/easy-chat/apps/im/ws/websocket"
	"imooc.com/easy-chat/apps/im/ws/ws"
	"imooc.com/easy-chat/apps/task/mq/mq"
	"imooc.com/easy-chat/pkg/constants"
	"imooc.com/easy-chat/pkg/wuid"
	"time"
)

// Chat 处理聊天消息发送
// 支持单聊和群聊消息，将消息推送到 Kafka 队列进行异步处理
//
// 业务流程:
//  1. 解析客户端发送的聊天消息数据
//  2. 根据聊天类型生成会话 ID（如果未提供）
//     - 单聊: 使用发送者和接收者 ID 组合生成唯一会话 ID
//     - 群聊: 使用群组 ID 作为会话 ID
//  3. 将消息推送到 Kafka 消息队列
//  4. 消息处理服务消费 Kafka 消息，进行持久化存储和推送
//
// 消息格式:
//   请求: {
//     "method": "conversation.chat",
//     "data": {
//       "conversationId": "会话ID（可选）",
//       "chatType": 1,  // 1-单聊, 2-群聊
//       "recvId": "接收者ID",
//       "msg": {
//         "mType": 1,  // 消息类型: 1-文本, 2-图片, 3-语音等
//         "content": "消息内容"
//       }
//     }
//   }
//
// 参数:
//   - svc: 服务上下文，包含 Kafka 客户端等依赖组件
//
// 返回:
//   - websocket.HandlerFunc: WebSocket 消息处理函数
func Chat(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		// 解析聊天消息数据
		var data ws.Chat
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			srv.Send(websocket.NewErrMessage(err), conn)
			return
		}

		// 如果客户端未提供会话 ID，根据聊天类型自动生成
		if data.ConversationId == "" {
			switch data.ChatType {
			case constants.SingleChatType:
				// 单聊: 使用发送者和接收者 ID 组合生成唯一会话 ID
				data.ConversationId = wuid.CombineId(conn.Uid, data.RecvId)
			case constants.GroupChatType:
				// 群聊: 使用群组 ID 作为会话 ID
				data.ConversationId = data.RecvId
			}
		}

		// 将消息推送到 Kafka 消息队列，由消息处理服务异步消费
		err := svc.MsgChatTransferClient.Push(&mq.MsgChatTransfer{
			ConversationId: data.ConversationId,
			ChatType:       data.ChatType,
			SendId:         conn.Uid,           // 发送者 ID（从连接上下文获取）
			RecvId:         data.RecvId,
			SendTime:       time.Now().UnixNano(), // 发送时间（纳秒时间戳）
			MType:          data.Msg.MType,
			Content:        data.Msg.Content,
		})
		if err != nil {
			srv.Send(websocket.NewErrMessage(err), conn)
			return
		}
	}
}

// MarkRead 处理消息已读标记
// 将消息已读状态推送到 Kafka 队列，由消息处理服务更新已读状态
//
// 业务流程:
//  1. 解析客户端发送的已读标记数据
//  2. 将已读信息推送到 Kafka 消息队列
//  3. 消息处理服务消费 Kafka 消息，更新消息已读状态
//  4. 通知发送者消息已被阅读（可选）
//
// 消息格式:
//   请求: {
//     "method": "conversation.markChat",
//     "data": {
//       "chatType": 1,  // 1-单聊, 2-群聊
//       "conversationId": "会话ID",
//       "recvId": "接收者ID",
//       "msgIds": ["消息ID1", "消息ID2"]  // 需要标记为已读的消息 ID 列表
//     }
//   }
//
// 参数:
//   - svc: 服务上下文，包含 Kafka 客户端等依赖组件
//
// 返回:
//   - websocket.HandlerFunc: WebSocket 消息处理函数
func MarkRead(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		// 解析已读标记数据
		var data ws.MarkRead
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			srv.Send(websocket.NewErrMessage(err), conn)
			return
		}

		// 将已读信息推送到 Kafka 消息队列，由消息处理服务异步消费
		err := svc.MsgReadTransferClient.Push(&mq.MsgMarkRead{
			ChatType:       data.ChatType,
			ConversationId: data.ConversationId,
			SendId:         conn.Uid,  // 标记已读的用户 ID（从连接上下文获取）
			RecvId:         data.RecvId,
			MsgIds:         data.MsgIds,
		})

		if err != nil {
			srv.Send(websocket.NewErrMessage(err), conn)
			return
		}
	}
}
