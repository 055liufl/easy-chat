// =============================================================================
// WebSocket Message - WebSocket 消息定义
// =============================================================================
// 定义 WebSocket 通信中使用的消息结构和消息类型
//
// 消息帧类型:
//   - FrameData: 业务数据帧，携带实际的业务数据
//   - FramePing: 心跳帧，用于保持连接活跃
//   - FrameAck: 确认帧，用于消息确认机制
//   - FrameNoAck: 无需确认帧，不需要 ACK 的消息
//   - FrameErr: 错误帧，携带错误信息
//
// 消息确认机制:
//   - 通过 Id 和 AckSeq 实现消息的可靠传输
//   - 支持消息重传和超时处理
//
// =============================================================================
package websocket

import "time"

// FrameType 消息帧类型
// 定义了 WebSocket 通信中的各种消息类型
type FrameType uint8

const (
	// FrameData 数据帧
	// 携带业务数据的消息，需要根据 Method 字段路由到对应的处理器
	FrameData FrameType = 0x0

	// FramePing 心跳帧
	// 用于保持连接活跃，防止连接超时断开
	FramePing FrameType = 0x1

	// FrameAck 确认帧
	// 用于消息确认机制，确保消息可靠传输
	FrameAck FrameType = 0x2

	// FrameNoAck 无需确认帧
	// 不需要 ACK 确认的消息，用于对可靠性要求不高的场景
	FrameNoAck FrameType = 0x3

	// FrameErr 错误帧
	// 携带错误信息，用于通知对方发生了错误
	FrameErr FrameType = 0x9

	//FrameHeaders      FrameType = 0x1
	//FramePriority     FrameType = 0x2
	//FrameRSTStream    FrameType = 0x3
	//FrameSettings     FrameType = 0x4
	//FramePushPromise  FrameType = 0x5
	//FrameGoAway       FrameType = 0x7
	//FrameWindowUpdate FrameType = 0x8
	//FrameContinuation FrameType = 0x9
)

// Message WebSocket 消息结构
// 定义了消息的完整格式，包括消息类型、确认信息、路由信息和数据
type Message struct {
	FrameType `json:"frameType"` // 消息帧类型
	Id        string      `json:"id"`       // 消息唯一标识，用于消息确认和去重
	AckSeq    int         `json:"ackSeq"`   // 确认序号，用于消息确认机制
	ackTime   time.Time   `json:"ackTime"`  // 确认时间，用于超时判断（不序列化）
	errCount  int         `json:"errCount"` // 错误计数，用于重试控制（不序列化）
	Method    string      `json:"method"`   // 方法名，用于路由到对应的处理器
	FormId    string      `json:"formId"`   // 发送方 ID，标识消息来源
	Data      interface{} `json:"data"`     // 消息数据，可以是任意类型
}

// NewMessage 创建新的数据消息
// 用于构造业务数据消息
//
// 参数:
//   - formId: 发送方 ID
//   - data: 消息数据
//
// 返回:
//   - *Message: 初始化的消息对象
func NewMessage(formId string, data interface{}) *Message {
	return &Message{
		FrameType: FrameData,
		FormId:    formId,
		Data:      data,
	}
}

// NewErrMessage 创建错误消息
// 用于构造错误通知消息
//
// 参数:
//   - err: 错误对象
//
// 返回:
//   - *Message: 包含错误信息的消息对象
func NewErrMessage(err error) *Message {
	return &Message{
		FrameType: FrameErr,
		Data:      err.Error(),
	}
}
