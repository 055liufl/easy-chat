// =============================================================================
// 获取聊天记录 Handler - HTTP 请求处理器
// =============================================================================
// API 接口: GET /v1/im/chatlog
//
// 功能说明:
//   查询指定会话的聊天记录
//   支持按时间范围和数量过滤
//
// 请求参数:
//   - conversationId: 会话 ID（必填）
//   - startSendTime: 开始时间，Unix 时间戳（可选）
//   - endSendTime: 结束时间，Unix 时间戳（可选）
//   - count: 查询数量（可选，默认返回所有）
//
// 响应格式:
//   {
//     "list": [
//       {
//         "id": "msg_123",
//         "conversationId": "conv_456",
//         "sendId": "user_001",
//         "recvId": "user_002",
//         "msgType": 1,
//         "msgContent": "Hello",
//         "chatType": 1,
//         "sendTime": 1678886400
//       }
//     ]
//   }
//
// 认证要求:
//   需要 JWT Token 认证
//
// =============================================================================
package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"imooc.com/easy-chat/apps/im/api/internal/logic"
	"imooc.com/easy-chat/apps/im/api/internal/svc"
	"imooc.com/easy-chat/apps/im/api/internal/types"
)

// getChatLogHandler 获取聊天记录处理器
// 返回一个 HTTP 处理函数，用于处理获取聊天记录的请求
//
// 参数:
//   - svcCtx: 服务上下文，包含所有依赖项
//
// 返回:
//   - http.HandlerFunc: HTTP 处理函数
//
// 处理流程:
//  1. 解析请求参数（conversationId、时间范围、数量）
//  2. 调用业务逻辑层查询聊天记录
//  3. 返回 JSON 响应或错误信息
func getChatLogHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 解析请求参数
		var req types.ChatLogReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		// 创建业务逻辑实例并执行
		l := logic.NewGetChatLogLogic(r.Context(), svcCtx)
		resp, err := l.GetChatLog(&req)
		if err != nil {
			// 返回错误响应
			httpx.Error(w, err)
		} else {
			// 返回成功响应（JSON 格式）
			httpx.OkJson(w, resp)
		}
	}
}
