// =============================================================================
// User RPC 配置模块
// =============================================================================
// 定义用户 RPC 服务的配置结构，包括:
//   - RPC 服务器配置（端口、超时等）
//   - MySQL 数据库配置
//   - 缓存配置（用于数据库查询缓存）
//   - Redis 配置（用于在线状态管理）
//   - JWT 配置（用于生成访问令牌）
//
// 配置来源:
//   通过配置文件（YAML）加载，支持动态配置中心（如 ETCD）
//
// 业务场景:
//   为用户 RPC 服务提供统一的配置管理，包括服务端口、数据库连接、
//   缓存配置和 JWT 密钥等
//
// =============================================================================
package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

// Config 用户 RPC 服务配置结构
type Config struct {
	zrpc.RpcServerConf // RPC 服务器配置（继承 go-zero 的 RPC 配置）

	// Mysql MySQL 数据库配置
	Mysql struct {
		DataSource string // 数据库连接字符串（格式: user:password@tcp(host:port)/dbname）
	}

	Cache cache.CacheConf // 缓存配置（用于数据库查询结果缓存）

	Redisx redis.RedisConf // Redis 配置（用于在线用户状态管理和系统令牌存储）

	// Jwt JWT 配置
	Jwt struct {
		AccessSecret string // JWT 签名密钥（用于生成和验证 token）
		AccessExpire int64  // JWT 过期时间（秒）
	}
}
