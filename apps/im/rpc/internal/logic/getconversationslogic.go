// =============================================================================
// 获取用户会话列表业务逻辑
// =============================================================================
// 提供用户会话列表查询功能，包括:
//   - 查询用户的所有会话
//   - 计算每个会话的未读消息数
//   - 更新会话显示状态
//
// 数据来源:
//   - Conversations 集合：用户的会话列表
//   - Conversation 集合：会话的详细信息（消息总数等）
//
// 业务场景:
//   - 用户打开聊天列表页面
//   - 刷新会话列表
//   - 显示未读消息数量
//
// =============================================================================
package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"imooc.com/easy-chat/apps/im/immodels"
	"imooc.com/easy-chat/pkg/xerr"

	"imooc.com/easy-chat/apps/im/rpc/im"
	"imooc.com/easy-chat/apps/im/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

// GetConversationsLogic 获取用户会话列表业务逻辑处理器
type GetConversationsLogic struct {
	ctx    context.Context        // 请求上下文，用于超时控制和链路追踪
	svcCtx *svc.ServiceContext    // 服务上下文，提供数据模型访问
	logx.Logger                   // 日志记录器
}

// NewGetConversationsLogic 创建获取用户会话列表业务逻辑实例
//
// 参数:
//   - ctx: 请求上下文
//   - svcCtx: 服务上下文，包含数据模型等依赖
//
// 返回:
//   - *GetConversationsLogic: 业务逻辑处理器实例
func NewGetConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationsLogic {
	return &GetConversationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetConversations 获取会话列表
// 查询用户的所有会话，并计算未读消息数
//
// 参数:
//   - in: 查询请求
//     - UserId: 用户 ID
//
// 返回:
//   - *im.GetConversationsResp: 会话列表响应
//     - ConversationList: 会话列表（map 结构，key 为会话 ID）
//       - Total: 会话总消息数
//       - ToRead: 未读消息数
//       - IsShow: 是否显示该会话
//   - error: 错误信息
//
// 业务流程:
//   1. 根据用户 ID 查询用户的会话列表
//   2. 提取所有会话 ID，批量查询会话详情
//   3. 对比用户已读消息数和会话总消息数，计算未读数
//   4. 如果有未读消息，更新会话为显示状态
func (l *GetConversationsLogic) GetConversations(in *im.GetConversationsReq) (*im.GetConversationsResp, error) {
	// todo: add your logic here and delete this line

	// 根据用户查询用户会话列表
	data, err := l.svcCtx.ConversationsModel.FindByUserId(l.ctx, in.UserId)
	if err != nil {
		if err == immodels.ErrNotFound {
			return &im.GetConversationsResp{}, nil
		}
		return nil, errors.Wrapf(xerr.NewDBErr(), "ConversationsModel.FindByUserId err %v, req %v", err, in.UserId)
	}
	var res im.GetConversationsResp
	copier.Copy(&res, &data)

	// 根据会话列表，查询具体的会话
	ids := make([]string, 0, len(data.ConversationList))
	for _, conversation := range data.ConversationList {
		ids = append(ids, conversation.ConversationId)
	}
	conversations, err := l.svcCtx.ConversationModel.ListByConversationIds(l.ctx, ids)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "ConversationModel.ListByConversationIds err %v, req %v", err, ids)
	}

	// 计算是否存在未读消息
	for _, conversation := range conversations {
		if _, ok := res.ConversationList[conversation.ConversationId]; !ok {
			continue
		}
		// 用户读取的消息量
		total := res.ConversationList[conversation.ConversationId].Total
		if total < int32(conversation.Total) {
			// 有新的消息
			res.ConversationList[conversation.ConversationId].Total = int32(conversation.Total)
			// 有多少是未读
			res.ConversationList[conversation.ConversationId].ToRead = int32(conversation.Total) - total
			// 更改当前会话为显示状态
			res.ConversationList[conversation.ConversationId].IsShow = true
		}
	}

	return &res, nil
}
