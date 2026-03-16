// =============================================================================
// 好友列表 Logic - 好友列表业务逻辑
// =============================================================================
// 实现获取用户好友列表的核心业务逻辑
//
// 业务流程:
//   1. 从上下文获取当前用户 ID
//   2. 调用 Social RPC 获取好友关系列表
//   3. 提取所有好友的用户 ID
//   4. 调用 User RPC 批量查询好友的详细信息（昵称、头像等）
//   5. 将好友关系和用户信息合并，组装返回数据
//
// 数据来源:
//   - Social RPC: 好友关系数据（好友ID列表）
//   - User RPC: 用户详细信息（昵称、头像等）
//
// =============================================================================
package friend

import (
	"context"
	"imooc.com/easy-chat/apps/social/rpc/socialclient"
	"imooc.com/easy-chat/apps/user/rpc/userclient"
	"imooc.com/easy-chat/pkg/ctxdata"

	"imooc.com/easy-chat/apps/social/api/internal/svc"
	"imooc.com/easy-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// FriendListLogic 好友列表业务逻辑处理器
type FriendListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewFriendListLogic 创建好友列表业务逻辑处理器实例
//
// 参数:
//   - ctx: 请求上下文，包含用户认证信息等
//   - svcCtx: 服务上下文，包含 RPC 客户端等依赖
//
// 返回:
//   - *FriendListLogic: 业务逻辑处理器实例
func NewFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendListLogic {
	return &FriendListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// FriendList 获取好友列表
// 查询当前用户的所有好友，并返回好友的详细信息
//
// 处理流程:
//   1. 从 JWT token 上下文中提取当前用户 ID
//   2. 调用 Social RPC 服务获取好友关系列表
//   3. 如果没有好友，直接返回空列表
//   4. 提取所有好友的用户 ID，构建查询列表
//   5. 调用 User RPC 服务批量查询好友的用户信息
//   6. 将用户信息构建为 map，便于快速查找
//   7. 遍历好友关系，关联用户信息，组装最终响应数据
//
// 参数:
//   - req: 好友列表请求参数（当前为空结构）
//
// 返回:
//   - resp: 好友列表响应，包含好友ID、昵称、头像等信息
//   - err: 错误信息
func (l *FriendListLogic) FriendList(req *types.FriendListReq) (resp *types.FriendListResp, err error) {
	// 从上下文获取当前登录用户的 ID
	uid := ctxdata.GetUId(l.ctx)

	// 调用 Social RPC 获取好友关系列表
	friends, err := l.svcCtx.Social.FriendList(l.ctx, &socialclient.FriendListReq{
		UserId: uid,
	})
	if err != nil {
		return nil, err
	}

	// 如果没有好友，返回空列表
	if len(friends.List) == 0 {
		return &types.FriendListResp{}, nil
	}

	// 提取所有好友的用户 ID，用于批量查询用户信息
	uids := make([]string, 0, len(friends.List))
	for _, i := range friends.List {
		uids = append(uids, i.FriendUid)
	}

	// 调用 User RPC 批量查询好友的用户信息
	users, err := l.svcCtx.User.FindUser(l.ctx, &userclient.FindUserReq{
		Ids: uids,
	})
	if err != nil {
		return &types.FriendListResp{}, nil
	}

	// 将用户信息构建为 map，key 为用户 ID，便于后续快速查找
	userRecords := make(map[string]*userclient.UserEntity, len(users.User))
	for i, _ := range users.User {
		userRecords[users.User[i].Id] = users.User[i]
	}

	// 组装响应数据：将好友关系和用户信息合并
	respList := make([]*types.Friends, 0, len(friends.List))
	for _, v := range friends.List {
		friend := &types.Friends{
			Id:        v.Id,
			FriendUid: v.FriendUid,
		}

		// 如果找到对应的用户信息，填充昵称和头像
		if u, ok := userRecords[v.FriendUid]; ok {
			friend.Nickname = u.Nickname
			friend.Avatar = u.Avatar
		}
		respList = append(respList, friend)
	}

	return &types.FriendListResp{
		List: respList,
	}, nil
}
