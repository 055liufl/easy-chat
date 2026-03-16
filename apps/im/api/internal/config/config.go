// =============================================================================
// IM API 配置 - 服务配置定义
// =============================================================================
// 定义 IM API 服务的配置结构，包括:
//   - REST 服务配置（端口、超时等）
//   - RPC 客户端配置（IM、User、Social 服务）
//   - JWT 认证配置
//
// 配置文件格式: YAML
// 配置来源: 本地文件 + ETCD 配置中心
//
// =============================================================================
package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

// Config IM API 服务配置结构
type Config struct {
	// RestConf REST 服务配置（继承 go-zero 的 REST 配置）
	// 包含: Host, Port, Timeout, MaxConns 等
	rest.RestConf

	// ImRpc IM RPC 服务客户端配置
	// 用于调用 IM RPC 服务（聊天记录、会话管理等）
	ImRpc zrpc.RpcClientConf

	// UserRpc User RPC 服务客户端配置
	// 用于调用 User RPC 服务（用户信息查询等）
	UserRpc zrpc.RpcClientConf

	// SocialRpc Social RPC 服务客户端配置
	// 用于调用 Social RPC 服务（群组信息、好友关系等）
	SocialRpc zrpc.RpcClientConf

	// JwtAuth JWT 认证配置
	JwtAuth struct {
		// AccessSecret JWT 签名密钥
		// 用于验证客户端请求中的 JWT Token
		AccessSecret string
	}
}
