// =============================================================================
// 群组入群申请逻辑
// =============================================================================
// 实现群组入群申请功能，包括:
//   - 用户主动申请入群
//   - 群成员邀请用户入群
//   - 群管理员/群主邀请用户入群（直接通过）
//   - 根据群组验证设置决定是否需要审核
//
// 业务流程:
//  1. 检查用户是否已是群成员
//  2. 检查是否已有待处理的入群申请
//  3. 查询群组信息，判断是否需要验证
//  4. 根据入群方式（申请/邀请）和邀请人角色决定处理方式:
//     - 群无验证：直接通过，创建群成员记录
//     - 用户申请：创建待审核的申请记录
//     - 普通成员邀请：创建待审核的申请记录
//     - 管理员/群主邀请：直接通过，创建群成员记录
//
// 数据流:
//   入群请求 -> 成员检查 -> 申请检查 -> 群组验证检查 -> 入群方式判断 -> 创建申请/成员 -> 响应
//
// 入群方式:
//   - PutInGroupJoinSource: 用户主动申请
//   - InviteGroupJoinSource: 成员邀请
// =============================================================================
package logic

import (
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"

	"imooc.com/easy-chat/apps/social/rpc/internal/svc"
	"imooc.com/easy-chat/apps/social/rpc/social"
	"imooc.com/easy-chat/apps/social/socialmodels"
	"imooc.com/easy-chat/pkg/constants"
	"imooc.com/easy-chat/pkg/xerr"
)

// GroupPutinLogic 群组入群申请逻辑处理器
type GroupPutinLogic struct {
	ctx    context.Context        // 请求上下文
	svcCtx *svc.ServiceContext     // 服务上下文
	logx.Logger                    // 日志记录器
}

// NewGroupPutinLogic 创建群组入群申请逻辑处理器实例
//
// 参数:
//   - ctx: 请求上下文
//   - svcCtx: 服务上下文
//
// 返回:
//   - *GroupPutinLogic: 初始化完成的逻辑处理器实例
func NewGroupPutinLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutinLogic {
	return &GroupPutinLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GroupPutin 处理群组入群申请/邀请
//
// 参数:
//   - in: 入群请求，包含申请人 ID、群组 ID、入群方式、邀请人 ID、申请消息等
//
// 返回:
//   - *social.GroupPutinResp: 入群响应，如果直接通过则包含群组 ID
//   - error: 错误信息（已是成员、已有申请、数据库错误等）
//
// 业务逻辑:
//  1. 检查申请人是否已是群成员
//  2. 检查是否已有未处理的入群申请
//  3. 查询群组信息，判断是否需要验证
//  4. 根据入群方式和邀请人角色决定处理方式
//
// 处理规则:
//  - 群无验证：直接通过，创建群成员记录
//  - 用户申请（PutInGroupJoinSource）：创建待审核的申请记录
//  - 普通成员邀请：创建待审核的申请记录
//  - 管理员/群主邀请：直接通过，创建群成员记录
func (l *GroupPutinLogic) GroupPutin(in *social.GroupPutinReq) (*social.GroupPutinResp, error) {
	// todo: add your logic here and delete this line

	//  1. 普通用户申请 ： 如果群无验证直接进入
	//  2. 群成员邀请： 如果群无验证直接进入
	//  3. 群管理员/群创建者邀请：直接进入群

	var (
		inviteGroupMember *socialmodels.GroupMembers // 邀请人的群成员信息
		userGroupMember   *socialmodels.GroupMembers // 申请人的群成员信息
		groupInfo         *socialmodels.Groups       // 群组信息

		err error
	)

	// 检查申请人是否已是群成员
	userGroupMember, err = l.svcCtx.GroupMembersModel.FindByGroudIdAndUserId(l.ctx, in.ReqId, in.GroupId)
	if err != nil && err != socialmodels.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find group member by groud id and  req id err %v, req %v, %v", err,
			in.GroupId, in.ReqId)
	}
	if userGroupMember != nil {
		// 已是群成员，直接返回成功
		return &social.GroupPutinResp{}, nil
	}

	// 检查是否已有未处理的入群申请
	groupReq, err := l.svcCtx.GroupRequestsModel.FindByGroupIdAndReqId(l.ctx, in.GroupId, in.ReqId)
	if err != nil && err != socialmodels.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find group req by groud id and user id err %v, req %v, %v", err,
			in.GroupId, in.ReqId)
	}
	if groupReq != nil {
		// 已有申请记录，直接返回成功
		return &social.GroupPutinResp{}, nil
	}

	// 构建入群申请记录
	groupReq = &socialmodels.GroupRequests{
		ReqId:   in.ReqId,   // 申请人 ID
		GroupId: in.GroupId, // 群组 ID
		ReqMsg: sql.NullString{ // 申请消息
			String: in.ReqMsg,
			Valid:  true,
		},
		ReqTime: sql.NullTime{ // 申请时间
			Time:  time.Unix(in.ReqTime, 0),
			Valid: true,
		},
		JoinSource: sql.NullInt64{ // 入群方式（申请/邀请）
			Int64: int64(in.JoinSource),
			Valid: true,
		},
		InviterUserId: sql.NullString{ // 邀请人 ID
			String: in.InviterUid,
			Valid:  true,
		},
		HandleResult: sql.NullInt64{ // 处理结果（初始为未处理）
			Int64: int64(constants.NoHandlerResult),
			Valid: true,
		},
	}

	// 定义延迟执行的创建群成员函数
	createGroupMember := func() {
		if err != nil {
			return
		}
		err = l.createGroupMember(in)
	}

	// 查询群组信息
	groupInfo, err = l.svcCtx.GroupsModel.FindOne(l.ctx, in.GroupId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find group by groud id err %v, req %v", err, in.GroupId)
	}

	// 验证群组是否需要验证
	if !groupInfo.IsVerify {
		// 群组不需要验证，直接通过
		defer createGroupMember()

		groupReq.HandleResult = sql.NullInt64{
			Int64: int64(constants.PassHandlerResult),
			Valid: true,
		}

		return l.createGroupReq(groupReq, true)
	}

	// 验证入群方式
	if constants.GroupJoinSource(in.JoinSource) == constants.PutInGroupJoinSource {
		// 用户主动申请，创建待审核的申请记录
		return l.createGroupReq(groupReq, false)
	}

	// 查询邀请人的群成员信息
	inviteGroupMember, err = l.svcCtx.GroupMembersModel.FindByGroudIdAndUserId(l.ctx, in.InviterUid, in.GroupId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find group member by groud id and user id err %v, req %v",
			in.InviterUid, in.GroupId)
	}

	// 判断邀请人角色
	if constants.GroupRoleLevel(inviteGroupMember.RoleLevel) == constants.CreatorGroupRoleLevel || constants.
		GroupRoleLevel(inviteGroupMember.RoleLevel) == constants.ManagerGroupRoleLevel {
		// 邀请人是管理员或群主，直接通过
		defer createGroupMember()

		groupReq.HandleResult = sql.NullInt64{
			Int64: int64(constants.PassHandlerResult),
			Valid: true,
		}
		groupReq.HandleUserId = sql.NullString{
			String: in.InviterUid,
			Valid:  true,
		}
		return l.createGroupReq(groupReq, true)
	}
	// 普通成员邀请，创建待审核的申请记录
	return l.createGroupReq(groupReq, false)

}

// createGroupReq 创建入群申请记录
//
// 参数:
//   - groupReq: 入群申请记录
//   - isPass: 是否直接通过（true: 直接通过，false: 待审核）
//
// 返回:
//   - *social.GroupPutinResp: 入群响应，如果直接通过则包含群组 ID
//   - error: 错误信息（数据库错误等）
//
// 业务逻辑:
//  1. 插入入群申请记录到数据库
//  2. 如果直接通过，返回包含群组 ID 的响应
//  3. 如果待审核，返回空响应
func (l *GroupPutinLogic) createGroupReq(groupReq *socialmodels.GroupRequests, isPass bool) (*social.GroupPutinResp, error) {

	_, err := l.svcCtx.GroupRequestsModel.Insert(l.ctx, groupReq)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "insert group req err %v req %v", err, groupReq)
	}

	if isPass {
		// 直接通过，返回群组 ID
		return &social.GroupPutinResp{GroupId: groupReq.GroupId}, nil
	}

	// 待审核，返回空响应
	return &social.GroupPutinResp{}, nil
}

// createGroupMember 创建群成员记录
//
// 参数:
//   - in: 入群请求，包含申请人 ID、群组 ID、邀请人 ID
//
// 返回:
//   - error: 错误信息（数据库错误等）
//
// 业务逻辑:
//  1. 构建群成员记录（角色为普通成员）
//  2. 插入群成员记录到数据库
func (l *GroupPutinLogic) createGroupMember(in *social.GroupPutinReq) error {
	groupMember := &socialmodels.GroupMembers{
		GroupId:     in.GroupId,                             // 群组 ID
		UserId:      in.ReqId,                               // 用户 ID（申请人）
		RoleLevel:   int(constants.AtLargeGroupRoleLevel),   // 角色等级（普通成员）
		OperatorUid: in.InviterUid,                          // 操作人 ID（邀请人）
	}
	_, err := l.svcCtx.GroupMembersModel.Insert(l.ctx, nil, groupMember)
	if err != nil {
		return errors.Wrapf(xerr.NewDBErr(), "insert friend err %v req %v", err, groupMember)
	}

	return nil
}
