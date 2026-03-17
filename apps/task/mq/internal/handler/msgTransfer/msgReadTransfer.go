// =============================================================================
// 已读消息传输处理器 - Kafka 已读消息消费者
// =============================================================================
// 本模块负责消费 Kafka 中的已读消息，完成以下任务：
//   - 更新 MongoDB 中的消息已读状态
//   - 推送已读回执给发送方
//   - 群聊已读消息支持延迟合并推送（减少推送频率）
//
// 数据来源:
//   - Kafka Topic: MsgReadTransfer（已读消息队列）
//
// 业务场景:
//   - 用户标记消息已读后，IM 服务将已读消息发送到 Kafka
//   - Task 服务消费消息，更新 MongoDB 已读状态
//   - 单聊已读：立即推送给发送方
//   - 群聊已读：支持延迟合并推送（避免频繁推送）
//
// 群聊已读合并策略:
//   - 立即推送模式（GroupMsgReadHandlerAtTransfer）：每次已读立即推送
//   - 延迟合并模式（GroupMsgReadHandlerDelayTransfer）：
//     * 时间阈值：超过 GroupMsgReadRecordDelayTime 后推送
//     * 数量阈值：累计 GroupMsgReadRecordDelayCount 条后推送
//     * 空闲清理：超过 2 倍延迟时间无新消息则清理
//
// 处理流程:
//   1. 从 Kafka 消费已读消息
//   2. 更新 MongoDB 中的消息已读记录（Bitmap）
//   3. 单聊：立即推送已读回执
//   4. 群聊：根据配置选择立即推送或延迟合并推送
//
// =============================================================================

package msgTransfer

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"imooc.com/easy-chat/apps/im/ws/ws"
	"imooc.com/easy-chat/apps/task/mq/internal/svc"
	"imooc.com/easy-chat/apps/task/mq/mq"
	"imooc.com/easy-chat/pkg/bitmap"
	"imooc.com/easy-chat/pkg/constants"
	"sync"
	"time"
)

var (
	// GroupMsgReadRecordDelayTime 群聊已读消息延迟推送时间阈值（默认 1 秒）
	GroupMsgReadRecordDelayTime = time.Second
	// GroupMsgReadRecordDelayCount 群聊已读消息延迟推送数量阈值（默认 10 条）
	GroupMsgReadRecordDelayCount = 10
)

const (
	// GroupMsgReadHandlerAtTransfer 立即推送模式：每次已读立即推送
	GroupMsgReadHandlerAtTransfer = iota
	// GroupMsgReadHandlerDelayTransfer 延迟合并模式：延迟合并后推送
	GroupMsgReadHandlerDelayTransfer
)

// MsgReadTransfer 已读消息传输处理器
// 实现 Kafka 消费者接口，处理消息已读状态更新和推送
type MsgReadTransfer struct {
	*baseMsgTransfer // 继承基础消息传输功能（推送逻辑）

	cache.Cache // 缓存接口（预留）

	mu sync.Mutex // 保护 groupMsgs 的并发访问

	groupMsgs map[string]*groupMsgRead // 群聊已读消息合并管理器（key: conversationId）
	push      chan *ws.Push            // 推送通道，用于异步推送消息
}

// NewMsgReadTransfer 创建已读消息传输处理器实例
// 初始化配置并启动异步推送协程
//
// 参数:
//   - svc: 服务上下文，包含配置、MongoDB 模型等
//
// 返回:
//   - kq.ConsumeHandler: Kafka 消费者处理器接口
func NewMsgReadTransfer(svc *svc.ServiceContext) kq.ConsumeHandler {
	m := &MsgReadTransfer{
		baseMsgTransfer: NewBaseMsgTransfer(svc),
		groupMsgs:       make(map[string]*groupMsgRead, 1),
		push:            make(chan *ws.Push, 1),
	}

	// 如果配置为延迟合并模式，加载配置参数
	if svc.Config.MsgReadHandler.GroupMsgReadHandler != GroupMsgReadHandlerAtTransfer {
		if svc.Config.MsgReadHandler.GroupMsgReadRecordDelayCount > 0 {
			GroupMsgReadRecordDelayCount = svc.Config.MsgReadHandler.GroupMsgReadRecordDelayCount
		}

		if svc.Config.MsgReadHandler.GroupMsgReadRecordDelayTime > 0 {
			GroupMsgReadRecordDelayTime = time.Duration(svc.Config.MsgReadHandler.GroupMsgReadRecordDelayTime) * time.Second
		}
	}

	// 启动异步推送协程
	go m.transfer()

	return m
}

// Consume Kafka 消费者接口实现
// 消费已读消息，更新 MongoDB 已读状态并推送
//
// 业务逻辑:
//   1. 解析 Kafka 消息（JSON 格式）
//   2. 调用 UpdateChatLogRead 更新 MongoDB 已读状态
//   3. 单聊：直接推送到 push 通道
//   4. 群聊：根据配置选择立即推送或延迟合并推送
//
// 参数:
//   - key: Kafka 消息 key
//   - value: Kafka 消息 value（JSON 格式的已读消息）
//
// 返回:
//   - error: 处理失败时返回错误，Kafka 会重试
func (m *MsgReadTransfer) Consume(key, value string) error {
	m.Info("MsgReadTransfer ", value)

	var (
		data mq.MsgMarkRead
		ctx  = context.Background()
	)
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return err
	}

	// 更新 MongoDB 中的消息已读状态
	readRecords, err := m.UpdateChatLogRead(ctx, &data)
	if err != nil {
		return err
	}
	// 构建推送数据
	push := &ws.Push{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		ContentType:    constants.ContentMakeRead, // 内容类型：已读回执
		ReadRecords:    readRecords,               // 已读记录（map[msgId]base64编码的Bitmap）
	}

	switch data.ChatType {
	case constants.SingleChatType:
		// 单聊：直接推送
		m.push <- push
	case constants.GroupChatType:
		// 群聊：判断是否开启延迟合并
		if m.svcCtx.Config.MsgReadHandler.GroupMsgReadHandler == GroupMsgReadHandlerAtTransfer {
			// 立即推送模式
			m.push <- push
		}

		m.mu.Lock()
		defer m.mu.Unlock()

		push.SendId = "" // 群聊已读推送不需要 SendId（合并多个用户的已读）

		if _, ok := m.groupMsgs[push.ConversationId]; ok {
			m.Infof("merge push %v", push.ConversationId)
			// 已存在该会话的合并管理器，合并已读记录
			m.groupMsgs[push.ConversationId].mergePush(push)
		} else {
			m.Infof("newGroupMsgRead push %v", push.ConversationId)
			// 创建新的合并管理器
			m.groupMsgs[push.ConversationId] = newGroupMsgRead(push, m.push)
		}
	}

	return nil
}

// UpdateChatLogRead 更新消息已读状态
// 根据消息 ID 列表更新 MongoDB 中的已读记录
//
// 业务逻辑:
//   1. 根据消息 ID 列表查询 ChatLog
//   2. 单聊：直接标记为已读（[]byte{1}）
//   3. 群聊：使用 Bitmap 记录已读用户 ID
//   4. 更新 MongoDB 中的 ReadRecords 字段
//
// 参数:
//   - ctx: 上下文
//   - data: 已读消息数据，包含消息 ID 列表和发送者 ID
//
// 返回:
//   - map[string]string: 已读记录（key: msgId, value: base64 编码的 Bitmap）
//   - error: 更新失败时返回错误
func (m *MsgReadTransfer) UpdateChatLogRead(ctx context.Context, data *mq.MsgMarkRead) (map[string]string, error) {

	res := make(map[string]string)

	// 根据消息 ID 列表查询 ChatLog
	chatLogs, err := m.svcCtx.ChatLogModel.ListByMsgIds(ctx, data.MsgIds)
	if err != nil {
		return nil, err
	}

	// 更新每条消息的已读记录
	for _, chatLog := range chatLogs {
		switch chatLog.ChatType {
		case constants.SingleChatType:
			// 单聊：直接标记为已读
			chatLog.ReadRecords = []byte{1}
		case constants.GroupChatType:
			// 群聊：使用 Bitmap 记录已读用户 ID
			readRecords := bitmap.Load(chatLog.ReadRecords)
			readRecords.Set(data.SendId) // 标记该用户已读
			chatLog.ReadRecords = readRecords.Export()
		}

		// 将已读记录编码为 base64（用于推送）
		res[chatLog.ID.Hex()] = base64.StdEncoding.EncodeToString(chatLog.ReadRecords)

		// 更新 MongoDB
		err = m.svcCtx.ChatLogModel.UpdateMakeRead(ctx, chatLog.ID, chatLog.ReadRecords)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

// transfer 异步推送协程
// 从 push 通道读取推送数据，调用 Transfer 推送给接收方
//
// 业务逻辑:
//   1. 从 push 通道读取推送数据
//   2. 如果有接收方（RecvId 或 RecvIds），调用 Transfer 推送
//   3. 单聊：直接返回
//   4. 群聊延迟合并模式：检查是否空闲，空闲则清理合并管理器
//
// 注意:
//   - 该协程在 NewMsgReadTransfer 中启动，随服务生命周期运行
func (m *MsgReadTransfer) transfer() {
	for push := range m.push {
		// 如果有接收方，推送消息
		if push.RecvId != "" || len(push.RecvIds) > 0 {
			if err := m.Transfer(context.Background(), push); err != nil {
				m.Errorf("m transfer err %v push %v", err, push)
			}
		}

		// 单聊直接返回
		if push.ChatType == constants.SingleChatType {
			continue
		}

		// 立即推送模式直接返回
		if m.svcCtx.Config.MsgReadHandler.GroupMsgReadHandler == GroupMsgReadHandlerAtTransfer {
			continue
		}
		// 延迟合并模式：清理空闲的合并管理器
		m.mu.Lock()
		//
		if _, ok := m.groupMsgs[push.ConversationId]; ok && m.groupMsgs[push.ConversationId].IsIdle() {
			m.groupMsgs[push.ConversationId].clear()
			delete(m.groupMsgs, push.ConversationId)
		}

		m.mu.Unlock()

	}
}
