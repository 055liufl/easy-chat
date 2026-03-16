// =============================================================================
// 群成员在线状态 Handler - 查询群成员在线状态接口
// =============================================================================
// 提供批量查询群成员在线状态的 HTTP 接口
//
// 接口信息:
//   - 路径: POST /group/users/online
//   - 认证: 需要 JWT token
//   - 请求参数:
//     {
//       "groupId": "群组ID"
//     }
//
// 响应格式:
//   {
//     "onlineList": {
//       "用户ID1": true,   // true-在线，false-离线
//       "用户ID2": false
//     }
//   }
//
// 业务流程:
//   1. 解析请求参数（群组ID）
//   2. 调用 Social RPC 获取群成员列表
//   3. 从 Redis 缓存中查询成员的在线状态
//   4. 返回在线状态映射表
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

// GroupUserOnlineHandler 群成员在线状态处理器
// 创建并返回处理群成员在线状态查询请求的 HTTP 处理函数
//
// 参数:
//   - svcCtx: 服务上下文，包含 RPC 客户端等依赖
//
// 返回:
//   - http.HandlerFunc: HTTP 处理函数
func GroupUserOnlineHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 解析请求参数
		var req types.GroupUserOnlineReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		// 创建业务逻辑处理器并执行
		l := group.NewGroupUserOnlineLogic(r.Context(), svcCtx)
		resp, err := l.GroupUserOnline(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
