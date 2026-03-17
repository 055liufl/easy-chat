// =============================================================================
// Token - JWT Token 生成工具
// =============================================================================
// 提供 JWT Token 的生成功能，用于用户认证和授权。
//
// 功能特性:
//   - 生成包含用户 ID 的 JWT Token
//   - 支持自定义过期时间
//   - 使用 HS256 签名算法
//
// 使用场景:
//   - 用户登录：生成访问令牌
//   - 用户注册：生成初始令牌
//   - Token 刷新：生成新的令牌
//
// 设计思路:
//   - 使用 JWT 标准格式，便于跨服务验证
//   - 将用户 ID 存储在 Claims 中，便于提取
//   - 使用统一的键名（Identify），避免冲突
//
// 项目中的应用:
//   - User 服务：用户登录时生成 Token
//   - API 网关：验证 Token 并提取用户信息
//   - 所有服务：从 Token 中获取用户身份
//
// 第三方依赖:
//   - github.com/golang-jwt/jwt: JWT Token 生成和解析库
//
// =============================================================================
package ctxdata

import "github.com/golang-jwt/jwt"

// Identify Context 中存储用户 ID 的键名
// 用于在 Context 和 JWT Claims 中标识用户 ID
// 使用域名作为键名，避免与其他库冲突
const Identify = "imooc.com"

// GetJwtToken 生成 JWT Token
// 根据用户 ID 和过期时间生成 JWT Token
//
// 参数:
//   - secretKey: JWT 签名密钥，用于签名和验证 Token
//   - iat: Token 签发时间（Unix 时间戳，秒）
//   - seconds: Token 有效期（秒），例如 3600 表示 1 小时
//   - uid: 用户 ID，将存储在 Token 的 Claims 中
//
// 返回:
//   - string: 生成的 JWT Token 字符串
//   - error: 生成 Token 时的错误
//
// Token 结构:
//   - exp: 过期时间（iat + seconds）
//   - iat: 签发时间
//   - imooc.com: 用户 ID（自定义 Claim）
//
// 使用场景:
//   - 用户登录成功后生成 Token
//   - Token 过期后刷新生成新 Token
//
// 示例:
//   token, err := ctxdata.GetJwtToken(
//       "your-secret-key",
//       time.Now().Unix(),
//       3600,  // 1 小时有效期
//       "user_123",
//   )
//   if err != nil {
//       return err
//   }
//   // 返回 Token 给客户端
//
// 注意:
//   - secretKey 必须保密，不能泄露
//   - 建议使用环境变量或配置文件管理 secretKey
//   - Token 过期后需要重新登录或刷新
func GetJwtToken(secretKey string, iat, seconds int64, uid string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims[Identify] = uid

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims

	return token.SignedString([]byte(secretKey))
}
