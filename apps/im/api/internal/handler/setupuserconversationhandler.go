// =============================================================================
// 建立用户会话 Handler - HTTP 请求处理器
// =============================================================================
// API 接口: POST /v1/im/setup/conversation
//
// 功能说明:
//   创建或初始化一个新的会话
//   用于建立单聊或群聊会话
//
// 请求参数:
//   {
//     "sendId": "user_001",    // 发送者用户 ID
//     "recvId": "user_002",    // 接收者 ID（单聊为用户 ID，群聊为群组 ID）
//     "chatType": 1            // 聊天类型（1:单聊 2:群聊）
//   }
//
// 响应格式:
//   {} (空对象，表示成功)
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

// setUpUserConversationHandler 建立用户会话处理器
// 返回一个 HTTP 处理函数，用于处理建立用户会话的请求
//
// 参数:
//   - svcCtx: 服务上下文，包含所有依赖项
//
// 返回:
//   - http.HandlerFunc: HTTP 处理函数
//
// 处理流程:
//  1. 解析请求参数（sendId、recvId、chatType）
//  2. 调用业务逻辑层创建会话
//  3. 返回 JSON 响应或错误信息
func setUpUserConversationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 解析请求参数
		var req types.SetUpUserConversationReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		// 创建业务逻辑实例并执行
		l := logic.NewSetUpUserConversationLogic(r.Context(), svcCtx)
		resp, err := l.SetUpUserConversation(&req)
		if err != nil {
			// 返回错误响应
			httpx.Error(w, err)
		} else {
			// 返回成功响应（JSON 格式）
			httpx.OkJson(w, resp)
		}
	}
}
