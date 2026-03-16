// =============================================================================
// 群组申请列表查询逻辑
// =============================================================================
// 实现群组申请列表查询功能，包括:
//   - 查询群组收到的未处理入群申请
//   - 数据模型转换
//
// 业务流程:
//  1. 根据群组 ID 查询未处理的入群申请列表
//  2. 将数据库模型转换为 RPC 响应模型
//  3. 返回申请列表
//
// 数据流:
//   群组 ID -> 数据库查询 -> 未处理申请列表 -> 模型转换 -> RPC 响应
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

// GroupPutinListLogic 群组申请列表查询逻辑处理器
type GroupPutinListLogic struct {
	ctx    context.Context        // 请求上下文
	svcCtx *svc.ServiceContext     // 服务上下文
	logx.Logger                    // 日志记录器
}

// NewGroupPutinListLogic 创建群组申请列表查询逻辑处理器实例
//
// 参数:
//   - ctx: 请求上下文
//   - svcCtx: 服务上下文
//
// 返回:
//   - *GroupPutinListLogic: 初始化完成的逻辑处理器实例
func NewGroupPutinListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutinListLogic {
	return &GroupPutinListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GroupPutinList 查询群组收到的未处理入群申请列表
//
// 参数:
//   - in: 群组申请列表查询请求，包含群组 ID
//
// 返回:
//   - *social.GroupPutinListResp: 群组申请列表响应
//   - error: 错误信息（数据库查询失败等）
//
// 业务逻辑:
//  1. 根据群组 ID 查询数据库中未处理的入群申请记录
//  2. 使用 copier 将数据库模型转换为 RPC 响应模型
//  3. 返回申请列表
func (l *GroupPutinListLogic) GroupPutinList(in *social.GroupPutinListReq) (*social.GroupPutinListResp, error) {
	// todo: add your logic here and delete this line

	// 查询群组收到的未处理入群申请列表
	groupReqs, err := l.svcCtx.GroupRequestsModel.ListNoHandler(l.ctx, in.GroupId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "list group req err %v req %v", err, in.GroupId)
	}

	// 将数据库模型转换为 RPC 响应模型
	var respList []*social.GroupRequests
	copier.Copy(&respList, groupReqs)

	return &social.GroupPutinListResp{
		List: respList,
	}, nil
}
