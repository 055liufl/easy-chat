// =============================================================================
// Social API 服务入口 - 社交服务 API 主程序
// =============================================================================
// Social API 服务提供社交关系管理的 HTTP 接口，包括:
//   - 好友管理: 好友申请、好友列表、好友在线状态查询
//   - 群组管理: 创建群组、入群申请、群成员管理、群成员在线状态查询
//
// 服务特性:
//   - JWT 认证: 所有接口都需要 JWT token 认证
//   - 幂等性控制: 防止重复请求导致的数据不一致
//   - 限流保护: 防止恶意请求或异常流量冲击
//   - 配置热更新: 支持从 ETCD 动态加载配置
//
// 依赖服务:
//   - Social RPC: 社交关系服务
//   - User RPC: 用户信息服务
//   - IM RPC: 即时通讯服务
//   - Redis: 缓存服务
//   - ETCD: 配置中心
//
// =============================================================================
package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/proc"
	"github.com/zeromicro/go-zero/rest/httpx"
	"imooc.com/easy-chat/pkg/configserver"
	"imooc.com/easy-chat/pkg/resultx"
	"sync"

	"imooc.com/easy-chat/apps/social/api/internal/config"
	"imooc.com/easy-chat/apps/social/api/internal/handler"
	"imooc.com/easy-chat/apps/social/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/dev/social.yaml", "the config file")
var wg sync.WaitGroup

// main 主函数
// 启动 Social API 服务，支持配置热更新
func main() {
	flag.Parse()

	var c config.Config
	// 从配置中心加载配置，支持热更新
	err := configserver.NewConfigServer(*configFile, configserver.NewSail(&configserver.Config{
		ETCDEndpoints:  "192.168.117.24:3379",        // ETCD 地址
		ProjectKey:     "98c6f2c2287f4c73cea3d40ae7ec3ff2", // 项目密钥
		Namespace:      "social",                     // 命名空间
		Configs:        "social-api.yaml",            // 配置文件名
		ConfigFilePath: "./etc/conf",                 // 配置文件路径
		LogLevel:       "DEBUG",                      // 日志级别
	})).MustLoad(&c, func(bytes []byte) error {
		// 配置更新回调函数
		var c config.Config
		configserver.LoadFromJsonBytes(bytes, &c)

		// 优雅关闭当前服务
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

	wg.Wait()

}

// Run 运行 Social API 服务
// 初始化服务上下文、注册路由、启动 HTTP 服务器
//
// 参数:
//   - c: 配置信息
func Run(c config.Config) {
	// 创建 REST 服务器
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// 初始化服务上下文
	ctx := svc.NewServiceContext(c)
	// 注册路由
	handler.RegisterHandlers(server, ctx)

	// 设置错误处理器
	httpx.SetErrorHandlerCtx(resultx.ErrHandler(c.Name))
	// 设置成功响应处理器
	httpx.SetOkHandler(resultx.OkHandler)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
