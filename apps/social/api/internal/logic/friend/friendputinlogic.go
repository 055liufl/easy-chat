// =============================================================================
// 好友申请 Logic - 发送好友申请业务逻辑
// =============================================================================
// 实现发送好友申请的核心业务逻辑
//
// 业务流程:
//   1. 从上下文获取当前用户 ID（申请发起人）
//   2. 调用 Social RPC 创建好友申请记录
//   3. 记录申请消息、申请时间等信息
//
// 数据来源:
//   - JWT token: 当前用户 ID
//   - 请求参数: 目标用户 ID、申请消息、申请时间
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

// FriendPutInLogic 好友申请业务逻辑处理器
type FriendPutInLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewFriendPutInLogic 创建好友申请业务逻辑处理器实例
//
// 参数:
//   - ctx: 请求上下文，包含用户认证信息等
//   - svcCtx: 服务上下文，包含 RPC 客户端等依赖
//
// 返回:
//   - *FriendPutInLogic: 业务逻辑处理器实例
func NewFriendPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInLogic {
	return &FriendPutInLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// FriendPutIn 发送好友申请
// 向指定用户发送好友申请，记录申请消息和时间
//
// 处理流程:
//   1. 从 JWT token 上下文中提取当前用户 ID（申请发起人）
//   2. 调用 Social RPC 服务创建好友申请记录
//   3. 传递申请人 ID、目标用户 ID、申请消息、申请时间
//
// 参数:
//   - req: 好友申请请求参数
//     - UserId: 目标用户 ID（被申请人）
//     - ReqMsg: 申请消息（如"你好，我想加你为好友"）
//     - ReqTime: 申请时间戳
//
// 返回:
//   - resp: 好友申请响应（空结构）
//   - err: 错误信息
func (l *FriendPutInLogic) FriendPutIn(req *types.FriendPutInReq) (resp *types.FriendPutInResp, err error) {
	// 从上下文获取当前登录用户的 ID（申请发起人）
	uid := ctxdata.GetUId(l.ctx)

	// 调用 Social RPC 创建好友申请记录
	_, err = l.svcCtx.Social.FriendPutIn(l.ctx, &socialclient.FriendPutInReq{
		UserId:  uid,        // 申请发起人 ID
		ReqUid:  req.UserId, // 目标用户 ID（被申请人）
		ReqMsg:  req.ReqMsg, // 申请消息
		ReqTime: req.ReqTime, // 申请时间
	})

	return
}
