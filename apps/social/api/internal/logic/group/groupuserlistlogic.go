// =============================================================================
// 群成员列表 Logic - 获取群成员列表业务逻辑
// =============================================================================
// 实现获取群组成员列表的核心业务逻辑
//
// 业务流程:
//   1. 调用 Social RPC 获取群成员列表
//   2. 提取所有成员的用户 ID
//   3. 调用 User RPC 批量查询成员的详细信息
//   4. 将成员关系和用户信息合并，组装返回数据
//
// 数据来源:
//   - Social RPC: 群成员关系数据（用户ID、角色等级等）
//   - User RPC: 用户详细信息（昵称、头像等）
//
// =============================================================================
package group

import (
	"context"
	"imooc.com/easy-chat/apps/social/rpc/socialclient"
	"imooc.com/easy-chat/apps/user/rpc/userclient"

	"imooc.com/easy-chat/apps/social/api/internal/svc"
	"imooc.com/easy-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// GroupUserListLogic 群成员列表业务逻辑处理器
type GroupUserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGroupUserListLogic 创建群成员列表业务逻辑处理器实例
//
// 参数:
//   - ctx: 请求上下文，包含用户认证信息等
//   - svcCtx: 服务上下文，包含 RPC 客户端等依赖
//
// 返回:
//   - *GroupUserListLogic: 业务逻辑处理器实例
func NewGroupUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupUserListLogic {
	return &GroupUserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GroupUserList 获取群成员列表
// 查询指定群组的所有成员，并返回成员的详细信息
//
// 处理流程:
//   1. 调用 Social RPC 服务获取群成员列表
//   2. 提取所有成员的用户 ID，构建查询列表
//   3. 调用 User RPC 服务批量查询成员的用户信息
//   4. 将用户信息构建为 map，便于快速查找
//   5. 遍历群成员关系，关联用户信息，组装最终响应数据
//
// 参数:
//   - req: 群成员列表请求参数
//     - GroupId: 群组 ID
//
// 返回:
//   - resp: 群成员列表响应，包含成员ID、用户ID、角色等级、昵称、头像等信息
//   - err: 错误信息
func (l *GroupUserListLogic) GroupUserList(req *types.GroupUserListReq) (resp *types.GroupUserListResp, err error) {
	// 调用 Social RPC 获取群成员列表
	groupUsers, err := l.svcCtx.Social.GroupUsers(l.ctx, &socialclient.GroupUsersReq{
		GroupId: req.GroupId,
	})

	// 提取所有成员的用户 ID，用于批量查询用户信息
	uids := make([]string, 0, len(groupUsers.List))
	for _, v := range groupUsers.List {
		uids = append(uids, v.UserId)
	}

	// 调用 User RPC 批量查询成员的用户信息
	userList, err := l.svcCtx.User.FindUser(l.ctx, &userclient.FindUserReq{Ids: uids})
	if err != nil {
		return nil, err
	}

	// 将用户信息构建为 map，key 为用户 ID，便于后续快速查找
	userRecords := make(map[string]*userclient.UserEntity, len(userList.User))
	for i, _ := range userList.User {
		userRecords[userList.User[i].Id] = userList.User[i]
	}

	// 组装响应数据：将群成员关系和用户信息合并
	respList := make([]*types.GroupMembers, 0, len(groupUsers.List))
	for _, v := range groupUsers.List {

		member := &types.GroupMembers{
			Id:        int64(v.Id),       // 成员记录 ID
			GroupId:   v.GroupId,         // 群组 ID
			UserId:    v.UserId,          // 用户 ID
			RoleLevel: int(v.RoleLevel),  // 角色等级（1-群主，2-管理员，3-普通成员）
		}
		// 如果找到对应的用户信息，填充昵称和头像
		if u, ok := userRecords[v.UserId]; ok {
			member.Nickname = u.Nickname
			member.UserAvatarUrl = u.Avatar
		}
		respList = append(respList, member)
	}

	return &types.GroupUserListResp{List: respList}, err
}
