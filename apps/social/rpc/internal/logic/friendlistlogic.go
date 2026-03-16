// =============================================================================
// 好友列表查询逻辑
// =============================================================================
// 实现好友列表查询功能，包括:
//   - 根据用户 ID 查询好友列表
//   - 数据模型转换（数据库模型 -> RPC 响应模型）
//
// 业务流程:
//  1. 根据用户 ID 查询数据库中的好友关系记录
//  2. 将数据库模型转换为 RPC 响应模型
//  3. 返回好友列表
//
// 数据流:
//   用户 ID -> 数据库查询 -> 好友关系列表 -> 模型转换 -> RPC 响应
// =============================================================================
package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"imooc.com/easy-chat/pkg/xerr"

	"imooc.com/easy-chat/apps/social/rpc/internal/svc"
	"imooc.com/easy-chat/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

// FriendListLogic 好友列表查询逻辑处理器
type FriendListLogic struct {
	ctx    context.Context        // 请求上下文
	svcCtx *svc.ServiceContext     // 服务上下文（包含数据模型等依赖）
	logx.Logger                    // 日志记录器
}

// NewFriendListLogic 创建好友列表查询逻辑处理器实例
//
// 参数:
//   - ctx: 请求上下文，用于超时控制和链路追踪
//   - svcCtx: 服务上下文，提供数据模型访问
//
// 返回:
//   - *FriendListLogic: 初始化完成的逻辑处理器实例
func NewFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendListLogic {
	return &FriendListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// FriendList 查询用户的好友列表
//
// 参数:
//   - in: 好友列表查询请求，包含用户 ID
//
// 返回:
//   - *social.FriendListResp: 好友列表响应，包含好友信息列表
//   - error: 错误信息（数据库查询失败等）
//
// 业务逻辑:
//  1. 根据用户 ID 查询数据库中的好友关系记录
//  2. 使用 copier 将数据库模型转换为 RPC 响应模型
//  3. 返回好友列表
func (l *FriendListLogic) FriendList(in *social.FriendListReq) (*social.FriendListResp, error) {
	// todo: add your logic here and delete this line

	// 查询用户的好友列表
	friendsList, err := l.svcCtx.FriendsModel.ListByUserid(l.ctx, in.UserId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "list friend by uid err %v req %v ", err,
			in.UserId)
	}

	// 将数据库模型转换为 RPC 响应模型
	var respList []*social.Friends
	copier.Copy(&respList, &friendsList)

	return &social.FriendListResp{
		List: respList,
	}, nil
}
