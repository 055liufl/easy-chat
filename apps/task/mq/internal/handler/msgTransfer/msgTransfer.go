// =============================================================================
// 消息传输基础模块 - 消息推送核心逻辑
// =============================================================================
// 本模块提供消息传输的基础功能，负责将消息推送给接收方，包括：
//   - 单聊消息推送：直接推送给指定用户
//   - 群聊消息推送：查询群成员后批量推送
//
// 数据来源:
//   - Kafka 消息队列（聊天消息、已读消息）
//   - Social RPC 服务（群成员信息）
//
// 业务场景:
//   - 聊天消息持久化后推送给接收方
//   - 已读消息更新后推送给发送方
//   - 群聊消息需要查询群成员列表后批量推送
//
// 推送流程:
//   1. 根据聊天类型（单聊/群聊）选择推送策略
//   2. 群聊消息需要先查询群成员列表
//   3. 通过 WebSocket 客户端推送消息
//
// =============================================================================

package msgTransfer

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"imooc.com/easy-chat/apps/im/ws/websocket"
	"imooc.com/easy-chat/apps/im/ws/ws"
	"imooc.com/easy-chat/apps/social/rpc/socialclient"
	"imooc.com/easy-chat/apps/task/mq/internal/svc"
	"imooc.com/easy-chat/pkg/constants"
)

// baseMsgTransfer 消息传输基础结构
// 提供消息推送的通用功能，被具体的消息处理器继承使用
type baseMsgTransfer struct {
	svcCtx *svc.ServiceContext // 服务上下文，包含配置、数据库、RPC 客户端等
	logx.Logger                // 日志记录器
}

// NewBaseMsgTransfer 创建消息传输基础实例
//
// 参数:
//   - svc: 服务上下文，包含 WebSocket 客户端、Social RPC 客户端等
//
// 返回:
//   - *baseMsgTransfer: 消息传输基础实例
func NewBaseMsgTransfer(svc *svc.ServiceContext) *baseMsgTransfer {
	return &baseMsgTransfer{
		svcCtx: svc,
		Logger: logx.WithContext(context.Background()),
	}
}

// Transfer 消息传输入口
// 根据聊天类型（单聊/群聊）选择对应的推送策略
//
// 参数:
//   - ctx: 上下文
//   - data: 推送数据，包含消息内容、接收方等信息
//
// 返回:
//   - error: 推送失败时返回错误
func (m *baseMsgTransfer) Transfer(ctx context.Context, data *ws.Push) error {
	var err error
	switch data.ChatType {
	case constants.GroupChatType:
		// 群聊消息：需要查询群成员后批量推送
		err = m.group(ctx, data)
	case constants.SingleChatType:
		// 单聊消息：直接推送给指定用户
		err = m.single(ctx, data)
	}
	return err
}

// single 单聊消息推送
// 直接通过 WebSocket 推送消息给指定接收方
//
// 参数:
//   - ctx: 上下文
//   - data: 推送数据，RecvId 字段指定接收方用户 ID
//
// 返回:
//   - error: 推送失败时返回错误
func (m *baseMsgTransfer) single(ctx context.Context, data *ws.Push) error {
	return m.svcCtx.WsClient.Send(websocket.Message{
		FrameType: websocket.FrameData,
		Method:    "push",
		FormId:    constants.SYSTEM_ROOT_UID, // 系统消息，来源为系统根用户
		Data:      data,
	})
}

// group 群聊消息推送
// 先查询群成员列表，然后批量推送消息（排除发送者自己）
//
// 业务逻辑:
//   1. 调用 Social RPC 查询群成员列表
//   2. 过滤掉发送者自己（发送者已经在客户端显示消息）
//   3. 将其他成员 ID 填充到 RecvIds 字段
//   4. 通过 WebSocket 批量推送
//
// 参数:
//   - ctx: 上下文
//   - data: 推送数据，RecvId 字段为群 ID，SendId 为发送者 ID
//
// 返回:
//   - error: 查询群成员或推送失败时返回错误
func (m *baseMsgTransfer) group(ctx context.Context, data *ws.Push) error {
	// 查询群成员列表
	users, err := m.svcCtx.Social.GroupUsers(ctx, &socialclient.GroupUsersReq{
		GroupId: data.RecvId,
	})
	if err != nil {
		return err
	}
	data.RecvIds = make([]string, 0, len(users.List))

	// 遍历群成员，排除发送者自己
	for _, members := range users.List {
		if members.UserId == data.SendId {
			continue
		}

		data.RecvIds = append(data.RecvIds, members.UserId)
	}

	// 批量推送给群成员
	return m.svcCtx.WsClient.Send(websocket.Message{
		FrameType: websocket.FrameData,
		Method:    "push",
		FormId:    constants.SYSTEM_ROOT_UID, // 系统消息，来源为系统根用户
		Data:      data,
	})
}
