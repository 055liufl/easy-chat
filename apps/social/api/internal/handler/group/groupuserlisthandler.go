// =============================================================================
// 群成员列表 Handler - 获取群成员列表接口
// =============================================================================
// 提供获取指定群组成员列表的 HTTP 接口
//
// 接口信息:
//   - 路径: GET /group/users
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
//         "id": "成员记录ID",
//         "groupId": "群组ID",
//         "userId": "用户ID",
//         "roleLevel": "角色等级（1-群主，2-管理员，3-普通成员）",
//         "nickname": "用户昵称",
//         "avatar": "用户头像"
//       }
//     ]
//   }
//
// 业务流程:
//   1. 解析请求参数（群组ID）
//   2. 调用 Social RPC 获取群成员列表
//   3. 调用 User RPC 批量获取成员详细信息
//   4. 组装响应数据返回
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

// GroupUserListHandler 群成员列表处理器
// 创建并返回处理群成员列表请求的 HTTP 处理函数
//
// 参数:
//   - svcCtx: 服务上下文，包含 RPC 客户端等依赖
//
// 返回:
//   - http.HandlerFunc: HTTP 处理函数
func GroupUserListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 解析请求参数
		var req types.GroupUserListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		// 创建业务逻辑处理器并执行
		l := group.NewGroupUserListLogic(r.Context(), svcCtx)
		resp, err := l.GroupUserList(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
