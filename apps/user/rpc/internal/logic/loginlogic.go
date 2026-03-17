// =============================================================================
// 用户登录业务逻辑（RPC）
// =============================================================================
// 处理用户登录的核心业务逻辑，包括:
//   - 验证手机号是否已注册
//   - 验证密码是否正确
//   - 生成 JWT 访问令牌（已注释，用于测试）
//
// 数据来源:
//   API 层通过 RPC 调用传递登录请求参数
//
// 业务场景:
//   用户通过手机号和密码进行登录，系统验证后返回访问令牌
//
// 业务流程:
//   1. 根据手机号查询用户信息
//   2. 如果用户不存在，返回 ErrPhoneNotRegister 错误
//   3. 验证密码是否正确（使用 bcrypt 验证）
//   4. 如果密码错误，返回 ErrUserPwdError 错误
//   5. 生成 JWT 访问令牌（当前已注释，用于测试）
//   6. 返回令牌和过期时间
//
// 注意:
//   当前代码中 JWT 生成部分已注释，返回测试错误
//
// =============================================================================
package logic

import (
	"context"
	"github.com/pkg/errors"
	"imooc.com/easy-chat/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
	"imooc.com/easy-chat/apps/user/models"
	"imooc.com/easy-chat/apps/user/rpc/internal/svc"
	"imooc.com/easy-chat/apps/user/rpc/user"
	"imooc.com/easy-chat/pkg/encrypt"
)

var (
	// ErrPhoneNotRegister 手机号未注册错误
	ErrPhoneNotRegister = xerr.New(xerr.SERVER_COMMON_ERROR, "手机号没有注册")
	// ErrUserPwdError 密码错误
	ErrUserPwdError = xerr.New(xerr.SERVER_COMMON_ERROR, "密码不正确")
)

// LoginLogic 用户登录业务逻辑结构
type LoginLogic struct {
	ctx    context.Context    // 请求上下文
	svcCtx *svc.ServiceContext // 服务上下文（包含配置和依赖）
	logx.Logger                // 日志记录器
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
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Login 执行用户登录业务逻辑
//
// 参数:
//   - in: 登录请求参数（包含手机号、密码）
//
// 返回:
//   - *user.LoginResp: 登录响应（包含用户 ID、JWT token 和过期时间）
//   - error: 错误信息
//
// 业务流程:
//   1. 根据手机号查询用户信息
//   2. 如果用户不存在，返回 ErrPhoneNotRegister 错误
//   3. 使用 bcrypt 验证密码是否正确
//   4. 如果密码错误，返回 ErrUserPwdError 错误
//   5. 生成 JWT 访问令牌（当前已注释，用于测试）
//   6. 返回用户 ID、令牌和过期时间
//
// 注意:
//   当前代码中 JWT 生成部分已注释，直接返回测试错误
func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {
	// todo: add your logic here and delete this line

	// 1. 验证用户是否注册，根据手机号码验证
	userEntity, err := l.svcCtx.UsersModel.FindByPhone(l.ctx, in.Phone)
	if err != nil {
		if err == models.ErrNotFound {
			return nil, errors.WithStack(ErrPhoneNotRegister)
		}
		return nil, errors.Wrapf(xerr.NewDBErr(), "find user by phone err %v , req %v", err, in.Phone)
	}

	// 密码验证（使用 bcrypt 验证密码哈希）
	if !encrypt.ValidatePasswordHash(in.Password, userEntity.Password.String) {
		return nil, errors.WithStack(ErrUserPwdError)
	}

	// 生成 JWT token（当前已注释，用于测试）
	//now := time.Now().Unix()
	//token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, l.svcCtx.Config.Jwt.AccessExpire,
	//	userEntity.Id)
	//if err != nil {
	//	return nil, errors.Wrapf(xerr.NewDBErr(), "ctxdata get jwt token err %v", err)
	//}

	// 返回测试错误（实际应返回下面注释的代码）
	return nil, errors.New("做测试")
	//return &user.LoginResp{
	//	Id:     userEntity.Id,
	//	Token:  token,
	//	Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
	//}, nil
}
