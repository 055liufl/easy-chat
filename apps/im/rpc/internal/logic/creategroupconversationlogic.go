// =============================================================================
// 创建群聊会话业务逻辑
// =============================================================================
// 提供群聊会话创建功能，包括:
//   - 创建群聊会话记录
//   - 为创建者建立群聊会话关系
//
// 数据操作:
//   - Conversation 集合：创建群聊会话
//   - Conversations 集合：为创建者添加群聊会话
//
// 业务场景:
//   - 用户创建群聊时调用
//   - 确保群聊会话的幂等性（重复创建不会出错）
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

// CreateGroupConversationLogic 创建群聊会话业务逻辑处理器
type CreateGroupConversationLogic struct {
	ctx    context.Context        // 请求上下文，用于超时控制和链路追踪
	svcCtx *svc.ServiceContext    // 服务上下文，提供数据模型访问
	logx.Logger                   // 日志记录器
}

// NewCreateGroupConversationLogic 创建群聊会话业务逻辑实例
//
// 参数:
//   - ctx: 请求上下文
//   - svcCtx: 服务上下文，包含数据模型等依赖
//
// 返回:
//   - *CreateGroupConversationLogic: 业务逻辑处理器实例
func NewCreateGroupConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGroupConversationLogic {
	return &CreateGroupConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CreateGroupConversation 创建群聊会话
// 为群聊创建会话记录，并为创建者建立会话关系
//
// 参数:
//   - in: 创建请求
//     - GroupId: 群组 ID（作为会话 ID）
//     - CreateId: 创建者用户 ID
//
// 返回:
//   - *im.CreateGroupConversationResp: 创建响应
//   - error: 错误信息
//
// 业务流程:
//   1. 检查群聊会话是否已存在，存在则直接返回
//   2. 创建群聊会话记录（使用群组 ID 作为会话 ID）
//   3. 为创建者建立群聊会话关系
//
// 说明:
//   - 群聊会话 ID 等于群组 ID
//   - 创建者会自动加入该群聊会话
func (l *CreateGroupConversationLogic) CreateGroupConversation(in *im.CreateGroupConversationReq) (*im.CreateGroupConversationResp, error) {
	// todo: add your logic here and delete this line

	res := &im.CreateGroupConversationResp{}

	_, err := l.svcCtx.ConversationModel.FindOne(l.ctx, in.GroupId)
	if err == nil {
		return res, nil
	}
	if err != immodels.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBErr(), "ConversationModel.FindOne err %v, req %v", err, in.GroupId)
	}

	err = l.svcCtx.ConversationModel.Insert(l.ctx, &immodels.Conversation{
		ConversationId: in.GroupId,
		ChatType:       constants.GroupChatType,
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "ConversationModel.Insert err %v", err)
	}

	_, err = NewSetUpUserConversationLogic(l.ctx, l.svcCtx).SetUpUserConversation(&im.SetUpUserConversationReq{
		SendId:   in.CreateId,
		RecvId:   in.GroupId,
		ChatType: int32(constants.GroupChatType),
	})

	return res, err
}
