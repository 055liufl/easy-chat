// =============================================================================
// 好友在线状态 Logic - 查询好友在线状态业务逻辑
// =============================================================================
// 实现批量查询好友在线状态的核心业务逻辑
//
// 业务流程:
//   1. 从上下文获取当前用户 ID
//   2. 调用 Social RPC 获取好友列表
//   3. 从 Redis 缓存中查询好友的在线状态
//   4. 返回好友在线状态映射表
//
// 数据来源:
//   - Social RPC: 好友关系列表
//   - Redis: 在线用户缓存（REDIS_ONLINE_USER）
//
// =============================================================================
package friend

import (
	"context"
	"imooc.com/easy-chat/apps/social/rpc/socialclient"
	"imooc.com/easy-chat/pkg/constants"
	"imooc.com/easy-chat/pkg/ctxdata"

	"imooc.com/easy-chat/apps/social/api/internal/svc"
	"imooc.com/easy-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// FriendsOnlineLogic 好友在线状态业务逻辑处理器
type FriendsOnlineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewFriendsOnlineLogic 创建好友在线状态业务逻辑处理器实例
//
// 参数:
//   - ctx: 请求上下文，包含用户认证信息等
//   - svcCtx: 服务上下文，包含 RPC 客户端等依赖
//
// 返回:
//   - *FriendsOnlineLogic: 业务逻辑处理器实例
func NewFriendsOnlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendsOnlineLogic {
	return &FriendsOnlineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// FriendsOnline 查询好友在线状态
// 批量查询当前用户所有好友的在线状态
//
// 处理流程:
//   1. 从 JWT token 上下文中提取当前用户 ID
//   2. 调用 Social RPC 服务获取好友列表
//   3. 如果没有好友或查询失败，返回空结果
//   4. 提取所有好友的用户 ID
//   5. 从 Redis 缓存中查询在线用户集合（REDIS_ONLINE_USER）
//   6. 遍历好友列表，检查每个好友是否在线
//   7. 构建在线状态映射表（好友ID -> 在线状态）
//
// 参数:
//   - req: 好友在线状态请求参数（当前为空结构）
//
// 返回:
//   - resp: 好友在线状态响应，包含好友ID到在线状态的映射
//   - err: 错误信息
func (l *FriendsOnlineLogic) FriendsOnline(req *types.FriendsOnlineReq) (resp *types.FriendsOnlineResp, err error) {
	// 从上下文获取当前登录用户的 ID
	uid := ctxdata.GetUId(l.ctx)

	// 调用 Social RPC 获取好友列表
	friendList, err := l.svcCtx.Social.FriendList(l.ctx, &socialclient.FriendListReq{
		UserId: uid,
	})
	if err != nil || len(friendList.List) == 0 {
		return &types.FriendsOnlineResp{}, err
	}

	// 提取所有好友的用户 ID
	uids := make([]string, 0, len(friendList.List))
	for _, friend := range friendList.List {
		uids = append(uids, friend.FriendUid)
	}

	// 从 Redis 缓存中查询在线用户集合
	// REDIS_ONLINE_USER 是一个 Hash 结构，key 为用户 ID，value 为在线信息
	onlines, err := l.svcCtx.Redis.Hgetall(constants.REDIS_ONLINE_USER)
	if err != nil {
		return nil, err
	}

	// 构建在线状态映射表
	resOnLineList := make(map[string]bool, len(uids))
	for _, s := range uids {
		// 如果用户 ID 存在于在线用户集合中，则标记为在线
		if _, ok := onlines[s]; ok {
			resOnLineList[s] = true
		} else {
			resOnLineList[s] = false
		}
	}

	return &types.FriendsOnlineResp{
		OnlineList: resOnLineList,
	}, nil
}
