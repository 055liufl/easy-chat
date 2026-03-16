// =============================================================================
// 好友申请逻辑
// =============================================================================
// 实现好友申请功能，包括:
//   - 验证是否已经是好友关系
//   - 验证是否已有待处理的申请
//   - 创建好友申请记录
//
// 业务流程:
//  1. 检查申请人与目标用户是否已经是好友
//  2. 检查是否已有未处理的好友申请
//  3. 创建新的好友申请记录（状态为待处理）
//
// 数据流:
//   申请请求 -> 好友关系检查 -> 申请记录检查 -> 创建申请记录 -> 响应
// =============================================================================
package logic

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"imooc.com/easy-chat/apps/social/socialmodels"
	"imooc.com/easy-chat/pkg/constants"
	"imooc.com/easy-chat/pkg/xerr"
	"time"

	"imooc.com/easy-chat/apps/social/rpc/internal/svc"
	"imooc.com/easy-chat/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

// FriendPutInLogic 好友申请逻辑处理器
type FriendPutInLogic struct {
	ctx    context.Context        // 请求上下文
	svcCtx *svc.ServiceContext     // 服务上下文
	logx.Logger                    // 日志记录器
}

// NewFriendPutInLogic 创建好友申请逻辑处理器实例
//
// 参数:
//   - ctx: 请求上下文
//   - svcCtx: 服务上下文
//
// 返回:
//   - *FriendPutInLogic: 初始化完成的逻辑处理器实例
func NewFriendPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInLogic {
	return &FriendPutInLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// FriendPutIn 处理好友申请请求
//
// 参数:
//   - in: 好友申请请求，包含申请人 ID、目标用户 ID、申请消息等
//
// 返回:
//   - *social.FriendPutInResp: 好友申请响应
//   - error: 错误信息（已是好友、已有申请、数据库错误等）
//
// 业务逻辑:
//  1. 检查申请人与目标用户是否已经是好友关系
//  2. 检查是否已有未处理的好友申请记录
//  3. 创建新的好友申请记录，状态为待处理
func (l *FriendPutInLogic) FriendPutIn(in *social.FriendPutInReq) (*social.FriendPutInResp, error) {
	// todo: add your logic here and delete this line

	// 检查申请人是否与目标用户已经是好友关系
	friends, err := l.svcCtx.FriendsModel.FindByUidAndFid(l.ctx, in.UserId, in.ReqUid)
	if err != nil && err != socialmodels.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friends by uid and fid err %v req %v ", err, in)
	}
	if friends != nil {
		// 已经是好友，直接返回成功
		return &social.FriendPutInResp{}, err
	}

	// 检查是否已经有过申请，且申请未完成
	friendReqs, err := l.svcCtx.FriendRequestsModel.FindByReqUidAndUserId(l.ctx, in.ReqUid, in.UserId)
	if err != nil && err != socialmodels.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friendsRequest by rid and uid err %v req %v ", err, in)
	}
	if friendReqs != nil {
		// 已有申请记录，直接返回成功
		return &social.FriendPutInResp{}, err
	}

	// 创建好友申请记录
	_, err = l.svcCtx.FriendRequestsModel.Insert(l.ctx, &socialmodels.FriendRequests{
		UserId: in.UserId,                // 被申请人 ID
		ReqUid: in.ReqUid,                // 申请人 ID
		ReqMsg: sql.NullString{           // 申请消息
			Valid:  true,
			String: in.ReqMsg,
		},
		ReqTime: time.Unix(in.ReqTime, 0), // 申请时间
		HandleResult: sql.NullInt64{       // 处理结果（初始为未处理）
			Int64: int64(constants.NoHandlerResult),
			Valid: true,
		},
	})

	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "insert friendRequest err %v req %v ", err, in)
	}

	return &social.FriendPutInResp{}, nil
}
