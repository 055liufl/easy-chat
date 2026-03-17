// =============================================================================
// Log Interceptor - RPC 服务端日志拦截器
// =============================================================================
// 提供 gRPC 服务端统一的错误日志记录和错误码转换功能。
//
// 功能特性:
//   - 统一错误日志记录
//   - 错误码转换（将业务错误码转换为 gRPC 错误码）
//   - 错误堆栈追踪
//   - 错误信息标准化
//
// 使用场景:
//   - 所有 RPC 服务端
//   - 统一错误处理
//   - 错误日志记录
//   - 错误监控和告警
//
// 设计思路:
//   - 拦截所有 RPC 请求的响应
//   - 如果响应包含错误，记录日志
//   - 将业务错误码转换为 gRPC 标准错误码
//   - 保留原始错误信息
//
// 项目中的应用:
//   - 所有 RPC 服务端
//   - 错误日志统一管理
//   - 错误监控和分析
//
// 第三方依赖:
//   - github.com/pkg/errors: 错误堆栈追踪
//   - github.com/zeromicro/x/errors: 业务错误码定义
//
// =============================================================================
package rpcserver

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	zerr "github.com/zeromicro/x/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LogInterceptor RPC 服务端日志拦截器
// 记录所有 RPC 请求的错误日志，并转换错误码
//
// 参数:
//   - ctx: 上下文对象
//   - req: 请求参数
//   - info: RPC 方法信息
//   - handler: RPC 处理函数
//
// 返回:
//   - resp: 响应结果
//   - err: 错误信息（已转换为 gRPC 错误）
//
// 工作流程:
//   1. 执行 RPC 处理函数
//   2. 如果没有错误，直接返回
//   3. 如果有错误，记录错误日志
//   4. 提取错误的根因（使用 errors.Cause）
//   5. 如果是业务错误码，转换为 gRPC 错误码
//   6. 返回转换后的错误
//
// 使用场景:
//   - 所有 RPC 服务端
//
// 示例:
//   server := grpc.NewServer(
//       grpc.UnaryInterceptor(rpcserver.LogInterceptor),
//   )
//
// 错误转换规则:
//   - 业务错误码（zerr.CodeMsg）-> gRPC 错误码
//   - 保留原始错误信息
func LogInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any,
	err error) {
	// 执行 RPC 处理函数
	resp, err = handler(ctx, req)
	if err == nil {
		return resp, nil
	}

	// 记录错误日志
	logx.WithContext(ctx).Errorf("【RPC SRV ERR】 %v", err)

	// 提取错误的根因
	causeErr := errors.Cause(err)
	// 如果是业务错误码，转换为 gRPC 错误码
	if e, ok := causeErr.(*zerr.CodeMsg); ok {
		err = status.Error(codes.Code(e.Code), e.Msg)
	}

	return resp, err
}