// =============================================================================
// WebSocket Default Values - WebSocket 默认配置值
// =============================================================================
// 定义 WebSocket 服务器的默认配置常量
//
// 默认值说明:
//   - defaultMaxConnectionIdle: 默认最大空闲时长为无限大（永不超时）
//   - defaultAckTimeout: 默认 ACK 超时时间为 30 秒
//   - defaultConcurrency: 默认并发处理任务数为 10
//
// =============================================================================
package websocket

import (
	"math"
	"time"
)

const (
	// defaultMaxConnectionIdle 默认最大空闲连接时长
	// 设置为 math.MaxInt64，表示连接永不因空闲而超时
	// 生产环境建议设置合理的超时时间（如 5 分钟）
	defaultMaxConnectionIdle = time.Duration(math.MaxInt64)

	// defaultAckTimeout 默认 ACK 超时时间
	// 在 RigorAck 模式下，如果客户端在此时间内未确认消息，则放弃该消息
	defaultAckTimeout = 30 * time.Second

	// defaultConcurrency 默认并发处理任务数
	// 控制同时处理的消息数量，避免过多并发导致资源耗尽
	defaultConcurrency = 10
)
