// =============================================================================
// JWT 认证处理器 - WebSocket 连接认证
// =============================================================================
// 提供基于 JWT Token 的 WebSocket 连接认证功能，包括:
//   - Token 解析与验证
//   - 用户身份识别
//   - 上下文信息注入
//
// 认证流程:
//   1. 从 HTTP 请求中提取 JWT Token
//   2. 使用配置的密钥验证 Token 有效性
//   3. 解析 Token 中的用户身份信息
//   4. 将用户 ID 注入到请求上下文中
//
// 使用场景:
//   - WebSocket 握手阶段的身份验证
//   - 确保只有合法用户才能建立 WebSocket 连接
//
// =============================================================================

package handler

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/token"
	"imooc.com/easy-chat/apps/im/ws/internal/svc"
	"imooc.com/easy-chat/pkg/ctxdata"
	"net/http"
)

// JwtAuth JWT 认证器结构体
// 负责处理 WebSocket 连接的 JWT Token 认证
type JwtAuth struct {
	svc    *svc.ServiceContext  // 服务上下文，包含配置和依赖
	parser *token.TokenParser   // Token 解析器，用于解析 JWT
	logx.Logger                 // 日志记录器
}

// NewJwtAuth 创建 JWT 认证器实例
// 初始化 Token 解析器和日志记录器
//
// 参数:
//   - svc: 服务上下文，包含 JWT 密钥等配置信息
//
// 返回:
//   - *JwtAuth: 初始化完成的 JWT 认证器实例
func NewJwtAuth(svc *svc.ServiceContext) *JwtAuth {
	return &JwtAuth{
		svc:    svc,
		parser: token.NewTokenParser(),
		Logger: logx.WithContext(context.Background()),
	}
}

// Auth 执行 JWT Token 认证
// 从 HTTP 请求中解析并验证 JWT Token，提取用户身份信息
//
// 认证流程:
//  1. 使用 TokenParser 从请求头中提取 Token
//  2. 使用配置的 AccessSecret 验证 Token 签名
//  3. 检查 Token 是否有效（未过期、签名正确）
//  4. 解析 Token 中的 Claims，提取用户标识
//  5. 将用户标识注入到请求上下文中，供后续处理使用
//
// 参数:
//   - w: HTTP 响应写入器（当前未使用，保留用于扩展）
//   - r: HTTP 请求对象，包含待验证的 Token
//
// 返回:
//   - bool: 认证是否成功，true 表示认证通过，false 表示认证失败
func (j *JwtAuth) Auth(w http.ResponseWriter, r *http.Request) bool {
	// 解析 Token，使用配置的密钥进行验证
	tok, err := j.parser.ParseToken(r, j.svc.Config.JwtAuth.AccessSecret, "")
	if err != nil {
		j.Errorf("parse token err %v ", err)
		return false
	}

	// 检查 Token 是否有效（签名验证、过期时间等）
	if !tok.Valid {
		return false
	}

	// 提取 Token 中的 Claims（声明信息）
	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	// 将用户标识注入到请求上下文中，供后续业务逻辑使用
	*r = *r.WithContext(context.WithValue(r.Context(), ctxdata.Identify, claims[ctxdata.Identify]))

	return true
}

// UserId 从请求上下文中获取用户 ID
// 必须在 Auth 方法成功执行后调用，否则返回空字符串
//
// 参数:
//   - r: HTTP 请求对象，其上下文中包含用户身份信息
//
// 返回:
//   - string: 用户 ID，如果上下文中不存在则返回空字符串
func (j *JwtAuth) UserId(r *http.Request) string {
	return ctxdata.GetUId(r.Context())
}
