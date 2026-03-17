// =============================================================================
// Task MQ 服务主程序 - Kafka 消息队列消费者服务
// =============================================================================
// 本服务是 Easy-Chat 系统的消息队列消费者，负责：
//   - 消费 Kafka 中的聊天消息并持久化到 MongoDB
//   - 消费 Kafka 中的已读消息并更新已读状态
//   - 通过 WebSocket 推送消息给在线用户
//
// 架构设计:
//   - 生产者：IM 服务将消息发送到 Kafka
//   - 消费者：Task 服务从 Kafka 消费消息
//   - 存储：MongoDB 持久化消息和会话
//   - 推送：WebSocket 推送消息给在线用户
//
// 服务启动流程:
//   1. 加载配置（本地文件 + ETCD 配置中心）
//   2. 初始化服务上下文（Redis、MongoDB、RPC 客户端等）
//   3. 创建 Kafka 消费者（MsgChatTransfer、MsgReadTransfer）
//   4. 启动消费者服务组
//   5. 支持配置热更新（ETCD 配置变更时自动重启）
//
// =============================================================================

package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/service"
	"imooc.com/easy-chat/apps/task/mq/internal/config"
	"imooc.com/easy-chat/apps/task/mq/internal/handler"
	"imooc.com/easy-chat/apps/task/mq/internal/svc"
	"imooc.com/easy-chat/pkg/configserver"
	"sync"
)

var configFile = flag.String("f", "etc/dev/task.yaml", "the config file")
var wg sync.WaitGroup

func main() {
	flag.Parse()

	var c config.Config
	// 加载配置（支持 ETCD 配置中心）
	err := configserver.NewConfigServer(*configFile, configserver.NewSail(&configserver.Config{
		ETCDEndpoints:  "192.168.117.24:3379",
		ProjectKey:     "98c6f2c2287f4c73cea3d40ae7ec3ff2",
		Namespace:      "task",
		Configs:        "task-mq.yaml",
		ConfigFilePath: "./etc/conf",
		LogLevel:       "DEBUG",
	})).MustLoad(&c, func(bytes []byte) error {
		// 配置热更新回调：配置变更时重新启动服务
		var c config.Config
		configserver.LoadFromJsonBytes(bytes, &c)

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

	wg.Wait()
}

// Run 启动 Task MQ 服务
// 初始化服务上下文并启动所有 Kafka 消费者
//
// 启动流程:
//   1. 设置日志和追踪配置
//   2. 创建服务上下文（初始化依赖）
//   3. 创建监听器（注册消费者）
//   4. 创建服务组并启动所有消费者
//
// 参数:
//   - c: 配置信息
func Run(c config.Config) {
	// 设置日志和追踪配置
	if err := c.SetUp(); err != nil {
		panic(err)
	}
	// 创建服务上下文
	ctx := svc.NewServiceContext(c)
	// 创建监听器（注册消费者）
	listen := handler.NewListen(ctx)

	// 创建服务组并添加所有消费者
	serviceGroup := service.NewServiceGroup()
	for _, s := range listen.Services() {
		serviceGroup.Add(s)
	}
	fmt.Println("Starting mqueue at ...")
	// 启动服务组（阻塞运行）
	serviceGroup.Start()
}
