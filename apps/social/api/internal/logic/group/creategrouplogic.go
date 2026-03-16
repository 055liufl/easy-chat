// =============================================================================
// 创建群组 Logic - 创建群组业务逻辑
// =============================================================================
// 实现创建新群组的核心业务逻辑
//
// 业务流程:
//   1. 从上下文获取当前用户 ID（群主）
//   2. 调用 Social RPC 创建群组记录
//   3. 调用 IM RPC 创建群组会话
//   4. 返回新创建的群组 ID
//
// 数据来源:
//   - JWT token: 当前用户 ID
//   - 请求参数: 群组名称、图标
//
// =============================================================================
package group

import (
	"context"
	"imooc.com/easy-chat/apps/im/rpc/imclient"
	"imooc.com/easy-chat/apps/social/rpc/socialclient"
	"imooc.com/easy-chat/pkg/ctxdata"

	"imooc.com/easy-chat/apps/social/api/internal/svc"
	"imooc.com/easy-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// CreateGroupLogic 创建群组业务逻辑处理器
type CreateGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewCreateGroupLogic 创建群组业务逻辑处理器实例
//
// 参数:
//   - ctx: 请求上下文，包含用户认证信息等
//   - svcCtx: 服务上下文，包含 RPC 客户端等依赖
//
// 返回:
//   - *CreateGroupLogic: 业务逻辑处理器实例
func NewCreateGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGroupLogic {
	return &CreateGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CreateGroup 创建群组
// 创建新的群组，并建立群组会话
//
// 处理流程:
//   1. 从 JWT token 上下文中提取当前用户 ID（群主）
//   2. 调用 Social RPC 服务创建群组记录
//   3. 如果群组创建成功，调用 IM RPC 服务创建群组会话
//   4. 群组会话用于群聊消息的发送和接收
//
// 参数:
//   - req: 创建群组请求参数
//     - Name: 群组名称
//     - Icon: 群组图标 URL
//
// 返回:
//   - resp: 创建群组响应（当前为空结构）
//   - err: 错误信息
func (l *CreateGroupLogic) CreateGroup(req *types.GroupCreateReq) (resp *types.GroupCreateResp, err error) {
	// 从上下文获取当前登录用户的 ID（群主）
	uid := ctxdata.GetUId(l.ctx)

	// 调用 Social RPC 创建群组
	res, err := l.svcCtx.Social.GroupCreate(l.ctx, &socialclient.GroupCreateReq{
		Name:       req.Name, // 群组名称
		Icon:       req.Icon, // 群组图标
		CreatorUid: uid,      // 创建者 ID（群主）
	})
	if err != nil {
		return nil, err
	}

	// 检查群组是否创建成功
	if res.Id == "" {
		return nil, err
	}

	// 调用 IM RPC 建立群组会话
	// 群组会话用于群聊消息的发送和接收
	_, err = l.svcCtx.Im.CreateGroupConversation(l.ctx, &imclient.CreateGroupConversationReq{
		GroupId:  res.Id, // 群组 ID
		CreateId: uid,    // 创建者 ID
	})

	return nil, err
}
