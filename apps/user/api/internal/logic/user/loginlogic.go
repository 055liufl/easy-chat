// =============================================================================
// 用户登录业务逻辑
// =============================================================================
// 处理用户登录的业务逻辑，包括:
//   - 调用 User RPC 服务进行登录验证
//   - 数据转换（API 层类型 -> RPC 层类型）
//   - 更新用户在线状态到 Redis
//   - 返回登录结果（JWT token 和过期时间）
//
// 数据来源:
//   HTTP API 层传递的登录请求参数
//
// 业务场景:
//   用户通过手机号和密码进行登录，登录成功后返回访问令牌并更新在线状态
//
// 业务流程:
//   1. 接收 API 层的登录请求
//   2. 调用 User RPC 服务的 Login 方法进行验证
//   3. 将 RPC 响应转换为 API 响应格式
//   4. 将用户 ID 写入 Redis 在线用户集合
//   5. 返回登录结果
//
// =============================================================================
package user

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"imooc.com/easy-chat/apps/user/rpc/user"
	"imooc.com/easy-chat/pkg/constants"

	"imooc.com/easy-chat/apps/user/api/internal/svc"
	"imooc.com/easy-chat/apps/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// LoginLogic 用户登录业务逻辑结构
type LoginLogic struct {
	logx.Logger                    // 日志记录器
	ctx         context.Context    // 请求上下文
	svcCtx      *svc.ServiceContext // 服务上下文（包含配置和依赖）
}

// NewLoginLogic 创建用户登录业务逻辑实例
//
// 参数:
//   - ctx: 请求上下文
//   - svcCtx: 服务上下文
//
// 返回:
//   - *LoginLogic: 业务逻辑实例
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Login 执行用户登录业务逻辑
//
// 参数:
//   - req: 登录请求参数（包含手机号、密码）
//
// 返回:
//   - resp: 登录响应（包含 JWT token 和过期时间）
//   - err: 错误信息
//
// 业务流程:
//   1. 打印数据库配置（调试用）
//   2. 调用 User RPC 服务的 Login 方法进行验证
//   3. 使用 copier 将 RPC 响应转换为 API 响应格式
//   4. 将用户 ID 写入 Redis 在线用户 Hash 表（key: REDIS_ONLINE_USER, field: userId, value: "1"）
//   5. 返回登录结果
func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// todo: add your logic here and delete this line

	// 打印数据库配置（调试用）
	fmt.Println(l.svcCtx.Config.Database)

	// 调用 User RPC 服务进行登录验证
	loginResp, err := l.svcCtx.User.Login(l.ctx, &user.LoginReq{
		Phone:    req.Phone,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	// 将 RPC 响应转换为 API 响应格式
	var res types.LoginResp
	copier.Copy(&res, loginResp)

	// 处理登录的业务：将用户 ID 写入 Redis 在线用户集合
	// 使用 Hash 结构存储在线用户，key 为 REDIS_ONLINE_USER，field 为用户 ID，value 为 "1"
	l.svcCtx.Redis.HsetCtx(l.ctx, constants.REDIS_ONLINE_USER, loginResp.Id, "1")

	return &res, nil
}
