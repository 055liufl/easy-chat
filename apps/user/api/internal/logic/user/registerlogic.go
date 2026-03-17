// =============================================================================
// 用户注册业务逻辑
// =============================================================================
// 处理用户注册的业务逻辑，包括:
//   - 调用 User RPC 服务进行用户注册
//   - 数据转换（API 层类型 -> RPC 层类型）
//   - 返回注册结果（JWT token 和过期时间）
//
// 数据来源:
//   HTTP API 层传递的注册请求参数
//
// 业务场景:
//   用户通过手机号和密码进行注册，注册成功后返回访问令牌
//
// 业务流程:
//   1. 接收 API 层的注册请求
//   2. 调用 User RPC 服务的 Register 方法
//   3. 将 RPC 响应转换为 API 响应格式
//   4. 返回注册结果
//
// =============================================================================
package user

import (
	"context"
	"github.com/jinzhu/copier"
	"imooc.com/easy-chat/apps/user/rpc/user"

	"imooc.com/easy-chat/apps/user/api/internal/svc"
	"imooc.com/easy-chat/apps/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// RegisterLogic 用户注册业务逻辑结构
type RegisterLogic struct {
	logx.Logger                    // 日志记录器
	ctx         context.Context    // 请求上下文
	svcCtx      *svc.ServiceContext // 服务上下文（包含配置和依赖）
}

// NewRegisterLogic 创建用户注册业务逻辑实例
//
// 参数:
//   - ctx: 请求上下文
//   - svcCtx: 服务上下文
//
// 返回:
//   - *RegisterLogic: 业务逻辑实例
func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Register 执行用户注册业务逻辑
//
// 参数:
//   - req: 注册请求参数（包含手机号、密码、昵称、性别、头像）
//
// 返回:
//   - resp: 注册响应（包含 JWT token 和过期时间）
//   - err: 错误信息
//
// 业务流程:
//   1. 调用 User RPC 服务的 Register 方法进行注册
//   2. 使用 copier 将 RPC 响应转换为 API 响应格式
//   3. 返回注册结果
func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	// todo: add your logic here and delete this line

	// 调用 User RPC 服务进行注册
	registerResp, err := l.svcCtx.User.Register(l.ctx, &user.RegisterReq{
		Phone:    req.Phone,
		Nickname: req.Nickname,
		Password: req.Password,
		Avatar:   req.Avatar,
		Sex:      int32(req.Sex),
	})
	if err != nil {
		return nil, err
	}

	// 将 RPC 响应转换为 API 响应格式
	var res types.RegisterResp
	copier.Copy(&res, registerResp)

	return &res, nil
}
