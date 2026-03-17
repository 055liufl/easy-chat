// =============================================================================
// Ping 业务逻辑（RPC）
// =============================================================================
// 处理 Ping 请求的业务逻辑，用于:
//   - 健康检查
//   - 服务可用性测试
//   - RPC 连接测试
//
// 数据来源:
//   客户端通过 RPC 调用发送 Ping 请求
//
// 业务场景:
//   用于监控系统检查服务是否正常运行，或测试 RPC 连接是否可用
//
// 业务流程:
//   1. 接收 Ping 请求
//   2. 返回固定的 Pong 响应（包含服务标识）
//
// =============================================================================
package logic

import (
	"context"
	"imooc.com/easy-chat/apps/user/rpc/internal/svc"
	"imooc.com/easy-chat/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

// PingLogic Ping 业务逻辑结构
type PingLogic struct {
	ctx    context.Context    // 请求上下文
	svcCtx *svc.ServiceContext // 服务上下文（包含配置和依赖）
	logx.Logger                // 日志记录器
}

// NewPingLogic 创建 Ping 业务逻辑实例
//
// 参数:
//   - ctx: 请求上下文
//   - svcCtx: 服务上下文
//
// 返回:
//   - *PingLogic: 业务逻辑实例
func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Ping 执行 Ping 业务逻辑
//
// 参数:
//   - in: Ping 请求参数
//
// 返回:
//   - *user.Response: Pong 响应（包含服务标识）
//   - error: 错误信息
//
// 业务流程:
//   1. 接收 Ping 请求
//   2. 返回固定的 Pong 响应（Pong: "imooc.com"）
func (l *PingLogic) Ping(in *user.Request) (*user.Response, error) {
	// todo: add your logic here and delete this line

	return &user.Response{
		Pong: "imooc.com",
	}, nil
}
