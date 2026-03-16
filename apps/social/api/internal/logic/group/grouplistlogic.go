// =============================================================================
// 群组列表 Logic - 获取群组列表业务逻辑
// =============================================================================
// 实现获取用户加入的群组列表的核心业务逻辑
//
// 业务流程:
//   1. 从上下文获取当前用户 ID
//   2. 调用 Social RPC 获取用户加入的群组列表
//   3. 使用 copier 库将 RPC 响应转换为 API 响应格式
//
// 数据来源:
//   - Social RPC: 群组列表数据
//
// =============================================================================
package group

import (
	"context"
	"github.com/jinzhu/copier"
	"imooc.com/easy-chat/apps/social/rpc/socialclient"
	"imooc.com/easy-chat/pkg/ctxdata"

	"imooc.com/easy-chat/apps/social/api/internal/svc"
	"imooc.com/easy-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// GroupListLogic 群组列表业务逻辑处理器
type GroupListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGroupListLogic 创建群组列表业务逻辑处理器实例
//
// 参数:
//   - ctx: 请求上下文，包含用户认证信息等
//   - svcCtx: 服务上下文，包含 RPC 客户端等依赖
//
// 返回:
//   - *GroupListLogic: 业务逻辑处理器实例
func NewGroupListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupListLogic {
	return &GroupListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GroupList 获取群组列表
// 查询当前用户加入的所有群组
//
// 处理流程:
//   1. 从 JWT token 上下文中提取当前用户 ID
//   2. 调用 Social RPC 服务获取用户加入的群组列表
//   3. 使用 copier 库将 RPC 响应数据结构转换为 API 响应格式
//   4. 返回群组列表，包含群组 ID、名称、图标、创建者等信息
//
// 参数:
//   - req: 群组列表请求参数（当前为空结构）
//
// 返回:
//   - resp: 群组列表响应，包含群组列表
//   - err: 错误信息
func (l *GroupListLogic) GroupList(req *types.GroupListRep) (resp *types.GroupListResp, err error) {
	// 从上下文获取当前登录用户的 ID
	uid := ctxdata.GetUId(l.ctx)

	// 调用 Social RPC 获取用户加入的群组列表
	list, err := l.svcCtx.Social.GroupList(l.ctx, &socialclient.GroupListReq{
		UserId: uid,
	})

	if err != nil {
		return nil, err
	}

	// 使用 copier 库将 RPC 响应转换为 API 响应格式
	var respList []*types.Groups
	copier.Copy(&respList, list.List)

	return &types.GroupListResp{List: respList}, nil
}
