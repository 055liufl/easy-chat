// =============================================================================
// CtxData - 上下文数据处理工具
// =============================================================================
// 提供从 Context 中获取用户信息的工具函数。
//
// 功能特性:
//   - 从 Context 中提取用户 ID
//   - 类型安全的数据获取
//   - 统一的上下文数据访问接口
//
// 使用场景:
//   - HTTP 请求处理：获取当前登录用户 ID
//   - RPC 调用：传递用户身份信息
//   - 日志记录：记录操作用户
//   - 权限验证：获取用户身份进行权限检查
//
// 设计思路:
//   - 使用 Context 传递用户信息，避免全局变量
//   - 提供类型安全的获取方法，避免类型断言错误
//   - 统一的键名（Identify），避免冲突
//
// 项目中的应用:
//   - 所有 HTTP API：获取当前登录用户
//   - 所有 RPC 服务：获取调用者身份
//   - 业务逻辑层：获取操作用户进行权限控制
//
// 配合使用:
//   - ctxdata.GetJwtToken: 生成包含用户 ID 的 JWT Token
//   - JWT 中间件: 解析 Token 并将用户 ID 存入 Context
//   - GetUId: 从 Context 中获取用户 ID
//
// =============================================================================
package ctxdata

import "context"

// GetUId 从 Context 中获取用户 ID
// 从上下文中提取当前登录用户的 ID
//
// 参数:
//   - ctx: 上下文对象，包含用户身份信息
//
// 返回:
//   - string: 用户 ID，如果不存在则返回空字符串
//
// 使用场景:
//   - HTTP 请求处理中获取当前登录用户
//   - RPC 服务中获取调用者身份
//   - 业务逻辑中进行权限验证
//
// 示例:
//   func (l *GetUserInfoLogic) GetUserInfo(ctx context.Context) (*types.User, error) {
//       uid := ctxdata.GetUId(ctx)
//       if uid == "" {
//           return nil, errors.New("未登录")
//       }
//       // 查询用户信息
//       return l.svcCtx.UserModel.FindOne(ctx, uid)
//   }
//
// 注意:
//   - 如果 Context 中不存在用户 ID，返回空字符串
//   - 调用前需要确保 JWT 中间件已经将用户 ID 存入 Context
func GetUId(ctx context.Context) string {
	if u, ok := ctx.Value(Identify).(string); ok {
		return u
	}
	return ""
}
