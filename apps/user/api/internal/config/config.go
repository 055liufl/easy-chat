// =============================================================================
// User API 配置模块
// =============================================================================
// 定义用户 API 服务的配置结构，包括:
//   - REST 服务配置（端口、超时等）
//   - 数据库连接配置
//   - Redis 缓存配置
//   - RPC 客户端配置（用于调用 User RPC 服务）
//   - JWT 认证配置（用于生成和验证访问令牌）
//
// 配置来源:
//   通过配置文件（YAML）加载，支持动态配置中心（如 ETCD）
//
// 业务场景:
//   为用户 API 服务提供统一的配置管理，包括服务端口、数据库连接、
//   缓存配置、RPC 调用配置和 JWT 认证密钥等
//
// =============================================================================
package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

// Config 用户 API 服务配置结构
type Config struct {
	rest.RestConf // REST 服务配置（继承 go-zero 的 REST 配置）
	Database      string // 数据库名称（预留字段）

	Redisx redis.RedisConf // Redis 配置（用于缓存和在线用户状态管理）

	UserRpc zrpc.RpcClientConf // User RPC 客户端配置（用于调用用户 RPC 服务）

	// JwtAuth JWT 认证配置
	JwtAuth struct {
		AccessSecret string // JWT 签名密钥（用于生成和验证 token）
		//AccessExpire int64  // JWT 过期时间（秒）- 已注释
	}
}
