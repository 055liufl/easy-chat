// =============================================================================
// 入群申请处理 Handler - 处理入群申请接口（同意/拒绝）
// =============================================================================
// 提供处理入群申请的 HTTP 接口，用于同意或拒绝入群申请
//
// 接口信息:
//   - 路径: PUT /group/putin/handle
//   - 认证: 需要 JWT token
//   - 请求参数:
//     {
//       "groupReqId": "入群申请记录ID",
//       "groupId": "群组ID",
//       "handleResult": "处理结果（1-同意，2-拒绝）"
//     }
//
// 响应格式:
//   {}  // 成功返回空对象
//
// 业务流程:
//   1. 解析请求参数（申请记录ID、群组ID、处理结果）
//   2. 从上下文获取当前用户 ID（处理人，通常是群主或管理员）
//   3. 调用 Social RPC 处理入群申请
//   4. 如果同意，将申请人加入群组
//   5. 返回处理结果
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

// GroupPutInHandleHandler 入群申请处理处理器
// 创建并返回处理入群申请处理请求的 HTTP 处理函数
//
// 参数:
//   - svcCtx: 服务上下文，包含 RPC 客户端等依赖
//
// 返回:
//   - http.HandlerFunc: HTTP 处理函数
func GroupPutInHandleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 解析请求参数
		var req types.GroupPutInHandleRep
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		// 创建业务逻辑处理器并执行
		l := group.NewGroupPutInHandleLogic(r.Context(), svcCtx)
		resp, err := l.GroupPutInHandle(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
