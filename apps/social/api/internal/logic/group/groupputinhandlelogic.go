// =============================================================================
// 入群申请处理 Logic - 处理入群申请业务逻辑（同意/拒绝）
// =============================================================================
// 实现处理入群申请的核心业务逻辑
//
// 业务流程:
//   1. 从上下文获取当前用户 ID（处理人）
//   2. 调用 Social RPC 处理入群申请
//   3. 如果同意申请，调用 IM RPC 建立申请人与群组的会话
//   4. 返回处理结果
//
// 数据来源:
//   - JWT token: 当前用户 ID
//   - 请求参数: 申请记录 ID、群组 ID、处理结果
//
// =============================================================================
package group

import (
	"context"
	"imooc.com/easy-chat/apps/im/rpc/imclient"
	"imooc.com/easy-chat/apps/social/rpc/socialclient"
	"imooc.com/easy-chat/pkg/constants"
	"imooc.com/easy-chat/pkg/ctxdata"

	"imooc.com/easy-chat/apps/social/api/internal/svc"
	"imooc.com/easy-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// GroupPutInHandleLogic 入群申请处理业务逻辑处理器
type GroupPutInHandleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGroupPutInHandleLogic 创建入群申请处理业务逻辑处理器实例
//
// 参数:
//   - ctx: 请求上下文，包含用户认证信息等
//   - svcCtx: 服务上下文，包含 RPC 客户端等依赖
//
// 返回:
//   - *GroupPutInHandleLogic: 业务逻辑处理器实例
func NewGroupPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInHandleLogic {
	return &GroupPutInHandleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GroupPutInHandle 处理入群申请
// 对收到的入群申请进行处理（同意或拒绝）
//
// 处理流程:
//   1. 从 JWT token 上下文中提取当前用户 ID（处理人，通常是群主或管理员）
//   2. 调用 Social RPC 服务处理入群申请
//   3. 如果处理结果不是"通过"，直接返回
//   4. 如果同意申请（HandleResult=1），调用 IM RPC 服务建立申请人与群组的会话
//   5. 会话用于后续的群聊消息发送和接收
//
// 参数:
//   - req: 入群申请处理请求参数
//     - GroupReqId: 入群申请记录 ID
//     - GroupId: 群组 ID
//     - HandleResult: 处理结果（1-同意，2-拒绝）
//
// 返回:
//   - resp: 入群申请处理响应（空结构）
//   - err: 错误信息
func (l *GroupPutInHandleLogic) GroupPutInHandle(req *types.GroupPutInHandleRep) (resp *types.GroupPutInHandleResp, err error) {
	// 从上下文获取当前登录用户的 ID（处理人）
	uid := ctxdata.GetUId(l.ctx)

	// 调用 Social RPC 处理入群申请
	res, err := l.svcCtx.Social.GroupPutInHandle(l.ctx, &socialclient.GroupPutInHandleReq{
		GroupReqId:   req.GroupReqId,    // 入群申请记录 ID
		GroupId:      req.GroupId,       // 群组 ID
		HandleUid:    uid,               // 处理人 ID
		HandleResult: req.HandleResult,  // 处理结果（1-同意，2-拒绝）
	})

	// 如果处理结果不是"通过"，直接返回
	if constants.HandlerResult(req.HandleResult) != constants.PassHandlerResult {
		return
	}

	// 以下是通过申请后的业务逻辑

	// 检查群组 ID 是否有效
	if res.GroupId == "" {
		return nil, err
	}

	// 调用 IM RPC 建立申请人与群组的会话
	// 会话类型为群聊（GroupChatType）
	_, err = l.svcCtx.Im.SetUpUserConversation(l.ctx, &imclient.SetUpUserConversationReq{
		SendId:   uid,                                  // 发送者 ID（处理人）
		RecvId:   res.GroupId,                          // 接收者 ID（群组 ID）
		ChatType: int32(constants.GroupChatType),       // 会话类型：群聊
	})

	return nil, err
}
