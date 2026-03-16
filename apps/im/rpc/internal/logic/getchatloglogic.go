// =============================================================================
// 获取聊天记录业务逻辑
// =============================================================================
// 提供聊天记录查询功能，支持:
//   - 根据消息 ID 查询单条聊天记录
//   - 根据会话 ID 和时间范围查询聊天记录列表
//
// 数据来源:
//   从 MongoDB 的 ChatLog 集合中查询
//
// 业务场景:
//   - 用户查看历史聊天记录
//   - 消息详情查看
//   - 聊天记录分页加载
//
// =============================================================================
package logic

import (
	"context"
	"github.com/pkg/errors"
	"imooc.com/easy-chat/pkg/xerr"

	"imooc.com/easy-chat/apps/im/rpc/im"
	"imooc.com/easy-chat/apps/im/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

// GetChatLogLogic 获取聊天记录业务逻辑处理器
type GetChatLogLogic struct {
	ctx    context.Context        // 请求上下文，用于超时控制和链路追踪
	svcCtx *svc.ServiceContext    // 服务上下文，提供数据模型访问
	logx.Logger                   // 日志记录器
}

// NewGetChatLogLogic 创建获取聊天记录业务逻辑实例
//
// 参数:
//   - ctx: 请求上下文
//   - svcCtx: 服务上下文，包含数据模型等依赖
//
// 返回:
//   - *GetChatLogLogic: 业务逻辑处理器实例
func NewGetChatLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatLogLogic {
	return &GetChatLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetChatLog 获取会话记录
// 支持两种查询方式:
//   1. 根据消息 ID 查询单条记录
//   2. 根据会话 ID 和时间范围查询记录列表
//
// 参数:
//   - in: 查询请求
//     - MsgId: 消息 ID（可选），如果提供则只查询该条消息
//     - ConversationId: 会话 ID，用于查询会话下的消息列表
//     - StartSendTime: 开始时间（时间戳），查询该时间之后的消息
//     - EndSendTime: 结束时间（时间戳），查询该时间之前的消息
//     - Count: 查询数量限制
//
// 返回:
//   - *im.GetChatLogResp: 聊天记录列表
//   - error: 错误信息
//
// 业务流程:
//   1. 如果提供了 MsgId，直接查询单条记录并返回
//   2. 否则根据会话 ID 和时间范围查询记录列表
//   3. 将数据库模型转换为 RPC 响应格式
func (l *GetChatLogLogic) GetChatLog(in *im.GetChatLogReq) (*im.GetChatLogResp, error) {
	// todo: add your logic here and delete this line

	// 根据id
	if in.MsgId != "" {
		chatlog, err := l.svcCtx.ChatLogModel.FindOne(l.ctx, in.MsgId)
		if err != nil {
			return nil, errors.Wrapf(xerr.NewDBErr(), "find chatLog by msgId err %v, req %v", err, in.MsgId)
		}

		return &im.GetChatLogResp{
			List: []*im.ChatLog{{
				Id:             chatlog.ID.Hex(),
				ConversationId: chatlog.ConversationId,
				SendId:         chatlog.SendId,
				RecvId:         chatlog.RecvId,
				MsgType:        int32(chatlog.MsgType),
				MsgContent:     chatlog.MsgContent,
				ChatType:       int32(chatlog.ChatType),
				SendTime:       chatlog.SendTime,
				ReadRecords:    chatlog.ReadRecords,
			}},
		}, nil
	}

	data, err := l.svcCtx.ChatLogModel.ListBySendTime(l.ctx, in.ConversationId, in.StartSendTime, in.EndSendTime, in.Count)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find chatLog list by SendTime err %v, req %v", err, in)
	}

	res := make([]*im.ChatLog, 0, len(data))
	for _, datum := range data {
		res = append(res, &im.ChatLog{
			Id:             datum.ID.Hex(),
			ConversationId: datum.ConversationId,
			SendId:         datum.SendId,
			RecvId:         datum.RecvId,
			MsgType:        int32(datum.MsgType),
			MsgContent:     datum.MsgContent,
			ChatType:       int32(datum.ChatType),
			SendTime:       datum.SendTime,
			ReadRecords:    datum.ReadRecords,
		})
	}

	return &im.GetChatLogResp{
		List: res,
	}, nil
}
