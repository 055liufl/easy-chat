// =============================================================================
// User RPC 服务上下文
// =============================================================================
// 定义用户 RPC 服务的上下文结构，包括:
//   - 配置信息（Config）
//   - Redis 客户端（用于在线状态和系统令牌管理）
//   - 用户数据模型（UsersModel，用于数据库操作）
//
// 数据来源:
//   在服务启动时初始化，从配置文件加载配置并创建各种客户端连接
//
// 业务场景:
//   为所有 RPC Logic 提供统一的依赖注入，包括配置、数据库模型和
//   Redis 客户端等
//
// =============================================================================
package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"imooc.com/easy-chat/apps/user/models"
	"imooc.com/easy-chat/apps/user/rpc/internal/config"
	"imooc.com/easy-chat/pkg/constants"
	"imooc.com/easy-chat/pkg/ctxdata"
	"time"
)

// ServiceContext 服务上下文结构
// 包含服务运行所需的所有依赖
type ServiceContext struct {
	Config config.Config // 服务配置

	*redis.Redis      // Redis 客户端（用于在线用户状态和系统令牌管理）
	models.UsersModel // 用户数据模型（用于用户表的 CRUD 操作）
}

// NewServiceContext 创建服务上下文实例
// 初始化所有依赖的客户端连接和数据模型
//
// 参数:
//   - c: 服务配置
//
// 返回:
//   - *ServiceContext: 初始化完成的服务上下文实例
//
// 初始化流程:
//   1. 创建 MySQL 数据库连接
//   2. 创建 Redis 客户端连接
//   3. 创建用户数据模型（带缓存）
//   4. 返回服务上下文实例
func NewServiceContext(c config.Config) *ServiceContext {
	// 创建 MySQL 数据库连接
	sqlConn := sqlx.NewMysql(c.Mysql.DataSource)

	return &ServiceContext{
		Config:     c,
		Redis:      redis.MustNewRedis(c.Redisx),
		UsersModel: models.NewUsersModel(sqlConn, c.Cache),
	}
}

// SetRootToken 设置系统根令牌
// 生成系统根用户的 JWT token 并存储到 Redis
//
// 返回:
//   - error: 错误信息
//
// 业务场景:
//   系统启动时生成一个长期有效的根令牌，用于系统内部服务调用
//
// 业务流程:
//   1. 使用系统根用户 ID（SYSTEM_ROOT_UID）生成 JWT token
//   2. 设置超长过期时间（999999999 秒，约 31 年）
//   3. 将 token 存储到 Redis（key: REDIS_SYSTEM_ROOT_TOKEN）
func (svc *ServiceContext) SetRootToken() error {
	// 生成 JWT token
	// 参数: 密钥、当前时间、过期时间、用户 ID
	systemToken, err := ctxdata.GetJwtToken(svc.Config.Jwt.AccessSecret, time.Now().Unix(), 999999999, constants.SYSTEM_ROOT_UID)
	if err != nil {
		return err
	}
	// 将系统根令牌写入到 Redis
	return svc.Redis.Set(constants.REDIS_SYSTEM_ROOT_TOKEN, systemToken)
}
