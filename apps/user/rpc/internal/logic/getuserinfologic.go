// =============================================================================
// 获取用户信息业务逻辑（RPC）
// =============================================================================
// 处理获取用户详细信息的业务逻辑，包括:
//   - 根据用户 ID 查询用户信息
//   - 数据转换（数据库模型 -> RPC 响应）
//   - 返回用户详细信息
//
// 数据来源:
//   API 层通过 RPC 调用传递用户 ID
//
// 业务场景:
//   根据用户 ID 获取用户的详细信息（昵称、头像、性别等）
//
// 业务流程:
//   1. 根据用户 ID 查询数据库
//   2. 如果用户不存在，返回 ErrUserNotFound 错误
//   3. 将数据库模型转换为 RPC 响应格式
//   4. 返回用户详细信息
//
// =============================================================================
package logic

import (
	"context"
	"errors"
	"github.com/jinzhu/copier"
	"imooc.com/easy-chat/apps/user/models"

	"imooc.com/easy-chat/apps/user/rpc/internal/svc"
	"imooc.com/easy-chat/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

// ErrUserNotFound 用户不存在错误
var ErrUserNotFound = errors.New("这个用户没有")

// GetUserInfoLogic 获取用户信息业务逻辑结构
type GetUserInfoLogic struct {
	ctx    context.Context    // 请求上下文
	svcCtx *svc.ServiceContext // 服务上下文（包含配置和依赖）
	logx.Logger                // 日志记录器
}

// NewGetUserInfoLogic 创建获取用户信息业务逻辑实例
//
// 参数:
//   - ctx: 请求上下文
//   - svcCtx: 服务上下文
//
// 返回:
//   - *GetUserInfoLogic: 业务逻辑实例
func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetUserInfo 执行获取用户信息业务逻辑
//
// 参数:
//   - in: 获取用户信息请求参数（包含用户 ID）
//
// 返回:
//   - *user.GetUserInfoResp: 用户信息响应（包含用户详细信息）
//   - error: 错误信息
//
// 业务流程:
//   1. 根据用户 ID 从数据库查询用户信息
//   2. 如果用户不存在（ErrNotFound），返回 ErrUserNotFound 错误
//   3. 使用 copier 将数据库模型转换为 RPC 响应格式
//   4. 返回用户详细信息
func (l *GetUserInfoLogic) GetUserInfo(in *user.GetUserInfoReq) (*user.GetUserInfoResp, error) {
	// todo: add your logic here and delete this line

	// 根据用户 ID 查询用户信息
	userEntiy, err := l.svcCtx.UsersModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == models.ErrNotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// 将数据库模型转换为 RPC 响应格式
	var resp user.UserEntity
	copier.Copy(&resp, userEntiy)

	return &user.GetUserInfoResp{
		User: &resp,
	}, nil
}
