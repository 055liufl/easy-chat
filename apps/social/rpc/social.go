// =============================================================================
// Social RPC 服务主程序
// =============================================================================
// 社交服务 RPC 主程序，负责:
//   - 加载配置（支持配置文件和配置中心）
//   - 启动 gRPC 服务器
//   - 注册服务接口
//   - 添加拦截器（幂等性等）
//   - 支持配置热更新和服务优雅重启
//
// 启动流程:
//  1. 解析命令行参数，获取配置文件路径
//  2. 加载配置（支持 ETCD 配置中心）
//  3. 创建服务上下文（初始化数据库连接等）
//  4. 创建 gRPC 服务器并注册服务
//  5. 添加拦截器（幂等性拦截器）
//  6. 启动服务监听
//
// 配置热更新:
//   监听配置中心变化，配置更新时优雅重启服务
//
// 服务注册:
//   在开发/测试模式下启用 gRPC 反射，方便调试
// =============================================================================
package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"imooc.com/easy-chat/apps/social/rpc/internal/config"
	"imooc.com/easy-chat/apps/social/rpc/internal/server"
	"imooc.com/easy-chat/apps/social/rpc/internal/svc"
	"imooc.com/easy-chat/apps/social/rpc/social"
	"imooc.com/easy-chat/pkg/configserver"
	"imooc.com/easy-chat/pkg/interceptor"
	"sync"
)

var configFile = flag.String("f", "etc/dev/social.yaml", "the config file") // 配置文件路径
var grpcServer *grpc.Server                                                  // gRPC 服务器实例
var wg sync.WaitGroup                                                        // 等待组，用于协程同步

// main 主函数
// 负责解析命令行参数、加载配置、启动服务
func main() {
	flag.Parse()

	var c config.Config
	// 加载配置，支持配置中心（ETCD）和配置热更新
	err := configserver.NewConfigServer(*configFile, configserver.NewSail(&configserver.Config{
		ETCDEndpoints:  "192.168.117.24:3379",      // ETCD 地址
		ProjectKey:     "98c6f2c2287f4c73cea3d40ae7ec3ff2", // 项目密钥
		Namespace:      "social",                    // 命名空间
		Configs:        "social-rpc.yaml",           // 配置文件名
		ConfigFilePath: "./etc/conf",                // 配置文件本地缓存路径
		LogLevel:       "DEBUG",                     // 日志级别
	})).MustLoad(&c, func(bytes []byte) error {
		// 配置更新回调函数
		var c config.Config
		configserver.LoadFromJsonBytes(bytes, &c)

		// 优雅停止旧服务
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

	// 启动初始服务
	wg.Add(1)
	go func(c config.Config) {
		defer wg.Done()

		Run(c)
	}(c)

	// 等待所有服务退出
	wg.Wait()

}

// Run 运行 gRPC 服务
//
// 参数:
//   - c: 服务配置
//
// 功能:
//  1. 创建服务上下文
//  2. 创建 gRPC 服务器
//  3. 注册服务接口
//  4. 添加拦截器
//  5. 启动服务监听
func Run(c config.Config) {
	// 创建服务上下文
	ctx := svc.NewServiceContext(c)

	// 创建 gRPC 服务器
	s := zrpc.MustNewServer(c.RpcServerConf, func(srv *grpc.Server) {

		grpcServer = srv

		// 注册社交服务
		social.RegisterSocialServer(grpcServer, server.NewSocialServer(ctx))

		// 在开发/测试模式下启用 gRPC 反射（方便使用 grpcurl 等工具调试）
		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	// 添加幂等性拦截器（防止重复请求）
	//s.AddUnaryInterceptors(rpcserver.LogInterceptor, rpcserver.SyncxLimitInterceptor(10))
	s.AddUnaryInterceptors(interceptor.NewIdempotenceServer(interceptor.NewDefaultIdempotent(c.Cache[0].RedisConf)))

	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
