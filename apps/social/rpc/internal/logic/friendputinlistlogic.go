// =============================================================================
// 好友申请列表查询逻辑
// =============================================================================
// 实现好友申请列表查询功能，包括:
//   - 查询用户收到的未处理好友申请
//   - 数据模型转换
//
// 业务流程:
//  1. 根据用户 ID 查询未处理的好友申请列表
//  2. 将数据库模型转换为 RPC 响应模型
//  3. 返回申请列表
//
// 数据流:
//   用户 ID -> 数据库查询 -> 未处理申请列表 -> 模型转换 -> RPC 响应
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

// FriendPutInListLogic 好友申请列表查询逻辑处理器
type FriendPutInListLogic struct {
	ctx    context.Context        // 请求上下文
	svcCtx *svc.ServiceContext     // 服务上下文
	logx.Logger                    // 日志记录器
}

// NewFriendPutInListLogic 创建好友申请列表查询逻辑处理器实例
//
// 参数:
//   - ctx: 请求上下文
//   - svcCtx: 服务上下文
//
// 返回:
//   - *FriendPutInListLogic: 初始化完成的逻辑处理器实例
func NewFriendPutInListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInListLogic {
	return &FriendPutInListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// FriendPutInList 查询用户收到的未处理好友申请列表
//
// 参数:
//   - in: 好友申请列表查询请求，包含用户 ID
//
// 返回:
//   - *social.FriendPutInListResp: 好友申请列表响应
//   - error: 错误信息（数据库查询失败等）
//
// 业务逻辑:
//  1. 根据用户 ID 查询数据库中未处理的好友申请记录
//  2. 使用 copier 将数据库模型转换为 RPC 响应模型
//  3. 返回申请列表
func (l *FriendPutInListLogic) FriendPutInList(in *social.FriendPutInListReq) (*social.FriendPutInListResp, error) {
	// todo: add your logic here and delete this line

	// 查询用户收到的未处理好友申请列表
	friendReqList, err := l.svcCtx.FriendRequestsModel.ListNoHandler(l.ctx, in.UserId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find list friend req err %v req %v", err, in.UserId)
	}

	// 将数据库模型转换为 RPC 响应模型
	var resp []*social.FriendRequests
	copier.Copy(&resp, &friendReqList)

	return &social.FriendPutInListResp{
		List: resp,
	}, nil

}
