// =============================================================================
// 获取聊天记录业务逻辑
// =============================================================================
// 实现聊天记录查询的核心业务逻辑
//
// 功能说明:
//   - 查询指定会话的聊天记录
//   - 支持按时间范围过滤
//   - 支持限制返回数量
//
// 业务流程:
//  1. 调用 IM RPC 服务查询聊天记录
//  2. 将 RPC 响应数据转换为 API 响应格式
//  3. 返回聊天记录列表
//
// =============================================================================
package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"imooc.com/easy-chat/apps/im/rpc/imclient"

	"imooc.com/easy-chat/apps/im/api/internal/svc"
	"imooc.com/easy-chat/apps/im/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// GetChatLogLogic 获取聊天记录业务逻辑结构
type GetChatLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetChatLogLogic 创建获取聊天记录业务逻辑实例
//
// 参数:
//   - ctx: 上下文对象，用于传递请求信息和控制超时
//   - svcCtx: 服务上下文，包含所有依赖项（RPC 客户端等）
//
// 返回:
//   - *GetChatLogLogic: 业务逻辑实例
func NewGetChatLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatLogLogic {
	return &GetChatLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetChatLog 获取聊天记录
// 查询指定会话的聊天记录列表
//
// 参数:
//   - req: 请求参数，包含会话 ID、时间范围、查询数量等
//
// 返回:
//   - resp: 响应数据，包含聊天记录列表
//   - err: 错误信息
//
// 业务逻辑:
//  1. 调用 IM RPC 服务查询聊天记录
//  2. 使用 copier 库将 RPC 响应数据复制到 API 响应结构
//  3. 返回聊天记录列表
func (l *GetChatLogLogic) GetChatLog(req *types.ChatLogReq) (resp *types.ChatLogResp, err error) {
	// 步骤 1: 调用 IM RPC 服务查询聊天记录
	data, err := l.svcCtx.GetChatLog(l.ctx, &imclient.GetChatLogReq{
		ConversationId: req.ConversationId,
		StartSendTime:  req.StartSendTime,
		EndSendTime:    req.EndSendTime,
		Count:          req.Count,
	})
	if err != nil {
		return nil, err
	}

	// 步骤 2: 将 RPC 响应数据复制到 API 响应结构
	var res types.ChatLogResp
	copier.Copy(&res, data)

	// 步骤 3: 返回聊天记录列表
	return &res, err
}
