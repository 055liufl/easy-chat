// =============================================================================
// ResultX - HTTP 统一响应格式
// =============================================================================
// 提供统一的 HTTP 响应格式，确保所有 API 返回一致的 JSON 结构。
//
// 功能特性:
//   - 统一的成功/失败响应格式
//   - 自动解析 RPC 错误码和错误消息
//   - 支持 gRPC 错误和自定义错误码
//   - 自动记录错误日志
//
// 响应格式:
//   {
//       "code": 200,       // 业务状态码（200 成功，其他为错误码）
//       "msg":  "",        // 错误消息（成功时为空）
//       "data": {}         // 响应数据（失败时为 null）
//   }
//
// 使用场景:
//   - 所有 HTTP API 的响应格式化
//   - Go-Zero REST 服务的全局响应处理
//
// 项目中的应用:
//   - 所有 API 服务的 main.go 中配置全局响应处理器
//   - httpx.SetOkHandler(resultx.OkHandler)
//   - httpx.SetErrorHandlerCtx(resultx.ErrHandler("服务名"))
//
// 错误处理流程:
//   1. 业务逻辑返回 error
//   2. ErrHandler 解析错误类型（自定义错误码 / gRPC 错误）
//   3. 提取错误码和错误消息
//   4. 记录错误日志
//   5. 返回统一格式的 JSON 响应
//
// =============================================================================
package resultx

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	zrpcErr "github.com/zeromicro/x/errors"
	"google.golang.org/grpc/status"
	"net/http"

	"imooc.com/easy-chat/pkg/xerr"
)

// Response HTTP 统一响应结构体
// 所有 API 接口都使用此结构体返回数据
//
// 字段说明:
//   - Code: 业务状态码，200 表示成功，其他为错误码（见 xerr/err_code.go）
//   - Msg:  错误消息，成功时为空字符串
//   - Data: 响应数据，失败时为 nil
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Success 创建成功响应
// 返回 code=200 的成功响应，携带业务数据
//
// 参数:
//   - data: 响应数据（任意类型）
//
// 返回:
//   - *Response: 成功响应对象
//
// 示例:
//   return resultx.Success(userInfo)
//   // 输出: {"code": 200, "msg": "", "data": {"id": "123", "name": "张三"}}
func Success(data interface{}) *Response {
	return &Response{
		Code: 200,
		Msg:  "",
		Data: data,
	}
}

// Fail 创建失败响应
// 返回指定错误码和错误消息的失败响应
//
// 参数:
//   - code: 业务错误码
//   - err:  错误消息字符串
//
// 返回:
//   - *Response: 失败响应对象
//
// 示例:
//   return resultx.Fail(100001, "服务器异常")
//   // 输出: {"code": 100001, "msg": "服务器异常", "data": null}
func Fail(code int, err string) *Response {
	return &Response{
		Code: code,
		Msg:  err,
		Data: nil,
	}
}

// OkHandler Go-Zero 全局成功响应处理器
// 配合 httpx.SetOkHandler 使用，统一处理所有成功响应
//
// 参数:
//   - ctx: 上下文（未使用）
//   - v:   业务逻辑返回的数据
//
// 返回:
//   - any: 统一格式的成功响应
//
// 配置方式:
//   httpx.SetOkHandler(resultx.OkHandler)
func OkHandler(_ context.Context, v interface{}) any {
	return Success(v)
}

// ErrHandler Go-Zero 全局错误响应处理器
// 配合 httpx.SetErrorHandlerCtx 使用，统一处理所有错误响应
//
// 参数:
//   - name: 服务名称，用于日志标识（如 "user-api"、"social-api"）
//
// 返回:
//   - func: 错误处理函数，接收 ctx 和 err，返回 HTTP 状态码和响应体
//
// 错误解析优先级:
//   1. 先尝试解析为自定义错误码（*zrpcErr.CodeMsg）
//   2. 再尝试解析为 gRPC 错误（status.FromError）
//   3. 都不匹配则使用默认错误码 SERVER_COMMON_ERROR
//
// 配置方式:
//   httpx.SetErrorHandlerCtx(resultx.ErrHandler("user-api"))
func ErrHandler(name string) func(ctx context.Context, err error) (int, any) {
	return func(ctx context.Context, err error) (int, any) {
		// 默认错误码和错误消息
		errcode := xerr.SERVER_COMMON_ERROR
		errmsg := xerr.ErrMsg(errcode)

		// 获取根因错误
		causeErr := errors.Cause(err)
		// 尝试解析为自定义错误码（来自 RPC 服务返回的错误）
		if e, ok := causeErr.(*zrpcErr.CodeMsg); ok {
			errcode = e.Code
			errmsg = e.Msg
		} else {
			// 尝试解析为 gRPC 标准错误
			if gstatus, ok := status.FromError(causeErr); ok {
				errcode = int(gstatus.Code())
				errmsg = gstatus.Message()
			}
		}

		// 记录错误日志，包含服务名和完整错误信息
		logx.WithContext(ctx).Errorf("【%s】 err %v", name, err)

		// 返回 HTTP 400 状态码和统一格式的错误响应
		return http.StatusBadRequest, Fail(errcode, errmsg)
	}
}
