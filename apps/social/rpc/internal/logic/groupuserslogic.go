// =============================================================================
// 群组成员查询逻辑
// =============================================================================
// 实现群组成员列表查询功能，包括:
//   - 根据群组 ID 查询所有成员
//   - 数据模型转换
//
// 业务流程:
//  1. 根据群组 ID 查询群成员列表
//  2. 将数据库模型转换为 RPC 响应模型
//  3. 返回成员列表
//
// 数据流:
//   群组 ID -> 数据库查询 -> 群成员列表 -> 模型转换 -> RPC 响应
// =============================================================================
package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"imooc.com/easy-chat/apps/social/rpc/internal/svc"
	"imooc.com/easy-chat/apps/social/rpc/social"
	"imooc.com/easy-chat/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

// GroupUsersLogic 群组成员查询逻辑处理器
type GroupUsersLogic struct {
	ctx    context.Context        // 请求上下文
	svcCtx *svc.ServiceContext     // 服务上下文
	logx.Logger                    // 日志记录器
}

// NewGroupUsersLogic 创建群组成员查询逻辑处理器实例
//
// 参数:
//   - ctx: 请求上下文
//   - svcCtx: 服务上下文
//
// 返回:
//   - *GroupUsersLogic: 初始化完成的逻辑处理器实例
func NewGroupUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupUsersLogic {
	return &GroupUsersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GroupUsers 查询群组的所有成员
//
// 参数:
//   - in: 群组成员查询请求，包含群组 ID
//
// 返回:
//   - *social.GroupUsersResp: 群组成员列表响应
//   - error: 错误信息（数据库查询失败等）
//
// 业务逻辑:
//  1. 根据群组 ID 查询数据库中的群成员记录
//  2. 使用 copier 将数据库模型转换为 RPC 响应模型
//  3. 返回成员列表
func (l *GroupUsersLogic) GroupUsers(in *social.GroupUsersReq) (*social.GroupUsersResp, error) {
	// todo: add your logic here and delete this line

	// 查询群组的所有成员
	groupMembers, err := l.svcCtx.GroupMembersModel.ListByGroupId(l.ctx, in.GroupId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "list group member err %v req %v", err, in.GroupId)
	}

	// 将数据库模型转换为 RPC 响应模型
	var respList []*social.GroupMembers
	copier.Copy(&respList, &groupMembers)

	//time.Sleep(5 * time.Second) // 延迟（已注释，可能用于测试）

	return &social.GroupUsersResp{
		List: respList,
	}, nil
}
