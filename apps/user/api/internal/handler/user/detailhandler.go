// =============================================================================
// 用户详情 Handler
// =============================================================================
// 处理获取用户详情的 HTTP 请求，包括:
//   - 从 JWT token 中解析用户 ID
//   - 调用业务逻辑层获取用户详细信息
//   - 返回用户信息
//
// 数据来源:
//   HTTP GET 请求，用户 ID 从 JWT token 中解析
//
// 业务场景:
//   已登录用户获取自己的详细信息（昵称、头像、性别等）
//
// 请求路径:
//   GET /v1/user/user（需要 JWT 认证）
//
// =============================================================================
package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"imooc.com/easy-chat/apps/user/api/internal/logic/user"
	"imooc.com/easy-chat/apps/user/api/internal/svc"
	"imooc.com/easy-chat/apps/user/api/internal/types"
)

// DetailHandler 用户详情处理器
// 创建并返回获取用户详情的 HTTP 处理函数
//
// 参数:
//   - svcCtx: 服务上下文（包含配置和依赖）
//
// 返回:
//   - http.HandlerFunc: HTTP 处理函数
//
// 处理流程:
//   1. 解析 HTTP 请求参数到 UserInfoReq 结构（实际无参数）
//   2. 创建 DetailLogic 实例
//   3. 调用业务逻辑层的 Detail 方法（从 context 中获取用户 ID）
//   4. 返回用户详细信息或错误信息
//
// 注意:
//   此接口需要 JWT 认证，用户 ID 从 JWT token 中解析
func DetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserInfoReq
		// 解析请求参数（实际无参数，用户 ID 从 JWT token 中获取）
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		// 创建业务逻辑实例并获取用户详情
		l := user.NewDetailLogic(r.Context(), svcCtx)
		resp, err := l.Detail(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
