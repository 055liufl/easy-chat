// =============================================================================
// IM RPC 配置模块
// =============================================================================
// 定义即时通讯 RPC 服务的配置结构，包括:
//   - RPC 服务器配置（端口、超时、日志等）
//   - MongoDB 数据库连接配置
//
// 配置来源:
//   通过配置文件（YAML）或配置中心加载
//
// 使用场景:
//   服务启动时加载配置，初始化 RPC 服务器和数据库连接
//
// =============================================================================
package config

import "github.com/zeromicro/go-zero/zrpc"

// Config IM RPC 服务配置结构
type Config struct {
	// RPC 服务器配置（继承 go-zero 的 RPC 配置）
	// 包含监听地址、超时设置、日志配置等
	zrpc.RpcServerConf

	// Mongo MongoDB 数据库配置
	Mongo struct {
		Url string // MongoDB 连接地址，格式: mongodb://host:port
		Db  string // 数据库名称
	}
}
