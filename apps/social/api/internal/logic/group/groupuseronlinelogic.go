// =============================================================================
// 群成员在线状态 Logic - 查询群成员在线状态业务逻辑
// =============================================================================
// 实现批量查询群成员在线状态的核心业务逻辑
//
// 业务流程:
//   1. 调用 Social RPC 获取群成员列表
//   2. 从 Redis 缓存中查询成员的在线状态
//   3. 返回成员在线状态映射表
//
// 数据来源:
//   - Social RPC: 群成员列表
//   - Redis: 在线用户缓存（REDIS_ONLINE_USER）
//
// =============================================================================
package group

import (
	"context"
	"imooc.com/easy-chat/apps/social/rpc/socialclient"
	"imooc.com/easy-chat/pkg/constants"

	"imooc.com/easy-chat/apps/social/api/internal/svc"
	"imooc.com/easy-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// GroupUserOnlineLogic 群成员在线状态业务逻辑处理器
type GroupUserOnlineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGroupUserOnlineLogic 创建群成员在线状态业务逻辑处理器实例
//
// 参数:
//   - ctx: 请求上下文，包含用户认证信息等
//   - svcCtx: 服务上下文，包含 RPC 客户端等依赖
//
// 返回:
//   - *GroupUserOnlineLogic: 业务逻辑处理器实例
func NewGroupUserOnlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupUserOnlineLogic {
	return &GroupUserOnlineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GroupUserOnline 查询群成员在线状态
// 批量查询指定群组所有成员的在线状态
//
// 处理流程:
//   1. 调用 Social RPC 服务获取群成员列表
//   2. 如果没有成员或查询失败，返回空结果
//   3. 提取所有成员的用户 ID
//   4. 从 Redis 缓存中查询在线用户集合（REDIS_ONLINE_USER）
//   5. 遍历成员列表，检查每个成员是否在线
//   6. 构建在线状态映射表（用户ID -> 在线状态）
//
// 参数:
//   - req: 群成员在线状态请求参数
//     - GroupId: 群组 ID
//
// 返回:
//   - resp: 群成员在线状态响应，包含用户ID到在线状态的映射
//   - err: 错误信息
func (l *GroupUserOnlineLogic) GroupUserOnline(req *types.GroupUserOnlineReq) (resp *types.GroupUserOnlineResp, err error) {
	// 调用 Social RPC 获取群成员列表
	groupUsers, err := l.svcCtx.Social.GroupUsers(l.ctx, &socialclient.GroupUsersReq{
		GroupId: req.GroupId,
	})
	if err != nil || len(groupUsers.List) == 0 {
		return &types.GroupUserOnlineResp{}, err
	}

	// 提取所有成员的用户 ID
	uids := make([]string, 0, len(groupUsers.List))
	for _, group := range groupUsers.List {
		uids = append(uids, group.UserId)
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

	return &types.GroupUserOnlineResp{
		OnlineList: resOnLineList,
	}, nil
}
