// =============================================================================
// 好友申请处理 Logic - 处理好友申请业务逻辑（同意/拒绝）
// =============================================================================
// 实现处理好友申请的核心业务逻辑
//
// 业务流程:
//   1. 从上下文获取当前用户 ID（申请接收人）
//   2. 调用 Social RPC 处理好友申请
//   3. 根据处理结果（同意/拒绝）更新申请状态
//   4. 如果同意，建立双向好友关系
//
// 数据来源:
//   - JWT token: 当前用户 ID
//   - 请求参数: 申请记录 ID、处理结果
//
// =============================================================================
package friend

import (
	"context"
	"imooc.com/easy-chat/apps/social/rpc/socialclient"
	"imooc.com/easy-chat/pkg/ctxdata"

	"imooc.com/easy-chat/apps/social/api/internal/svc"
	"imooc.com/easy-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// FriendPutInHandleLogic 好友申请处理业务逻辑处理器
type FriendPutInHandleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewFriendPutInHandleLogic 创建好友申请处理业务逻辑处理器实例
//
// 参数:
//   - ctx: 请求上下文，包含用户认证信息等
//   - svcCtx: 服务上下文，包含 RPC 客户端等依赖
//
// 返回:
//   - *FriendPutInHandleLogic: 业务逻辑处理器实例
func NewFriendPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInHandleLogic {
	return &FriendPutInHandleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// FriendPutInHandle 处理好友申请
// 对收到的好友申请进行处理（同意或拒绝）
//
// 处理流程:
//   1. 从 JWT token 上下文中提取当前用户 ID（申请接收人）
//   2. 调用 Social RPC 服务处理好友申请
//   3. 传递申请记录 ID、当前用户 ID、处理结果
//   4. RPC 服务会根据处理结果更新申请状态
//   5. 如果同意（HandleResult=1），会建立双向好友关系
//
// 参数:
//   - req: 好友申请处理请求参数
//     - FriendReqId: 好友申请记录 ID
//     - HandleResult: 处理结果（1-同意，2-拒绝）
//
// 返回:
//   - resp: 好友申请处理响应（空结构）
//   - err: 错误信息
func (l *FriendPutInHandleLogic) FriendPutInHandle(req *types.FriendPutInHandleReq) (resp *types.FriendPutInHandleResp, err error) {
	// 调用 Social RPC 处理好友申请
	_, err = l.svcCtx.Social.FriendPutInHandle(l.ctx, &socialclient.FriendPutInHandleReq{
		FriendReqId:  req.FriendReqId,           // 好友申请记录 ID
		UserId:       ctxdata.GetUId(l.ctx),     // 当前用户 ID（申请接收人）
		HandleResult: req.HandleResult,          // 处理结果（1-同意，2-拒绝）
	})

	return
}
