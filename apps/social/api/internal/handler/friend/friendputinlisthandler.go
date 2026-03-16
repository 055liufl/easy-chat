// =============================================================================
// 好友申请列表 Handler - 获取好友申请列表接口
// =============================================================================
// 提供获取当前用户收到的好友申请列表的 HTTP 接口
//
// 接口信息:
//   - 路径: GET /friend/putin/list
//   - 认证: 需要 JWT token
//   - 请求参数: 无（从 JWT token 中获取用户 ID）
//
// 响应格式:
//   {
//     "list": [
//       {
//         "id": "申请记录ID",
//         "userId": "申请人用户ID",
//         "reqMsg": "申请消息",
//         "reqTime": "申请时间",
//         "handleResult": "处理状态（0-待处理，1-已同意，2-已拒绝）",
//         "nickname": "申请人昵称",
//         "avatar": "申请人头像"
//       }
//     ]
//   }
//
// 业务流程:
//   1. 解析请求参数
//   2. 从上下文获取当前用户 ID
//   3. 调用 Social RPC 获取好友申请列表
//   4. 调用 User RPC 批量获取申请人详细信息
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

// FriendPutInListHandler 好友申请列表处理器
// 创建并返回处理好友申请列表请求的 HTTP 处理函数
//
// 参数:
//   - svcCtx: 服务上下文，包含 RPC 客户端等依赖
//
// 返回:
//   - http.HandlerFunc: HTTP 处理函数
func FriendPutInListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 解析请求参数
		var req types.FriendPutInListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		// 创建业务逻辑处理器并执行
		l := friend.NewFriendPutInListLogic(r.Context(), svcCtx)
		resp, err := l.FriendPutInList(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
