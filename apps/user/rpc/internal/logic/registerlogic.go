// =============================================================================
// 用户注册业务逻辑（RPC）
// =============================================================================
// 处理用户注册的核心业务逻辑，包括:
//   - 验证手机号是否已注册
//   - 密码加密处理
//   - 生成唯一用户 ID
//   - 创建用户记录
//   - 生成 JWT 访问令牌
//
// 数据来源:
//   API 层通过 RPC 调用传递注册请求参数
//
// 业务场景:
//   用户通过手机号和密码进行注册，系统验证后创建用户并返回访问令牌
//
// 业务流程:
//   1. 根据手机号查询用户是否已注册
//   2. 如果已注册，返回错误
//   3. 生成唯一用户 ID（基于数据库连接字符串）
//   4. 对密码进行 bcrypt 加密
//   5. 插入用户记录到数据库
//   6. 生成 JWT 访问令牌
//   7. 返回令牌和过期时间
//
// =============================================================================
package logic

import (
	"context"
	"database/sql"
	"errors"
	"imooc.com/easy-chat/apps/user/models"
	"imooc.com/easy-chat/pkg/ctxdata"
	"imooc.com/easy-chat/pkg/encrypt"
	"imooc.com/easy-chat/pkg/wuid"
	"time"

	"imooc.com/easy-chat/apps/user/rpc/internal/svc"
	"imooc.com/easy-chat/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	// ErrPhoneIsRegister 手机号已注册错误
	ErrPhoneIsRegister = errors.New("手机号已经注册过")
)

// RegisterLogic 用户注册业务逻辑结构
type RegisterLogic struct {
	ctx    context.Context    // 请求上下文
	svcCtx *svc.ServiceContext // 服务上下文（包含配置和依赖）
	logx.Logger                // 日志记录器
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
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Register 执行用户注册业务逻辑
//
// 参数:
//   - in: 注册请求参数（包含手机号、密码、昵称、性别、头像）
//
// 返回:
//   - *user.RegisterResp: 注册响应（包含 JWT token 和过期时间）
//   - error: 错误信息
//
// 业务流程:
//   1. 根据手机号查询用户是否已注册
//   2. 如果已注册，返回 ErrPhoneIsRegister 错误
//   3. 生成唯一用户 ID（使用 wuid 算法）
//   4. 如果提供了密码，使用 bcrypt 进行加密
//   5. 插入用户记录到数据库
//   6. 生成 JWT 访问令牌（包含用户 ID）
//   7. 返回令牌和过期时间
func (l *RegisterLogic) Register(in *user.RegisterReq) (*user.RegisterResp, error) {
	// todo: add your logic here and delete this line

	// 1. 验证用户是否注册，根据手机号码验证
	userEntity, err := l.svcCtx.UsersModel.FindByPhone(l.ctx, in.Phone)
	if err != nil && err != models.ErrNotFound {
		return nil, err
	}

	// 如果用户已存在，返回错误
	if userEntity != nil {
		return nil, ErrPhoneIsRegister
	}

	// 定义用户数据
	userEntity = &models.Users{
		Id:       wuid.GenUid(l.svcCtx.Config.Mysql.DataSource), // 生成唯一用户 ID
		Avatar:   in.Avatar,
		Nickname: in.Nickname,
		Phone:    in.Phone,
		Sex: sql.NullInt64{
			Int64: int64(in.Sex),
			Valid: true,
		},
	}

	// 如果提供了密码，进行加密处理
	if len(in.Password) > 0 {
		genPassword, err := encrypt.GenPasswordHash([]byte(in.Password))
		if err != nil {
			return nil, err
		}
		userEntity.Password = sql.NullString{
			String: string(genPassword),
			Valid:  true,
		}
	}

	// 插入用户记录到数据库
	_, err = l.svcCtx.UsersModel.Insert(l.ctx, userEntity)
	if err != nil {
		return nil, err
	}

	// 生成 JWT token
	now := time.Now().Unix()
	token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, l.svcCtx.Config.Jwt.AccessExpire,
		userEntity.Id)
	if err != nil {
		return nil, err
	}

	return &user.RegisterResp{
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
	}, nil
}
