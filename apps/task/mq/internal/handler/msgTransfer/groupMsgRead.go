// =============================================================================
// 群聊已读消息合并管理器 - 延迟合并推送优化
// =============================================================================
// 本模块负责群聊已读消息的延迟合并推送，避免频繁推送导致的性能问题
//
// 业务场景:
//   - 群聊中多个用户同时标记消息已读
//   - 如果每次已读都立即推送，会导致大量推送请求
//   - 通过延迟合并，将多个用户的已读记录合并后一次性推送
//
// 合并策略:
//   - 时间阈值：超过 GroupMsgReadRecordDelayTime 后推送
//   - 数量阈值：累计 GroupMsgReadRecordDelayCount 条后推送
//   - 空闲清理：超过 2 倍延迟时间无新消息则标记为空闲
//
// 工作流程:
//   1. 创建管理器时启动 transfer 协程
//   2. 接收新的已读消息，合并到 push.ReadRecords
//   3. 定时检查是否满足推送条件（时间或数量）
//   4. 满足条件后推送到 pushCh，由 MsgReadTransfer 统一推送
//   5. 空闲时清理资源
//
// =============================================================================

package msgTransfer

import (
	"github.com/zeromicro/go-zero/core/logx"
	"imooc.com/easy-chat/apps/im/ws/ws"
	"imooc.com/easy-chat/pkg/constants"
	"sync"
	"time"
)

// groupMsgRead 群聊已读消息合并管理器
// 负责单个会话的已读消息延迟合并推送
type groupMsgRead struct {
	mu             sync.Mutex   // 保护并发访问
	conversationId string       // 会话 ID
	push           *ws.Push     // 待推送的数据（累积多个用户的已读记录）
	pushCh         chan *ws.Push // 推送通道，推送到 MsgReadTransfer
	count          int          // 累计已读消息数量
	pushTime       time.Time    // 上次推送时间
	done           chan struct{} // 停止信号
}

// newGroupMsgRead 创建群聊已读消息合并管理器
// 创建后立即启动 transfer 协程进行延迟合并处理
//
// 参数:
//   - push: 初始推送数据（第一条已读消息）
//   - pushCh: 推送通道，用于将合并后的数据推送到 MsgReadTransfer
//
// 返回:
//   - *groupMsgRead: 合并管理器实例
func newGroupMsgRead(push *ws.Push, pushCh chan *ws.Push) *groupMsgRead {
	m := &groupMsgRead{
		conversationId: push.ConversationId,
		push:           push,
		pushCh:         pushCh,
		count:          1,
		pushTime:       time.Now(),
		done:           make(chan struct{}),
	}

	// 启动延迟合并协程
	go m.transfer()
	return m
}

// mergePush 合并已读消息
// 将新的已读记录合并到待推送数据中
//
// 参数:
//   - push: 新的已读消息数据
func (m *groupMsgRead) mergePush(push *ws.Push) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.count++
	// 合并已读记录（map[msgId]base64编码的Bitmap）
	for msgId, read := range push.ReadRecords {
		m.push.ReadRecords[msgId] = read
	}
}

// transfer 延迟合并推送协程
// 定时检查是否满足推送条件，满足则推送
//
// 推送条件:
//   1. 超过时间阈值（GroupMsgReadRecordDelayTime）
//   2. 超过数量阈值（GroupMsgReadRecordDelayCount）
//
// 工作流程:
//   1. 使用定时器检查时间阈值
//   2. 使用 default 分支检查数量阈值
//   3. 满足条件后推送到 pushCh
//   4. 检查是否空闲，空闲则发送空消息通知 MsgReadTransfer 清理
func (m *groupMsgRead) transfer() {
	// 定时器：检查时间阈值
	timer := time.NewTimer(GroupMsgReadRecordDelayTime / 2)
	defer timer.Stop()

	for {
		select {
		case <-m.done:
			// 停止信号
			return
		case <-timer.C:
			// 定时器触发：检查是否超过时间阈值
			m.mu.Lock()

			pushTime := m.pushTime
			val := GroupMsgReadRecordDelayTime - time.Since(pushTime)
			push := m.push
			logx.Infof("timer.C %v val %v", time.Now(), val)
			// 未达到推送条件：时间未到且数量未达标
			if val > 0 && m.count < GroupMsgReadRecordDelayCount || push == nil {
				if val > 0 {
					timer.Reset(val)
				}

				// 未达标
				m.mu.Unlock()
				continue
			}

			// 达到推送条件：重置状态并推送
			m.pushTime = time.Now()
			m.push = nil
			m.count = 0
			timer.Reset(GroupMsgReadRecordDelayTime / 2)
			m.mu.Unlock()

			// 推送合并后的已读记录
			logx.Infof("超过 合并的条件推送 %v ", push)
			m.pushCh <- push
		default:
			// 检查数量阈值
			m.mu.Lock()

			if m.count >= GroupMsgReadRecordDelayCount {
				// 达到数量阈值：立即推送
				push := m.push
				m.push = nil
				m.count = 0
				m.mu.Unlock()

				// 推送
				logx.Infof("default 超过 合并的条件推送 %v ", push)
				m.pushCh <- push
				continue
			}

			if m.isIdle() {
				// 空闲状态：发送空消息通知 MsgReadTransfer 清理
				m.mu.Unlock()
				// 使得 msgReadTransfer 释放
				m.pushCh <- &ws.Push{
					ChatType:       constants.GroupChatType,
					ConversationId: m.conversationId,
				}
				continue
			}
			m.mu.Unlock()

			// 短暂休眠，避免 CPU 空转
			tempDelay := GroupMsgReadRecordDelayTime / 4
			if tempDelay > time.Second {
				tempDelay = time.Second
			}
			time.Sleep(tempDelay)
		}
	}
}

// IsIdle 检查是否空闲（线程安全）
//
// 返回:
//   - bool: true 表示空闲，false 表示活跃
func (m *groupMsgRead) IsIdle() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.isIdle()
}

// isIdle 检查是否空闲（内部方法，需要持有锁）
// 空闲条件：超过 2 倍延迟时间无新消息，且无待推送数据
//
// 返回:
//   - bool: true 表示空闲，false 表示活跃
func (m *groupMsgRead) isIdle() bool {
	pushTime := m.pushTime
	val := GroupMsgReadRecordDelayTime*2 - time.Since(pushTime)

	if val <= 0 && m.push == nil && m.count == 0 {
		return true
	}

	return false
}

// clear 清理资源
// 停止 transfer 协程并清空待推送数据
func (m *groupMsgRead) clear() {
	select {
	case <-m.done:
	default:
		close(m.done)
	}

	m.push = nil
}
