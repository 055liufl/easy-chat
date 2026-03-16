// =============================================================================
// 群组入群申请处理逻辑
// =============================================================================
// 实现群组入群申请的处理功能，包括:
//   - 验证申请状态（是否已处理）
//   - 更新申请处理结果
//   - 通过申请时创建群成员记录
//
// 业务流程:
//  1. 查询入群申请记录
//  2. 验证申请是否已被处理（通过或拒绝）
//  3. 更新申请处理结果
//  4. 如果通过申请，创建群成员记录（使用事务保证一致性）
//
// 数据流:
//   处理请求 -> 查询申请记录 -> 状态验证 -> 更新申请 -> 创建群成员 -> 响应
//
// 事务保证:
//   申请状态更新和群成员创建在同一事务中，保证数据一致性
// =============================================================================
package logic

import (
	"context"
	"database/sql"
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
	ErrGroupReqBeforePass   = xerr.NewMsg("请求已通过")   // 申请已通过错误
	ErrGroupReqBeforeRefuse = xerr.NewMsg("请求已拒绝") // 申请已拒绝错误
)

// GroupPutInHandleLogic 群组入群申请处理逻辑处理器
type GroupPutInHandleLogic struct {
	ctx    context.Context        // 请求上下文
	svcCtx *svc.ServiceContext     // 服务上下文
	logx.Logger                    // 日志记录器
}

// NewGroupPutInHandleLogic 创建群组入群申请处理逻辑处理器实例
//
// 参数:
//   - ctx: 请求上下文
//   - svcCtx: 服务上下文
//
// 返回:
//   - *GroupPutInHandleLogic: 初始化完成的逻辑处理器实例
func NewGroupPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInHandleLogic {
	return &GroupPutInHandleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GroupPutInHandle 处理群组入群申请（通过或拒绝）
//
// 参数:
//   - in: 入群申请处理请求，包含申请 ID、处理结果（通过/拒绝）、处理人 ID
//
// 返回:
//   - *social.GroupPutInHandleResp: 处理响应，如果通过则包含群组 ID
//   - error: 错误信息（申请不存在、已处理、数据库错误等）
//
// 业务逻辑:
//  1. 根据申请 ID 查询入群申请记录
//  2. 验证申请状态，如果已处理则返回错误
//  3. 使用事务更新申请状态
//  4. 如果处理结果为通过，创建群成员记录
func (l *GroupPutInHandleLogic) GroupPutInHandle(in *social.GroupPutInHandleReq) (*social.GroupPutInHandleResp, error) {
	// todo: add your logic here and delete this line

	// 查询入群申请记录
	groupReq, err := l.svcCtx.GroupRequestsModel.FindOne(l.ctx, int64(in.GroupReqId))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friend req err %v req %v", err, in.GroupReqId)
	}

	// 验证申请是否已被处理
	switch constants.HandlerResult(groupReq.HandleResult.Int64) {
	case constants.PassHandlerResult:
		return nil, errors.WithStack(ErrGroupReqBeforePass)
	case constants.RefuseHandlerResult:
		return nil, errors.WithStack(ErrGroupReqBeforeRefuse)
	}

	// 更新处理结果
	groupReq.HandleResult = sql.NullInt64{
		Int64: int64(in.HandleResult),
		Valid: true,
	}

	// 使用事务：更新申请状态 -> 如果通过则创建群成员记录
	err = l.svcCtx.GroupRequestsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// 更新入群申请记录
		if err := l.svcCtx.GroupRequestsModel.Update(l.ctx, session, groupReq); err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "update friend req err %v req %v", err, groupReq)
		}

		// 如果不是通过，直接返回
		if constants.HandlerResult(groupReq.HandleResult.Int64) != constants.PassHandlerResult {
			return nil
		}

		// 创建群成员记录
		groupMember := &socialmodels.GroupMembers{
			GroupId:     groupReq.GroupId,                   // 群组 ID
			UserId:      groupReq.ReqId,                     // 用户 ID（申请人）
			RoleLevel:   int(constants.AtLargeGroupRoleLevel), // 角色等级（普通成员）
			OperatorUid: in.HandleUid,                       // 操作人 ID（处理人）
		}
		_, err = l.svcCtx.GroupMembersModel.Insert(l.ctx, session, groupMember)
		if err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "insert friend err %v req %v", err, groupMember)
		}

		return nil
	})

	// 如果不是通过，返回空响应
	if constants.HandlerResult(groupReq.HandleResult.Int64) != constants.PassHandlerResult {
		return &social.GroupPutInHandleResp{}, err
	}

	// 通过申请，返回群组 ID
	return &social.GroupPutInHandleResp{
		GroupId: groupReq.GroupId,
	}, err
}
