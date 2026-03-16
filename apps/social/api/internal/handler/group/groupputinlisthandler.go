// =============================================================================
// 入群申请列表 Handler - 获取入群申请列表接口
// =============================================================================
// 提供获取群组收到的入群申请列表的 HTTP 接口
//
// 接口信息:
//   - 路径: GET /group/putin/list
//   - 认证: 需要 JWT token
//   - 请求参数:
//     {
//       "groupId": "群组ID"
//     }
//
// 响应格式:
//   {
//     "list": [
//       {
//         "id": "申请记录ID",
//         "groupId": "群组ID",
//         "reqId": "申请人用户ID",
//         "reqMsg": "申请消息",
//         "reqTime": "申请时间",
//         "handleResult": "处理状态（0-待处理，1-已同意，2-已拒绝）"
//       }
//     ]
//   }
//
// 业务流程:
//   1. 解析请求参数（群组ID）
//   2. 调用 Social RPC 获取入群申请列表
//   3. 返回申请记录列表
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

// GroupPutInListHandler 入群申请列表处理器
// 创建并返回处理入群申请列表请求的 HTTP 处理函数
//
// 参数:
//   - svcCtx: 服务上下文，包含 RPC 客户端等依赖
//
// 返回:
//   - http.HandlerFunc: HTTP 处理函数
func GroupPutInListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 解析请求参数
		var req types.GroupPutInListRep
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		// 创建业务逻辑处理器并执行
		l := group.NewGroupPutInListLogic(r.Context(), svcCtx)
		resp, err := l.GroupPutInList(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
