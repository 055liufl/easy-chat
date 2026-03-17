// =============================================================================
// XErr - 错误创建工具
// =============================================================================
// 提供便捷的错误创建函数，封装 go-zero/x 的 errors 包。
//
// 功能特性:
//   - 创建带错误码的业务错误
//   - 提供常用错误的快捷创建方法
//   - 错误可通过 RPC 传递，客户端能解析错误码
//
// 使用场景:
//   - RPC 服务中返回业务错误
//   - 业务逻辑中抛出特定类型的错误
//
// 设计思路:
//   - 封装 go-zero/x 的 errors.CodeMsg，支持错误码传递
//   - 提供 New / NewMsg / NewDBErr / NewInternalErr 四种创建方式
//   - 错误通过 gRPC 传递后，客户端可通过 resultx.ErrHandler 解析
//
// 错误传递流程:
//   1. RPC 服务: return nil, xerr.NewDBErr()
//   2. gRPC 传输: 错误码和消息通过 gRPC status 传递
//   3. API 层: resultx.ErrHandler 解析错误码和消息
//   4. 客户端: 收到 {"code": 100003, "msg": "数据库繁忙，稍后再尝试"}
//
// 第三方依赖:
//   - github.com/zeromicro/x/errors: Go-Zero 扩展错误包，支持错误码
//
// =============================================================================
package xerr

import "github.com/zeromicro/x/errors"

// New 创建带错误码和自定义消息的业务错误
//
// 参数:
//   - code: 业务错误码（见 err_code.go）
//   - msg:  自定义错误消息
//
// 返回:
//   - error: 包含错误码的错误对象
//
// 示例:
//   return xerr.New(xerr.REQUEST_PARAM_ERROR, "手机号格式不正确")
func New(code int, msg string) error {
	return errors.New(code, msg)
}

// NewMsg 创建带自定义消息的通用服务器错误
// 错误码固定为 SERVER_COMMON_ERROR（100001）
//
// 参数:
//   - msg: 自定义错误消息
//
// 返回:
//   - error: 包含通用错误码的错误对象
//
// 示例:
//   return xerr.NewMsg("用户不存在")
func NewMsg(msg string) error {
	return errors.New(SERVER_COMMON_ERROR, msg)
}

// NewDBErr 创建数据库错误
// 错误码为 DB_ERROR（100003），消息为预定义的数据库错误提示
//
// 返回:
//   - error: 数据库错误对象
//
// 使用场景:
//   - 数据库查询失败
//   - 数据库连接异常
//   - 数据库写入失败
//
// 示例:
//   user, err := l.svcCtx.UserModel.FindOne(ctx, uid)
//   if err != nil {
//       return nil, xerr.NewDBErr()
//   }
func NewDBErr() error {
	return errors.New(DB_ERROR, ErrMsg(DB_ERROR))
}

// NewInternalErr 创建内部服务器错误
// 错误码为 SERVER_COMMON_ERROR（100001），消息为预定义的通用错误提示
//
// 返回:
//   - error: 内部服务器错误对象
//
// 使用场景:
//   - 未预期的内部错误
//   - 第三方服务调用失败
//   - 其他无法归类的错误
//
// 示例:
//   token, err := ctxdata.GetJwtToken(secretKey, now, expire, uid)
//   if err != nil {
//       return nil, xerr.NewInternalErr()
//   }
func NewInternalErr() error {
	return errors.New(SERVER_COMMON_ERROR, ErrMsg(SERVER_COMMON_ERROR))
}
