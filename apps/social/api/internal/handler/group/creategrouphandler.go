// =============================================================================
// 创建群组 Handler - 创建群组接口
// =============================================================================
// 提供创建新群组的 HTTP 接口
//
// 接口信息:
//   - 路径: POST /group/create
//   - 认证: 需要 JWT token
//   - 请求参数:
//     {
//       "name": "群组名称",
//       "icon": "群组图标URL",
//       "desc": "群组描述"
//     }
//
// 响应格式:
//   {
//     "id": "群组ID"
//   }
//
// 业务流程:
//   1. 解析请求参数（群组名称、图标、描述）
//   2. 从上下文获取当前用户 ID（群主）
//   3. 调用 Social RPC 创建群组
//   4. 返回新创建的群组 ID
//
// =============================================================================
package group

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"imooc.com/easy-chat/apps/social/api/internal/logic/group"
	"imooc.com/easy-chat/apps/social/api/internal/svc"
	"imooc.com/easy-chat/apps/social/api/internal/types"
)

// CreateGroupHandler 创建群组处理器
// 创建并返回处理创建群组请求的 HTTP 处理函数
//
// 参数:
//   - svcCtx: 服务上下文，包含 RPC 客户端等依赖
//
// 返回:
//   - http.HandlerFunc: HTTP 处理函数
func CreateGroupHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 解析请求参数
		var req types.GroupCreateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		// 创建业务逻辑处理器并执行
		l := group.NewCreateGroupLogic(r.Context(), svcCtx)
		resp, err := l.CreateGroup(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
