// =============================================================================
// 好友申请处理逻辑
// =============================================================================
// 实现好友申请的处理功能，包括:
//   - 验证申请状态（是否已处理）
//   - 更新申请处理结果
//   - 通过申请时建立双向好友关系
//
// 业务流程:
//  1. 查询好友申请记录
//  2. 验证申请是否已被处理（通过或拒绝）
//  3. 更新申请处理结果
//  4. 如果通过申请，创建双向好友关系记录（使用事务保证一致性）
//
// 数据流:
//   处理请求 -> 查询申请记录 -> 状态验证 -> 更新申请 -> 创建好友关系 -> 响应
//
// 事务保证:
//   申请状态更新和好友关系创建在同一事务中，保证数据一致性
// =============================================================================
package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"imooc.com/easy-chat/apps/social/socialmodels"
	"imooc.com/easy-chat/pkg/constants"
	"imooc.com/easy-chat/pkg/xerr"

	"imooc.com/easy-chat/apps/social/rpc/internal/svc"
	"imooc.com/easy-chat/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrFriendReqBeforePass   = xerr.NewMsg("好友申请并已经通过")   // 申请已通过错误
	ErrFriendReqBeforeRefuse = xerr.NewMsg("好友申请已经被拒绝") // 申请已拒绝错误
)

// FriendPutInHandleLogic 好友申请处理逻辑处理器
type FriendPutInHandleLogic struct {
	ctx    context.Context        // 请求上下文
	svcCtx *svc.ServiceContext     // 服务上下文
	logx.Logger                    // 日志记录器
}

// NewFriendPutInHandleLogic 创建好友申请处理逻辑处理器实例
//
// 参数:
//   - ctx: 请求上下文
//   - svcCtx: 服务上下文
//
// 返回:
//   - *FriendPutInHandleLogic: 初始化完成的逻辑处理器实例
func NewFriendPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInHandleLogic {
	return &FriendPutInHandleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// FriendPutInHandle 处理好友申请（通过或拒绝）
//
// 参数:
//   - in: 好友申请处理请求，包含申请 ID、处理结果（通过/拒绝）
//
// 返回:
//   - *social.FriendPutInHandleResp: 处理响应
//   - error: 错误信息（申请不存在、已处理、数据库错误等）
//
// 业务逻辑:
//  1. 根据申请 ID 查询好友申请记录
//  2. 验证申请状态，如果已处理则返回错误
//  3. 使用事务更新申请状态
//  4. 如果处理结果为通过，创建双向好友关系记录
func (l *FriendPutInHandleLogic) FriendPutInHandle(in *social.FriendPutInHandleReq) (*social.FriendPutInHandleResp, error) {
	// todo: add your logic here and delete this line

	// 获取好友申请记录
	friendReq, err := l.svcCtx.FriendRequestsModel.FindOne(l.ctx, int64(in.FriendReqId))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friendsRequest by friendReqid err %v req %v ", err,
			in.FriendReqId)
	}

	// 验证申请是否已被处理
	switch constants.HandlerResult(friendReq.HandleResult.Int64) {
	case constants.PassHandlerResult:
		return nil, errors.WithStack(ErrFriendReqBeforePass)
	case constants.RefuseHandlerResult:
		return nil, errors.WithStack(ErrFriendReqBeforeRefuse)
	}

	// 更新处理结果
	friendReq.HandleResult.Int64 = int64(in.HandleResult)

	// 使用事务：修改申请结果 -> 如果通过则建立两条好友关系记录
	err = l.svcCtx.FriendRequestsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// 更新好友申请记录
		if err := l.svcCtx.FriendRequestsModel.Update(l.ctx, session, friendReq); err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "update friend request err %v, req %v", err, friendReq)
		}

		// 如果不是通过，直接返回
		if constants.HandlerResult(in.HandleResult) != constants.PassHandlerResult {
			return nil
		}

		// 创建双向好友关系记录
		friends := []*socialmodels.Friends{
			{
				UserId:    friendReq.UserId,  // 被申请人 -> 申请人
				FriendUid: friendReq.ReqUid,
			}, {
				UserId:    friendReq.ReqUid,  // 申请人 -> 被申请人
				FriendUid: friendReq.UserId,
			},
		}

		// 批量插入好友关系记录
		_, err = l.svcCtx.FriendsModel.Inserts(l.ctx, session, friends...)
		if err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "friends inserts err %v, req %v", err, friends)
		}
		return nil
	})

	return &social.FriendPutInHandleResp{}, err
}
