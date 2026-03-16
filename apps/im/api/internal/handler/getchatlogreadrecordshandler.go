// =============================================================================
// 获取消息已读记录 Handler - HTTP 请求处理器
// =============================================================================
// API 接口: GET /v1/im/chatlog/readRecords
//
// 功能说明:
//   查询指定消息的已读/未读用户列表
//   支持单聊和群聊消息的已读状态查询
//
// 请求参数:
//   - msgId: 消息 ID（必填，通过 Query 参数传递）
//
// 响应格式:
//   {
//     "reads": ["13800138000", "13900139000"],   // 已读用户手机号列表
//     "unReads": ["13700137000"]                 // 未读用户手机号列表
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

// getChatLogReadRecordsHandler 获取消息已读记录处理器
// 返回一个 HTTP 处理函数，用于处理获取消息已读记录的请求
//
// 参数:
//   - svcCtx: 服务上下文，包含所有依赖项
//
// 返回:
//   - http.HandlerFunc: HTTP 处理函数
//
// 处理流程:
//  1. 解析请求参数（msgId）
//  2. 调用业务逻辑层处理
//  3. 返回 JSON 响应或错误信息
func getChatLogReadRecordsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 解析请求参数
		var req types.GetChatLogReadRecordsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		// 创建业务逻辑实例并执行
		l := logic.NewGetChatLogReadRecordsLogic(r.Context(), svcCtx)
		resp, err := l.GetChatLogReadRecords(&req)
		if err != nil {
			// 返回错误响应
			httpx.Error(w, err)
		} else {
			// 返回成功响应（JSON 格式）
			httpx.OkJson(w, resp)
		}
	}
}
