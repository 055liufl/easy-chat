// =============================================================================
// Social RPC 配置模块
// =============================================================================
// 定义社交服务 RPC 的配置结构，包括:
//   - RPC 服务器配置（端口、超时等）
//   - MySQL 数据库连接配置
//   - Redis 缓存配置
//
// 配置来源:
//   通过配置文件（YAML）或配置中心（ETCD）加载
//
// 使用场景:
//   在服务启动时加载配置，初始化数据库连接、缓存等资源
// =============================================================================
package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

// Config 社交服务 RPC 配置结构
type Config struct {
	// RPC 服务器配置（继承 go-zero 的 RPC 服务器配置）
	zrpc.RpcServerConf

	// MySQL 数据库配置
	Mysql struct {
		DataSource string // 数据库连接字符串，格式: user:password@tcp(host:port)/dbname
	}

	// Redis 缓存配置
	Cache cache.CacheConf // 缓存配置，用于数据库查询结果缓存
}
