// =============================================================================
// 好友申请 Handler - 发送好友申请接口
// =============================================================================
// 提供发送好友申请的 HTTP 接口
//
// 接口信息:
//   - 路径: POST /friend/putin
//   - 认证: 需要 JWT token
//   - 请求参数:
//     {
//       "userId": "目标用户ID",
//       "reqMsg": "申请消息",
//       "reqTime": "申请时间戳"
//     }
//
// 响应格式:
//   {}  // 成功返回空对象
//
// 业务流程:
//   1. 解析请求参数（目标用户ID、申请消息、申请时间）
//   2. 从上下文获取当前用户 ID
//   3. 调用 Social RPC 创建好友申请记录
//   4. 返回处理结果
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

// FriendPutInHandler 好友申请处理器
// 创建并返回处理好友申请请求的 HTTP 处理函数
//
// 参数:
//   - svcCtx: 服务上下文，包含 RPC 客户端等依赖
//
// 返回:
//   - http.HandlerFunc: HTTP 处理函数
func FriendPutInHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 解析请求参数
		var req types.FriendPutInReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		// 创建业务逻辑处理器并执行
		l := friend.NewFriendPutInLogic(r.Context(), svcCtx)
		resp, err := l.FriendPutIn(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
