// =============================================================================
// 获取会话列表 Handler - HTTP 请求处理器
// =============================================================================
// API 接口: GET /v1/im/conversation
//
// 功能说明:
//   查询当前用户的所有会话列表
//   包括单聊和群聊会话
//
// 请求参数:
//   无（用户 ID 从 JWT Token 中获取）
//
// 响应格式:
//   {
//     "userId": "user_001",
//     "conversationList": {
//       "conv_001": {
//         "conversationId": "conv_001",
//         "chatType": 1,
//         "targetId": "user_002",
//         "isShow": true,
//         "seq": 1678886400,
//         "read": 10,
//         "total": 15,
//         "unread": 5
//       }
//     }
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

// getConversationsHandler 获取会话列表处理器
// 返回一个 HTTP 处理函数，用于处理获取会话列表的请求
//
// 参数:
//   - svcCtx: 服务上下文，包含所有依赖项
//
// 返回:
//   - http.HandlerFunc: HTTP 处理函数
//
// 处理流程:
//  1. 解析请求参数（无参数，用户 ID 从 JWT 中获取）
//  2. 调用业务逻辑层查询会话列表
//  3. 返回 JSON 响应或错误信息
func getConversationsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 解析请求参数（无参数）
		var req types.GetConversationsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		// 创建业务逻辑实例并执行
		l := logic.NewGetConversationsLogic(r.Context(), svcCtx)
		resp, err := l.GetConversations(&req)
		if err != nil {
			// 返回错误响应
			httpx.Error(w, err)
		} else {
			// 返回成功响应（JSON 格式）
			httpx.OkJson(w, resp)
		}
	}
}
