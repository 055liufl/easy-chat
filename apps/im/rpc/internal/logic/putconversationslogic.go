// =============================================================================
// 更新用户会话业务逻辑
// =============================================================================
// 提供会话更新功能，包括:
//   - 更新会话的已读消息数
//   - 更新会话的显示状态
//   - 更新会话的消息序列号
//
// 数据操作:
//   - Conversations 集合：更新用户的会话列表
//
// 业务场景:
//   - 用户阅读消息后更新已读数
//   - 用户隐藏或显示会话
//   - 同步会话状态
//
// =============================================================================
package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"imooc.com/easy-chat/apps/im/immodels"
	"imooc.com/easy-chat/apps/im/rpc/im"
	"imooc.com/easy-chat/apps/im/rpc/internal/svc"
	"imooc.com/easy-chat/pkg/constants"
	"imooc.com/easy-chat/pkg/xerr"
)

// PutConversationsLogic 更新用户会话业务逻辑处理器
type PutConversationsLogic struct {
	ctx    context.Context        // 请求上下文，用于超时控制和链路追踪
	svcCtx *svc.ServiceContext    // 服务上下文，提供数据模型访问
	logx.Logger                   // 日志记录器
}

// NewPutConversationsLogic 创建更新用户会话业务逻辑实例
//
// 参数:
//   - ctx: 请求上下文
//   - svcCtx: 服务上下文，包含数据模型等依赖
//
// 返回:
//   - *PutConversationsLogic: 业务逻辑处理器实例
func NewPutConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PutConversationsLogic {
	return &PutConversationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// PutConversations 更新会话
// 批量更新用户的会话信息，包括已读数、显示状态等
//
// 参数:
//   - in: 更新请求
//     - UserId: 用户 ID
//     - ConversationList: 要更新的会话列表（map 结构）
//       - ConversationId: 会话 ID
//       - ChatType: 聊天类型
//       - IsShow: 是否显示
//       - Read: 本次已读消息数（会累加到总已读数）
//       - Seq: 消息序列号
//
// 返回:
//   - *im.PutConversationsResp: 更新响应
//   - error: 错误信息
//
// 业务流程:
//   1. 查询用户的会话列表
//   2. 初始化会话列表（如果不存在）
//   3. 遍历要更新的会话，累加已读数并更新其他字段
//   4. 保存更新后的会话列表到数据库
func (l *PutConversationsLogic) PutConversations(in *im.PutConversationsReq) (*im.PutConversationsResp, error) {
	// todo: add your logic here and delete this line

	data, err := l.svcCtx.ConversationsModel.FindByUserId(l.ctx, in.UserId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "ConversationsModel.FindByUserId err %v, req %v", err, in.UserId)
	}

	if data.ConversationList == nil {
		data.ConversationList = make(map[string]*immodels.Conversation)
	}

	for s, conversation := range in.ConversationList {
		var oldTotal int
		if data.ConversationList[s] != nil {
			oldTotal = data.ConversationList[s].Total
		}

		data.ConversationList[s] = &immodels.Conversation{
			ConversationId: conversation.ConversationId,
			ChatType:       constants.ChatType(conversation.ChatType),
			IsShow:         conversation.IsShow,
			Total:          int(conversation.Read) + oldTotal,
			Seq:            conversation.Seq,
		}
	}

	err = l.svcCtx.ConversationsModel.Update(l.ctx, data)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "ConversationsModel.Update err %v, req %v", err, data)
	}

	return &im.PutConversationsResp{}, nil
}
