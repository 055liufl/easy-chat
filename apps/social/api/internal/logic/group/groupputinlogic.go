// =============================================================================
// 入群申请 Logic - 发送入群申请业务逻辑
// =============================================================================
// 实现发送入群申请的核心业务逻辑
//
// 业务流程:
//   1. 从上下文获取当前用户 ID（申请人）
//   2. 调用 Social RPC 创建入群申请记录
//   3. 调用 IM RPC 建立用户与群组的会话
//   4. 返回处理结果
//
// 数据来源:
//   - JWT token: 当前用户 ID
//   - 请求参数: 群组 ID、申请消息、申请时间、加入来源
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

// GroupPutInLogic 入群申请业务逻辑处理器
type GroupPutInLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGroupPutInLogic 创建入群申请业务逻辑处理器实例
//
// 参数:
//   - ctx: 请求上下文，包含用户认证信息等
//   - svcCtx: 服务上下文，包含 RPC 客户端等依赖
//
// 返回:
//   - *GroupPutInLogic: 业务逻辑处理器实例
func NewGroupPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInLogic {
	return &GroupPutInLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GroupPutIn 发送入群申请
// 向指定群组发送入群申请，并建立会话
//
// 处理流程:
//   1. 从 JWT token 上下文中提取当前用户 ID（申请人）
//   2. 调用 Social RPC 服务创建入群申请记录
//   3. 如果申请创建成功，调用 IM RPC 服务建立用户与群组的会话
//   4. 会话用于后续的群聊消息发送和接收
//
// 参数:
//   - req: 入群申请请求参数
//     - GroupId: 群组 ID
//     - ReqMsg: 申请消息（如"我想加入这个群"）
//     - ReqTime: 申请时间戳
//     - JoinSource: 加入来源（如搜索、邀请等）
//
// 返回:
//   - resp: 入群申请响应（空结构）
//   - err: 错误信息
func (l *GroupPutInLogic) GroupPutIn(req *types.GroupPutInRep) (resp *types.GroupPutInResp, err error) {
	// 从上下文获取当前登录用户的 ID（申请人）
	uid := ctxdata.GetUId(l.ctx)

	// 调用 Social RPC 创建入群申请记录
	res, err := l.svcCtx.Social.GroupPutin(l.ctx, &socialclient.GroupPutinReq{
		GroupId:    req.GroupId,             // 群组 ID
		ReqId:      uid,                     // 申请人 ID
		ReqMsg:     req.ReqMsg,              // 申请消息
		ReqTime:    req.ReqTime,             // 申请时间
		JoinSource: int32(req.JoinSource),   // 加入来源
	})

	// 检查申请是否创建成功
	if res.GroupId == "" {
		return nil, err
	}

	// 调用 IM RPC 建立用户与群组的会话
	// 会话类型为群聊（GroupChatType）
	_, err = l.svcCtx.Im.SetUpUserConversation(l.ctx, &imclient.SetUpUserConversationReq{
		SendId:   uid,                                  // 发送者 ID（申请人）
		RecvId:   res.GroupId,                          // 接收者 ID（群组 ID）
		ChatType: int32(constants.GroupChatType),       // 会话类型：群聊
	})

	return nil, err
}
