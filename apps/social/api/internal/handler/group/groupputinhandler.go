// =============================================================================
// 入群申请 Handler - 发送入群申请接口
// =============================================================================
// 提供发送入群申请的 HTTP 接口
//
// 接口信息:
//   - 路径: POST /group/putin
//   - 认证: 需要 JWT token
//   - 请求参数:
//     {
//       "groupId": "群组ID",
//       "reqMsg": "申请消息",
//       "reqTime": "申请时间戳"
//     }
//
// 响应格式:
//   {}  // 成功返回空对象
//
// 业务流程:
//   1. 解析请求参数（群组ID、申请消息、申请时间）
//   2. 从上下文获取当前用户 ID
//   3. 调用 Social RPC 创建入群申请记录
//   4. 返回处理结果
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

// GroupPutInHandler 入群申请处理器
// 创建并返回处理入群申请请求的 HTTP 处理函数
//
// 参数:
//   - svcCtx: 服务上下文，包含 RPC 客户端等依赖
//
// 返回:
//   - http.HandlerFunc: HTTP 处理函数
func GroupPutInHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 解析请求参数
		var req types.GroupPutInRep
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		// 创建业务逻辑处理器并执行
		l := group.NewGroupPutInLogic(r.Context(), svcCtx)
		resp, err := l.GroupPutIn(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
