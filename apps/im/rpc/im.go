// =============================================================================
// IM RPC 服务主程序
// =============================================================================
// 即时通讯 RPC 服务的启动入口，提供以下功能:
//   - 聊天记录查询
//   - 会话管理（建立、查询、更新）
//   - 群聊会话创建
//
// 服务特性:
//   - 支持配置热更新（通过配置中心）
//   - 支持优雅关闭和重启
//   - 集成日志拦截器
//   - 开发/测试模式下启用 gRPC 反射
//
// 配置来源:
//   - 本地配置文件（默认: etc/dev/im.yaml）
//   - 远程配置中心（ETCD）
//
// =============================================================================
package main

import (
	"flag"
	"fmt"
	"imooc.com/easy-chat/pkg/configserver"
	"imooc.com/easy-chat/pkg/interceptor/rpcserver"
	"sync"

	"imooc.com/easy-chat/apps/im/rpc/im"
	"imooc.com/easy-chat/apps/im/rpc/internal/config"
	"imooc.com/easy-chat/apps/im/rpc/internal/server"
	"imooc.com/easy-chat/apps/im/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/dev/im.yaml", "the config file")  // 配置文件路径
var grpcServer *grpc.Server                                               // gRPC 服务器实例
var wg sync.WaitGroup                                                     // 等待组，用于协调多个 goroutine

func main() {
	flag.Parse()

	var c config.Config
	// 加载配置，支持配置热更新
	err := configserver.NewConfigServer(*configFile, configserver.NewSail(&configserver.Config{
		ETCDEndpoints:  "192.168.117.24:3379",
		ProjectKey:     "98c6f2c2287f4c73cea3d40ae7ec3ff2",
		Namespace:      "im",
		Configs:        "im-rpc.yaml",
		ConfigFilePath: "./etc/conf",
		LogLevel:       "DEBUG",
	})).MustLoad(&c, func(bytes []byte) error {
		// 配置更新回调：优雅关闭旧服务，启动新服务
		var c config.Config
		configserver.LoadFromJsonBytes(bytes, &c)

		grpcServer.GracefulStop()

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

	wg.Add(1)
	go func(c config.Config) {
		defer wg.Done()

		Run(c)
	}(c)

	wg.Wait()

}

// Run 启动 RPC 服务
// 初始化服务上下文、注册 gRPC 服务、添加拦截器并启动服务
//
// 参数:
//   - c: 服务配置
//
// 说明:
//   - 在开发/测试模式下启用 gRPC 反射，方便调试
//   - 添加日志拦截器记录所有 RPC 调用
//   - 服务启动后会阻塞，直到收到停止信号
func Run(c config.Config) {
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(srv *grpc.Server) {
		grpcServer = srv

		im.RegisterImServer(grpcServer, server.NewImServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	s.AddUnaryInterceptors(rpcserver.LogInterceptor)

	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
