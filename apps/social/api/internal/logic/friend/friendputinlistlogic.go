// =============================================================================
// 好友申请列表 Logic - 获取好友申请列表业务逻辑
// =============================================================================
// 实现获取用户收到的好友申请列表的核心业务逻辑
//
// 业务流程:
//   1. 从上下文获取当前用户 ID
//   2. 调用 Social RPC 获取好友申请列表
//   3. 使用 copier 库将 RPC 响应转换为 API 响应格式
//
// 数据来源:
//   - Social RPC: 好友申请记录列表
//
// =============================================================================
package friend

import (
	"context"
	"github.com/jinzhu/copier"
	"imooc.com/easy-chat/apps/social/rpc/socialclient"
	"imooc.com/easy-chat/pkg/ctxdata"

	"imooc.com/easy-chat/apps/social/api/internal/svc"
	"imooc.com/easy-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// FriendPutInListLogic 好友申请列表业务逻辑处理器
type FriendPutInListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewFriendPutInListLogic 创建好友申请列表业务逻辑处理器实例
//
// 参数:
//   - ctx: 请求上下文，包含用户认证信息等
//   - svcCtx: 服务上下文，包含 RPC 客户端等依赖
//
// 返回:
//   - *FriendPutInListLogic: 业务逻辑处理器实例
func NewFriendPutInListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInListLogic {
	return &FriendPutInListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// FriendPutInList 获取好友申请列表
// 查询当前用户收到的所有好友申请记录
//
// 处理流程:
//   1. 从 JWT token 上下文中提取当前用户 ID
//   2. 调用 Social RPC 服务获取好友申请列表
//   3. 使用 copier 库将 RPC 响应数据结构转换为 API 响应格式
//   4. 返回申请列表，包含申请人信息、申请消息、处理状态等
//
// 参数:
//   - req: 好友申请列表请求参数（当前为空结构）
//
// 返回:
//   - resp: 好友申请列表响应，包含申请记录列表
//   - err: 错误信息
func (l *FriendPutInListLogic) FriendPutInList(req *types.FriendPutInListReq) (resp *types.FriendPutInListResp, err error) {
	// 调用 Social RPC 获取好友申请列表
	list, err := l.svcCtx.Social.FriendPutInList(l.ctx, &socialclient.FriendPutInListReq{
		UserId: ctxdata.GetUId(l.ctx), // 当前用户 ID
	})
	if err != nil {
		return nil, err
	}

	// 使用 copier 库将 RPC 响应转换为 API 响应格式
	var respList []*types.FriendRequests
	copier.Copy(&respList, list.List)

	return &types.FriendPutInListResp{List: respList}, nil
}
