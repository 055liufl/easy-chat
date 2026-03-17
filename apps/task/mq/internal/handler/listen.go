// =============================================================================
// Kafka 消费者监听器 - 消息队列服务注册
// =============================================================================
// 本模块负责注册和管理所有 Kafka 消费者服务，包括：
//   - 聊天消息消费者（MsgChatTransfer）
//   - 已读消息消费者（MsgReadTransfer）
//
// 业务场景:
//   - 启动时创建所有消费者实例
//   - 统一管理消费者生命周期
//   - 支持动态扩展新的消费者
//
// =============================================================================

package handler

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
	"imooc.com/easy-chat/apps/task/mq/internal/handler/msgTransfer"
	"imooc.com/easy-chat/apps/task/mq/internal/svc"
)

// Listen Kafka 消费者监听器
// 负责创建和管理所有消息队列消费者服务
type Listen struct {
	svc *svc.ServiceContext // 服务上下文，包含配置和依赖
}

// NewListen 创建监听器实例
//
// 参数:
//   - svc: 服务上下文，包含配置、数据库连接、RPC 客户端等
//
// 返回:
//   - *Listen: 监听器实例
func NewListen(svc *svc.ServiceContext) *Listen {
	return &Listen{svc: svc}
}

// Services 获取所有消费者服务列表
// 返回所有需要启动的 Kafka 消费者服务
//
// 消费者列表:
//   - MsgReadTransfer: 已读消息消费者，处理消息已读状态更新
//   - MsgChatTransfer: 聊天消息消费者，处理消息持久化和推送
//
// 返回:
//   - []service.Service: 消费者服务列表
func (l *Listen) Services() []service.Service {
	return []service.Service{
		// 已读消息消费者：消费 Kafka 中的已读消息，更新 MongoDB 已读状态并推送
		kq.MustNewQueue(l.svc.Config.MsgReadTransfer, msgTransfer.NewMsgReadTransfer(l.svc)),
		// 聊天消息消费者：消费 Kafka 中的聊天消息，持久化到 MongoDB 并推送给接收方
		kq.MustNewQueue(l.svc.Config.MsgChatTransfer, msgTransfer.NewMsgChatTransfer(l.svc)),
	}
}
