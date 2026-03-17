// =============================================================================
// 查找用户业务逻辑（RPC）
// =============================================================================
// 处理查找用户的业务逻辑，支持多种查询方式:
//   - 根据手机号查询单个用户
//   - 根据昵称模糊查询用户列表
//   - 根据用户 ID 列表批量查询
//
// 数据来源:
//   API 层或其他服务通过 RPC 调用传递查询参数
//
// 业务场景:
//   支持用户搜索、好友查找、批量获取用户信息等场景
//
// 业务流程:
//   1. 根据请求参数选择查询方式（手机号/昵称/ID 列表）
//   2. 调用对应的数据模型方法查询用户
//   3. 将数据库模型转换为 RPC 响应格式
//   4. 返回用户列表
//
// 注意:
//   当前代码包含测试逻辑，返回测试错误
//
// =============================================================================
package logic

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"imooc.com/easy-chat/apps/user/models"
	"imooc.com/easy-chat/apps/user/rpc/internal/svc"
	"imooc.com/easy-chat/apps/user/rpc/user"
	"imooc.com/easy-chat/pkg/xerr"
)

// FindUserLogic 查找用户业务逻辑结构
type FindUserLogic struct {
	ctx    context.Context    // 请求上下文
	svcCtx *svc.ServiceContext // 服务上下文（包含配置和依赖）
	logx.Logger                // 日志记录器
}

// NewFindUserLogic 创建查找用户业务逻辑实例
//
// 参数:
//   - ctx: 请求上下文
//   - svcCtx: 服务上下文
//
// 返回:
//   - *FindUserLogic: 业务逻辑实例
func NewFindUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserLogic {
	return &FindUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// cint 测试计数器（用于调试）
var cint int

// FindUser 执行查找用户业务逻辑
//
// 参数:
//   - in: 查找用户请求参数（支持手机号/昵称/ID 列表三种查询方式）
//
// 返回:
//   - *user.FindUserResp: 用户列表响应
//   - error: 错误信息
//
// 业务流程:
//   1. 根据请求参数选择查询方式:
//      - 如果提供了 Phone，根据手机号查询单个用户
//      - 如果提供了 Name，根据昵称模糊查询用户列表
//      - 如果提供了 Ids，根据 ID 列表批量查询用户
//   2. 调用对应的数据模型方法查询用户
//   3. 使用 copier 将数据库模型转换为 RPC 响应格式
//   4. 返回用户列表
//
// 注意:
//   当前代码包含测试逻辑，返回 DeadlineExceeded 错误用于测试
func (l *FindUserLogic) FindUser(in *user.FindUserReq) (*user.FindUserResp, error) {
	// todo: add your logic here and delete this line

	var (
		userEntitys []*models.Users
		err         error
	)

	// 根据请求参数选择查询方式
	if in.Phone != "" {
		// 根据手机号查询单个用户
		userEntity, err := l.svcCtx.UsersModel.FindByPhone(l.ctx, in.Phone)
		if err == nil {
			userEntitys = append(userEntitys, userEntity)
		}
	} else if in.Name != "" {
		// 根据昵称模糊查询用户列表
		userEntitys, err = l.svcCtx.UsersModel.ListByName(l.ctx, in.Name)
	} else if len(in.Ids) > 0 {
		// 根据 ID 列表批量查询用户
		userEntitys, err = l.svcCtx.UsersModel.ListByIds(l.ctx, in.Ids)
	}

	if err != nil {
		return nil, err
	}

	// 将数据库模型转换为 RPC 响应格式
	var resp []*user.UserEntity
	copier.Copy(&resp, &userEntitys)

	// 测试计数器（用于调试）
	cint++
	fmt.Println("------ finduserlogic cint ", cint)

	// 返回用户列表（当前返回测试错误）
	return &user.FindUserResp{
		User: resp,
	}, errors.WithStack(xerr.New(int(codes.DeadlineExceeded), "测试"))
}
