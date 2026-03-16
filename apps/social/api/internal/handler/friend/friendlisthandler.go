// =============================================================================
// 好友列表 Handler - 获取用户好友列表接口
// =============================================================================
// 提供获取当前用户好友列表的 HTTP 接口
//
// 接口信息:
//   - 路径: GET /friend/list
//   - 认证: 需要 JWT token
//   - 请求参数: 无（从 JWT token 中获取用户 ID）
//
// 响应格式:
//   {
//     "list": [
//       {
//         "id": "好友关系ID",
//         "friendUid": "好友用户ID",
//         "nickname": "好友昵称",
//         "avatar": "好友头像"
//       }
//     ]
//   }
//
// 业务流程:
//   1. 解析请求参数
//   2. 从上下文获取当前用户 ID
//   3. 调用 Social RPC 获取好友关系列表
//   4. 调用 User RPC 批量获取好友详细信息
//   5. 组装响应数据返回
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

// FriendListHandler 好友列表处理器
// 创建并返回处理好友列表请求的 HTTP 处理函数
//
// 参数:
//   - svcCtx: 服务上下文，包含 RPC 客户端等依赖
//
// 返回:
//   - http.HandlerFunc: HTTP 处理函数
func FriendListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 解析请求参数
		var req types.FriendListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		// 创建业务逻辑处理器并执行
		l := friend.NewFriendListLogic(r.Context(), svcCtx)
		resp, err := l.FriendList(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
