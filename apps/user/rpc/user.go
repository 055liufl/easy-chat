// =============================================================================
// User RPC 服务主程序
// =============================================================================
// 用户 RPC 服务的启动入口，包括:
//   - 配置加载（支持动态配置中心）
//   - gRPC 服务器初始化
//   - 服务注册
//   - 拦截器配置
//   - 系统根令牌初始化
//   - 服务启动和优雅关闭
//
// 配置来源:
//   支持本地配置文件和 ETCD 动态配置中心
//
// 业务场景:
//   提供用户相关的 RPC 服务，包括注册、登录、获取用户信息、查找用户等
//
// 服务特性:
//   - 支持配置热更新（通过 ETCD 配置中心）
//   - 支持优雅关闭（配置更新时自动重启）
//   - 支持服务反射（开发和测试模式）
//   - 统一的日志拦截器
//
// =============================================================================
package main

import (
	"flag"
	"fmt"
	"imooc.com/easy-chat/pkg/configserver"
	"imooc.com/easy-chat/pkg/interceptor/rpcserver"
	"sync"

	"imooc.com/easy-chat/apps/user/rpc/internal/config"
	"imooc.com/easy-chat/apps/user/rpc/internal/server"
	"imooc.com/easy-chat/apps/user/rpc/internal/svc"
	"imooc.com/easy-chat/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// configFile 配置文件路径（命令行参数）
var configFile = flag.String("f", "etc/dev/user.yaml", "the config file")

// grpcServer gRPC 服务器实例（用于优雅关闭）
var grpcServer *grpc.Server

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
	// 创建配置服务器并加载配置
	// 支持从 ETCD 配置中心加载配置，并监听配置变化
	err := configserver.NewConfigServer(*configFile, configserver.NewSail(&configserver.Config{
		ETCDEndpoints:  "192.168.117.24:3379",              // ETCD 服务器地址
		ProjectKey:     "98c6f2c2287f4c73cea3d40ae7ec3ff2", // 项目密钥
		Namespace:      "user",                             // 命名空间
		Configs:        "user-rpc.yaml",                    // 配置文件名
		ConfigFilePath: "./etc/conf",                       // 本地配置文件路径
		LogLevel:       "DEBUG",                            // 日志级别
	})).MustLoad(&c, func(bytes []byte) error {
		// 配置变化回调函数
		// 当配置发生变化时，优雅关闭当前服务并启动新服务
		var c config.Config
		configserver.LoadFromJsonBytes(bytes, &c)

		// 优雅关闭当前 gRPC 服务器
		grpcServer.GracefulStop()

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

// Run 运行 gRPC 服务器
// 负责服务器初始化、服务注册和启动
//
// 参数:
//   - c: 服务配置
//
// 运行流程:
//   1. 创建服务上下文（初始化依赖）
//   2. 设置系统根令牌
//   3. 创建 gRPC 服务器并注册服务
//   4. 添加日志拦截器
//   5. 启动服务器
//   6. 等待服务器停止（优雅关闭）
func Run(c config.Config) {
	// 创建服务上下文（初始化 MySQL、Redis 等依赖）
	ctx := svc.NewServiceContext(c)

	// 设置系统根令牌（用于系统内部服务调用）
	if err := ctx.SetRootToken(); err != nil {
		panic(err)
	}

	// 创建 gRPC 服务器并注册服务
	s := zrpc.MustNewServer(c.RpcServerConf, func(srv *grpc.Server) {
		grpcServer = srv

		// 注册用户 RPC 服务
		user.RegisterUserServer(grpcServer, server.NewUserServer(ctx))

		// 在开发和测试模式下启用服务反射（用于 grpcurl 等工具）
		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})

	// 添加日志拦截器（记录所有 RPC 调用）
	s.AddUnaryInterceptors(rpcserver.LogInterceptor)
	defer s.Stop()

	// 打印启动信息并启动服务器
	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
