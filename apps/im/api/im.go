// =============================================================================
// IM API 服务 - 即时通讯 HTTP 接口服务
// =============================================================================
// 提供即时通讯相关的 RESTful API 接口，包括:
//   - 聊天记录查询
//   - 会话管理
//   - 消息已读/未读状态
//
// 服务特性:
//   - 支持动态配置热更新（通过 ETCD 配置中心）
//   - JWT 认证保护
//   - 集成 IM RPC、User RPC、Social RPC 服务
//
// 配置来源:
//   - 本地配置文件（etc/dev/im.yaml）
//   - ETCD 配置中心（支持配置热更新）
//
// 启动方式:
//   go run im.go -f etc/dev/im.yaml
//
// =============================================================================
package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/proc"
	"imooc.com/easy-chat/pkg/configserver"
	"sync"

	"imooc.com/easy-chat/apps/im/api/internal/config"
	"imooc.com/easy-chat/apps/im/api/internal/handler"
	"imooc.com/easy-chat/apps/im/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

// configFile 配置文件路径，默认为 etc/dev/im.yaml
var configFile = flag.String("f", "etc/dev/im.yaml", "the config file")

// wg 等待组，用于管理多个服务实例的生命周期
var wg sync.WaitGroup

// main 主函数 - IM API 服务入口
// 负责加载配置、启动 HTTP 服务器
//
// 启动流程:
//  1. 解析命令行参数
//  2. 从配置中心加载配置（支持热更新）
//  3. 启动两个服务实例（初始启动 + 热更新后的新实例）
//  4. 等待所有服务实例退出
func main() {
	flag.Parse()

	var c config.Config

	// 创建配置服务器，从 ETCD 配置中心加载配置
	// 支持配置热更新，当配置变更时会自动重启服务
	err := configserver.NewConfigServer(*configFile, configserver.NewSail(&configserver.Config{
		ETCDEndpoints:  "192.168.117.24:3379",
		ProjectKey:     "98c6f2c2287f4c73cea3d40ae7ec3ff2",
		Namespace:      "im",
		Configs:        "im-api.yaml",
		ConfigFilePath: "./etc/conf",
		LogLevel:       "DEBUG",
	})).MustLoad(&c, func(bytes []byte) error {
		// 配置热更新回调函数
		// 当 ETCD 中的配置发生变化时，此函数会被调用
		var c config.Config
		configserver.LoadFromJsonBytes(bytes, &c)

		// 优雅关闭旧服务实例
		proc.WrapUp()

		// 启动新的服务实例
		wg.Add(1)
		go func(c config.Config) {
			defer wg.Done()

			Run(c)
		}(c)
		return nil
	})
	if err != nil {
		panic(err)
	}

	// 启动初始服务实例
	wg.Add(1)
	go func(c config.Config) {
		defer wg.Done()

		Run(c)
	}(c)

	// 等待所有服务实例退出
	wg.Wait()
}

// Run 运行 HTTP 服务器
// 创建并启动 go-zero REST 服务器，注册所有 API 路由
//
// 参数:
//   - c: 服务配置对象，包含端口、RPC 客户端配置等
//
// 服务生命周期:
//  1. 创建 REST 服务器
//  2. 初始化服务上下文（RPC 客户端连接）
//  3. 注册所有 HTTP 路由
//  4. 启动服务器监听
//  5. 优雅关闭（defer server.Stop()）
func Run(c config.Config) {
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}