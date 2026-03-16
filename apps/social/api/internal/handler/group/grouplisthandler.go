// =============================================================================
// 群组列表 Handler - 获取用户群组列表接口
// =============================================================================
// 提供获取当前用户加入的群组列表的 HTTP 接口
//
// 接口信息:
//   - 路径: GET /group/list
//   - 认证: 需要 JWT token
//   - 请求参数: 无（从 JWT token 中获取用户 ID）
//
// 响应格式:
//   {
//     "list": [
//       {
//         "id": "群组ID",
//         "name": "群组名称",
//         "icon": "群组图标",
//         "desc": "群组描述",
//         "creatorUid": "创建者ID"
//       }
//     ]
//   }
//
// 业务流程:
//   1. 解析请求参数
//   2. 从上下文获取当前用户 ID
//   3. 调用 Social RPC 获取用户加入的群组列表
//   4. 返回群组详细信息
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

// GroupListHandler 群组列表处理器
// 创建并返回处理群组列表请求的 HTTP 处理函数
//
// 参数:
//   - svcCtx: 服务上下文，包含 RPC 客户端等依赖
//
// 返回:
//   - http.HandlerFunc: HTTP 处理函数
func GroupListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 解析请求参数
		var req types.GroupListRep
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		// 创建业务逻辑处理器并执行
		l := group.NewGroupListLogic(r.Context(), svcCtx)
		resp, err := l.GroupList(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
