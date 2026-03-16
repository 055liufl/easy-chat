// =============================================================================
// 好友在线状态 Handler - 查询好友在线状态接口
// =============================================================================
// 提供批量查询好友在线状态的 HTTP 接口
//
// 接口信息:
//   - 路径: POST /friend/online
//   - 认证: 需要 JWT token
//   - 请求参数:
//     {
//       "friendIds": ["好友ID1", "好友ID2", ...]
//     }
//
// 响应格式:
//   {
//     "onlineList": {
//       "好友ID1": true,   // true-在线，false-离线
//       "好友ID2": false
//     }
//   }
//
// 业务流程:
//   1. 解析请求参数（好友ID列表）
//   2. 调用 IM RPC 批量查询好友在线状态
//   3. 返回在线状态映射表
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

// FriendsOnlineHandler 好友在线状态处理器
// 创建并返回处理好友在线状态查询请求的 HTTP 处理函数
//
// 参数:
//   - svcCtx: 服务上下文，包含 RPC 客户端等依赖
//
// 返回:
//   - http.HandlerFunc: HTTP 处理函数
func FriendsOnlineHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 解析请求参数
		var req types.FriendsOnlineReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		// 创建业务逻辑处理器并执行
		l := friend.NewFriendsOnlineLogic(r.Context(), svcCtx)
		resp, err := l.FriendsOnline(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
