// =============================================================================
// 群组列表查询逻辑
// =============================================================================
// 实现用户群组列表查询功能，包括:
//   - 查询用户加入的所有群组
//   - 根据群组 ID 批量查询群组详情
//   - 数据模型转换
//
// 业务流程:
//  1. 根据用户 ID 查询群成员记录，获取群组 ID 列表
//  2. 根据群组 ID 列表批量查询群组详情
//  3. 将数据库模型转换为 RPC 响应模型
//  4. 返回群组列表
//
// 数据流:
//   用户 ID -> 查询群成员 -> 提取群组 ID -> 批量查询群组 -> 模型转换 -> RPC 响应
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

// GroupListLogic 群组列表查询逻辑处理器
type GroupListLogic struct {
	ctx    context.Context        // 请求上下文
	svcCtx *svc.ServiceContext     // 服务上下文
	logx.Logger                    // 日志记录器
}

// NewGroupListLogic 创建群组列表查询逻辑处理器实例
//
// 参数:
//   - ctx: 请求上下文
//   - svcCtx: 服务上下文
//
// 返回:
//   - *GroupListLogic: 初始化完成的逻辑处理器实例
func NewGroupListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupListLogic {
	return &GroupListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GroupList 查询用户加入的群组列表
//
// 参数:
//   - in: 群组列表查询请求，包含用户 ID
//
// 返回:
//   - *social.GroupListResp: 群组列表响应，包含群组信息列表
//   - error: 错误信息（数据库查询失败等）
//
// 业务逻辑:
//  1. 根据用户 ID 查询群成员记录
//  2. 提取群组 ID 列表
//  3. 根据群组 ID 列表批量查询群组详情
//  4. 使用 copier 将数据库模型转换为 RPC 响应模型
func (l *GroupListLogic) GroupList(in *social.GroupListReq) (*social.GroupListResp, error) {
	// todo: add your logic here and delete this line

	// 查询用户加入的群组成员记录
	userGroup, err := l.svcCtx.GroupMembersModel.ListByUserId(l.ctx, in.UserId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "list group member err %v req %v", err, in.UserId)
	}
	if len(userGroup) == 0 {
		// 用户未加入任何群组，返回空列表
		return &social.GroupListResp{}, nil
	}

	// 提取群组 ID 列表
	ids := make([]string, 0, len(userGroup))
	for _, v := range userGroup {
		ids = append(ids, v.GroupId)
	}

	// 根据群组 ID 列表批量查询群组详情
	groups, err := l.svcCtx.GroupsModel.ListByGroupIds(l.ctx, ids)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "list group err %v req %v", err, ids)
	}

	// 将数据库模型转换为 RPC 响应模型
	var respList []*social.Groups
	copier.Copy(&respList, &groups)

	return &social.GroupListResp{
		List: respList,
	}, nil
}
