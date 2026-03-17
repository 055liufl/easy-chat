// =============================================================================
// 用户注册 Handler
// =============================================================================
// 处理用户注册的 HTTP 请求，包括:
//   - 解析注册请求参数（手机号、密码、昵称、性别、头像）
//   - 调用业务逻辑层进行注册处理
//   - 返回注册结果（JWT token 和过期时间）
//
// 数据来源:
//   HTTP POST 请求，JSON 格式的请求体
//
// 业务场景:
//   用户通过手机号和密码进行注册，注册成功后返回访问令牌
//
// 请求路径:
//   POST /v1/user/register
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

// RegisterHandler 用户注册处理器
// 创建并返回用户注册的 HTTP 处理函数
//
// 参数:
//   - svcCtx: 服务上下文（包含配置和依赖）
//
// 返回:
//   - http.HandlerFunc: HTTP 处理函数
//
// 处理流程:
//   1. 解析 HTTP 请求参数到 RegisterReq 结构
//   2. 创建 RegisterLogic 实例
//   3. 调用业务逻辑层的 Register 方法
//   4. 返回注册结果或错误信息
func RegisterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterReq
		// 解析请求参数（支持 JSON、Form 等格式）
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		// 创建业务逻辑实例并执行注册
		l := user.NewRegisterLogic(r.Context(), svcCtx)
		resp, err := l.Register(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
