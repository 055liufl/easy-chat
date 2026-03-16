// =============================================================================
// Social API 配置 - 社交服务 API 配置定义
// =============================================================================
// 定义社交服务 API 层的所有配置项，包括:
//   - REST 服务配置（端口、超时等）
//   - Redis 缓存配置
//   - RPC 客户端配置（Social、User、IM 服务）
//   - JWT 认证配置
//
// 配置来源:
//   通过配置文件（如 social-api.yaml）加载
//
// 依赖服务:
//   - Social RPC: 社交关系服务（好友、群组管理）
//   - User RPC: 用户信息服务
//   - IM RPC: 即时通讯服务（在线状态查询）
//   - Redis: 缓存服务（幂等性、限流等）
//
// =============================================================================
package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

// Config 社交服务 API 配置结构
// 包含 REST 服务、缓存、RPC 客户端和认证相关的所有配置
type Config struct {
	rest.RestConf // REST 服务基础配置（端口、超时、日志等）

	Redisx redis.RedisConf // Redis 缓存配置，用于幂等性控制和限流

	SocialRpc zrpc.RpcClientConf // Social RPC 客户端配置，调用社交关系服务
	UserRpc   zrpc.RpcClientConf // User RPC 客户端配置，调用用户信息服务

	ImRpc zrpc.RpcClientConf // IM RPC 客户端配置，调用即时通讯服务

	// JwtAuth JWT 认证配置
	JwtAuth struct {
		AccessSecret string // JWT 签名密钥，用于验证 token 有效性
		//AccessExpire int64 // JWT 过期时间（秒），暂未使用
	}
}
