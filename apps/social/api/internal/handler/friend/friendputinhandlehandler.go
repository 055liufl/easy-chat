// =============================================================================
// 好友申请处理 Handler - 处理好友申请接口（同意/拒绝）
// =============================================================================
// 提供处理好友申请的 HTTP 接口，用于同意或拒绝好友申请
//
// 接口信息:
//   - 路径: PUT /friend/putin/handle
//   - 认证: 需要 JWT token
//   - 请求参数:
//     {
//       "friendReqId": "好友申请记录ID",
//       "handleResult": "处理结果（1-同意，2-拒绝）"
//     }
//
// 响应格式:
//   {}  // 成功返回空对象
//
// 业务流程:
//   1. 解析请求参数（申请记录ID、处理结果）
//   2. 从上下文获取当前用户 ID
//   3. 调用 Social RPC 处理好友申请
//   4. 如果同意，建立双向好友关系
//   5. 返回处理结果
//
// =============================================================================
package friend

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"imooc.com/easy-chat/apps/social/api/internal/logic/friend"
	"imooc.com/easy-chat/apps/social/api/internal/svc"
	"imooc.com/easy-chat/apps/social/api/internal/types"
)

// FriendPutInHandleHandler 好友申请处理处理器
// 创建并返回处理好友申请处理请求的 HTTP 处理函数
//
// 参数:
//   - svcCtx: 服务上下文，包含 RPC 客户端等依赖
//
// 返回:
//   - http.HandlerFunc: HTTP 处理函数
func FriendPutInHandleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 解析请求参数
		var req types.FriendPutInHandleReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		// 创建业务逻辑处理器并执行
		l := friend.NewFriendPutInHandleLogic(r.Context(), svcCtx)
		resp, err := l.FriendPutInHandle(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
