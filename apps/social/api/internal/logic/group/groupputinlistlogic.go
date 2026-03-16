// =============================================================================
// 入群申请列表 Logic - 获取入群申请列表业务逻辑
// =============================================================================
// 实现获取群组收到的入群申请列表的核心业务逻辑
//
// 业务流程:
//   1. 调用 Social RPC 获取指定群组的入群申请列表
//   2. 使用 copier 库将 RPC 响应转换为 API 响应格式
//
// 数据来源:
//   - Social RPC: 入群申请记录列表
//
// =============================================================================
package group

import (
	"context"
	"github.com/jinzhu/copier"
	"imooc.com/easy-chat/apps/social/rpc/socialclient"

	"imooc.com/easy-chat/apps/social/api/internal/svc"
	"imooc.com/easy-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// GroupPutInListLogic 入群申请列表业务逻辑处理器
type GroupPutInListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGroupPutInListLogic 创建入群申请列表业务逻辑处理器实例
//
// 参数:
//   - ctx: 请求上下文，包含用户认证信息等
//   - svcCtx: 服务上下文，包含 RPC 客户端等依赖
//
// 返回:
//   - *GroupPutInListLogic: 业务逻辑处理器实例
func NewGroupPutInListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInListLogic {
	return &GroupPutInListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GroupPutInList 获取入群申请列表
// 查询指定群组收到的所有入群申请记录
//
// 处理流程:
//   1. 调用 Social RPC 服务获取入群申请列表
//   2. 使用 copier 库将 RPC 响应数据结构转换为 API 响应格式
//   3. 返回申请列表，包含申请人信息、申请消息、处理状态等
//
// 参数:
//   - req: 入群申请列表请求参数
//     - GroupId: 群组 ID
//
// 返回:
//   - resp: 入群申请列表响应，包含申请记录列表
//   - err: 错误信息
func (l *GroupPutInListLogic) GroupPutInList(req *types.GroupPutInListRep) (resp *types.GroupPutInListResp, err error) {
	// 调用 Social RPC 获取入群申请列表
	list, err := l.svcCtx.Social.GroupPutinList(l.ctx, &socialclient.GroupPutinListReq{
		GroupId: req.GroupId, // 群组 ID
	})

	// 使用 copier 库将 RPC 响应转换为 API 响应格式
	var respList []*types.GroupRequests
	copier.Copy(&respList, list.List)

	return &types.GroupPutInListResp{List: respList}, nil
}
