// =============================================================================
// Kafka 消息生产者客户端 - 消息发送封装
// =============================================================================
// 本模块提供 Kafka 消息生产者客户端，用于其他服务（如 IM 服务）发送消息到 Kafka
//
// 客户端列表:
//   - MsgChatTransferClient: 聊天消息生产者，发送聊天消息到 Kafka
//   - MsgReadTransferClient: 已读消息生产者，发送已读消息到 Kafka
//
// 业务场景:
//   - IM 服务接收到用户发送的消息后，通过 MsgChatTransferClient 发送到 Kafka
//   - IM 服务接收到用户标记已读后，通过 MsgReadTransferClient 发送到 Kafka
//   - Task 服务从 Kafka 消费消息并处理
//
// =============================================================================

package mqclient

import (
	"encoding/json"
	"github.com/zeromicro/go-queue/kq"
	"imooc.com/easy-chat/apps/task/mq/mq"
)

// MsgChatTransferClient 聊天消息生产者客户端接口
type MsgChatTransferClient interface {
	Push(msg *mq.MsgChatTransfer) error
}

// msgChatTransferClient 聊天消息生产者客户端实现
type msgChatTransferClient struct {
	pusher *kq.Pusher // Kafka 生产者
}

// NewMsgChatTransferClient 创建聊天消息生产者客户端
//
// 参数:
//   - addr: Kafka 服务器地址列表
//   - topic: Kafka Topic 名称
//   - opts: 可选配置项
//
// 返回:
//   - MsgChatTransferClient: 聊天消息生产者客户端实例
func NewMsgChatTransferClient(addr []string, topic string, opts ...kq.PushOption) MsgChatTransferClient {
	return &msgChatTransferClient{
		pusher: kq.NewPusher(addr, topic),
	}
}

// Push 发送聊天消息到 Kafka
//
// 参数:
//   - msg: 聊天消息数据
//
// 返回:
//   - error: 发送失败时返回错误
func (c *msgChatTransferClient) Push(msg *mq.MsgChatTransfer) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return c.pusher.Push(string(body))
}

// MsgReadTransferClient 已读消息生产者客户端接口
type MsgReadTransferClient interface {
	Push(msg *mq.MsgMarkRead) error
}

// msgReadTransferClient 已读消息生产者客户端实现
type msgReadTransferClient struct {
	pusher *kq.Pusher // Kafka 生产者
}

// NewMsgReadTransferClient 创建已读消息生产者客户端
//
// 参数:
//   - addr: Kafka 服务器地址列表
//   - topic: Kafka Topic 名称
//   - opts: 可选配置项
//
// 返回:
//   - MsgReadTransferClient: 已读消息生产者客户端实例
func NewMsgReadTransferClient(addr []string, topic string, opts ...kq.PushOption) MsgReadTransferClient {
	return &msgReadTransferClient{
		pusher: kq.NewPusher(addr, topic),
	}
}

// Push 发送已读消息到 Kafka
//
// 参数:
//   - msg: 已读消息数据
//
// 返回:
//   - error: 发送失败时返回错误
func (c *msgReadTransferClient) Push(msg *mq.MsgMarkRead) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return c.pusher.Push(string(body))
}
