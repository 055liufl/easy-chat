// =============================================================================
// IM WebSocket 服务 - 即时通讯 WebSocket 服务主程序
// =============================================================================
// 提供即时通讯的 WebSocket 服务，支持:
//   - 实时消息推送
//   - 用户在线状态管理
//   - 会话消息处理
//   - 消息已读未读标记
//
// 服务特性:
//   - 支持动态配置热更新（通过 ETCD 配置中心）
//   - JWT Token 认证保障连接安全
//   - 支持多实例并发运行
//   - 优雅关闭机制
//
// 配置来源:
//   - 本地配置文件（默认: etc/dev/im.yaml）
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
	"imooc.com/easy-chat/apps/im/ws/internal/config"
	"imooc.com/easy-chat/apps/im/ws/internal/handler"
	"imooc.com/easy-chat/apps/im/ws/internal/svc"
	"imooc.com/easy-chat/apps/im/ws/websocket"
	"imooc.com/easy-chat/pkg/configserver"
	"sync"
)

// configFile 配置文件路径，通过命令行参数 -f 指定
var configFile = flag.String("f", "etc/dev/im.yaml", "the config file")

// wg 等待组，用于等待所有服务实例优雅退出
var wg sync.WaitGroup

// main 主函数
// 负责初始化配置、启动 WebSocket 服务
//
// 执行流程:
//  1. 解析命令行参数
//  2. 加载配置（支持本地文件和 ETCD 配置中心）
//  3. 启动 WebSocket 服务实例
//  4. 监听配置变更，支持热更新（自动重启服务）
//  5. 等待所有服务实例退出
func main() {
	flag.Parse()

	var c config.Config

	// 初始化配置服务器，支持从 ETCD 加载配置并监听配置变更
	err := configserver.NewConfigServer(*configFile, configserver.NewSail(&configserver.Config{
		ETCDEndpoints:  "192.168.117.24:3379",           // ETCD 服务地址
		ProjectKey:     "98c6f2c2287f4c73cea3d40ae7ec3ff2", // 项目唯一标识
		Namespace:      "im",                             // 配置命名空间
		Configs:        "im-ws.yaml",                     // 配置文件名
		ConfigFilePath: "./etc/conf",                     // 本地配置缓存路径
		LogLevel:       "DEBUG",                          // 日志级别
	})).MustLoad(&c, func(bytes []byte) error {
		// 配置变更回调函数，当 ETCD 中的配置发生变化时触发
		var c config.Config
		configserver.LoadFromJsonBytes(bytes, &c)

		// 启动新的服务实例（配置热更新）
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

// Run 运行 WebSocket 服务实例
// 初始化服务上下文、创建 WebSocket 服务器、注册路由处理器
//
// 服务组件:
//  - ServiceContext: 服务上下文，包含配置、数据库连接、MQ 客户端等
//  - WebSocket Server: WebSocket 服务器，处理连接管理和消息路由
//  - JWT 认证: 连接建立时的身份验证
//  - 路由处理器: 处理各类业务消息（聊天、在线状态、已读标记等）
//
// 参数:
//   - c: 服务配置对象，包含监听地址、JWT 密钥、数据库配置等
func Run(c config.Config) {
	// 执行配置初始化（如日志、追踪等）
	if err := c.SetUp(); err != nil {
		panic(err)
	}

	// 创建服务上下文，初始化依赖组件
	ctx := svc.NewServiceContext(c)

	// 创建 WebSocket 服务器
	srv := websocket.NewServer(c.ListenOn,
		websocket.WithServerAuthentication(handler.NewJwtAuth(ctx)), // 配置 JWT 认证
		//websocket.WithServerAck(websocket.RigorAck),                // 可选: 启用严格消息确认机制
		//websocket.WithServerMaxConnectionIdle(10*time.Second),      // 可选: 设置连接空闲超时时间
	)
	defer srv.Stop() // 确保服务优雅关闭

	// 注册业务路由处理器
	handler.RegisterHandlers(srv, ctx)

	fmt.Println("start websocket server at ", c.ListenOn, " ..... ")
	srv.Start() // 启动服务（阻塞）
}
