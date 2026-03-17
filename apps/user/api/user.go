// =============================================================================
// User API 服务主程序
// =============================================================================
// 用户 API 服务的启动入口，包括:
//   - 配置加载（支持动态配置中心）
//   - REST 服务器初始化
//   - 路由注册
//   - 错误处理器和响应处理器配置
//   - 服务启动和优雅关闭
//
// 配置来源:
//   支持本地配置文件和 ETCD 动态配置中心
//
// 业务场景:
//   提供用户相关的 HTTP API 服务，包括注册、登录、获取用户信息等
//
// 服务特性:
//   - 支持配置热更新（通过 ETCD 配置中心）
//   - 支持优雅关闭（配置更新时自动重启）
//   - 统一的错误处理和响应格式
//
// =============================================================================
package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/proc"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
	"imooc.com/easy-chat/apps/user/api/internal/config"
	"imooc.com/easy-chat/apps/user/api/internal/handler"
	"imooc.com/easy-chat/apps/user/api/internal/svc"
	"imooc.com/easy-chat/pkg/configserver"
	"imooc.com/easy-chat/pkg/resultx"
	"sync"
)

// configFile 配置文件路径（命令行参数）
var configFile = flag.String("f", "etc/dev/user.yaml", "the config file")

// wg 等待组（用于等待所有 goroutine 完成）
var wg sync.WaitGroup

// main 主函数
// 负责配置加载和服务启动
//
// 启动流程:
//   1. 解析命令行参数
//   2. 加载配置（支持 ETCD 动态配置中心）
//   3. 启动服务（首次启动）
//   4. 监听配置变化并自动重启服务
//   5. 等待所有服务完成
func main() {
	flag.Parse()

	var c config.Config
	//conf.MustLoad(*configFile, &c)

	// 创建配置服务器并加载配置
	// 支持从 ETCD 配置中心加载配置，并监听配置变化
	err := configserver.NewConfigServer(*configFile, configserver.NewSail(&configserver.Config{
		ETCDEndpoints:  "192.168.117.24:3379",              // ETCD 服务器地址
		ProjectKey:     "98c6f2c2287f4c73cea3d40ae7ec3ff2", // 项目密钥
		Namespace:      "user",                             // 命名空间
		Configs:        "user-api.yaml",                    // 配置文件名
		ConfigFilePath: "./etc/conf",                       // 本地配置文件路径
		LogLevel:       "DEBUG",                            // 日志级别
	})).MustLoad(&c, func(bytes []byte) error {
		// 配置变化回调函数
		// 当配置发生变化时，优雅关闭当前服务并启动新服务
		var c config.Config
		configserver.LoadFromJsonBytes(bytes, &c)

		// 优雅关闭当前服务
		proc.WrapUp()

		// 启动新服务
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

	// 首次启动服务
	wg.Add(1)
	go func(c config.Config) {
		defer wg.Done()

		Run(c)
	}(c)

	// 等待所有服务完成
	wg.Wait()
}

// Run 运行 REST 服务器
// 负责服务器初始化、路由注册和启动
//
// 参数:
//   - c: 服务配置
//
// 运行流程:
//   1. 创建 REST 服务器实例
//   2. 创建服务上下文（初始化依赖）
//   3. 注册 HTTP 路由
//   4. 配置错误处理器和响应处理器
//   5. 启动服务器
//   6. 等待服务器停止（优雅关闭）
func Run(c config.Config) {
	// 创建 REST 服务器实例
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// 创建服务上下文（初始化 Redis、RPC 客户端等依赖）
	ctx := svc.NewServiceContext(c)
	// 注册 HTTP 路由
	handler.RegisterHandlers(server, ctx)

	// 配置统一的错误处理器（返回标准错误格式）
	httpx.SetErrorHandlerCtx(resultx.ErrHandler(c.Name))
	// 配置统一的成功响应处理器（返回标准成功格式）
	httpx.SetOkHandler(resultx.OkHandler)

	// 打印启动信息并启动服务器
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
