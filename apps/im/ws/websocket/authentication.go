// =============================================================================
// WebSocket Authentication - WebSocket 身份认证
// =============================================================================
// 定义 WebSocket 连接的身份认证接口和默认实现
//
// 认证流程:
//  1. 客户端发起 WebSocket 连接请求
//  2. 服务器调用 Auth 方法验证请求是否合法
//  3. 认证通过后，调用 UserId 方法提取用户 ID
//  4. 建立用户 ID 与连接的映射关系
//
// 默认实现:
//   - Auth: 始终返回 true（无认证）
//   - UserId: 从 URL 查询参数中提取 userId，如果不存在则使用时间戳
//
// 自定义认证:
//   实现 Authentication 接口，可以集成 JWT、OAuth 等认证方式
//
// =============================================================================
package websocket

import (
	"fmt"
	"net/http"
	"time"
)

// Authentication 身份认证接口
// 定义了 WebSocket 连接的认证和用户识别方法
type Authentication interface {
	// Auth 验证连接请求是否合法
	// 可以在此方法中验证 token、session 等认证信息
	//
	// 参数:
	//   - w: HTTP 响应写入器，可用于返回认证失败信息
	//   - r: HTTP 请求对象，包含认证信息（如 token、cookie 等）
	//
	// 返回:
	//   - bool: true 表示认证通过，false 表示认证失败
	Auth(w http.ResponseWriter, r *http.Request) bool

	// UserId 从请求中提取用户 ID
	// 用于建立用户 ID 与连接的映射关系
	//
	// 参数:
	//   - r: HTTP 请求对象
	//
	// 返回:
	//   - string: 用户 ID
	UserId(r *http.Request) string
}

// authentication 默认的身份认证实现
// 提供了最基本的认证逻辑，适用于开发和测试环境
type authentication struct{}

// Auth 默认认证实现
// 始终返回 true，不进行任何认证检查
// 生产环境应该实现自己的认证逻辑
//
// 参数:
//   - w: HTTP 响应写入器
//   - r: HTTP 请求对象
//
// 返回:
//   - bool: 始终返回 true
func (*authentication) Auth(w http.ResponseWriter, r *http.Request) bool {
	return true
}

// UserId 默认用户 ID 提取实现
// 从 URL 查询参数中提取 userId
// 如果不存在则使用当前时间戳作为用户 ID
//
// 参数:
//   - r: HTTP 请求对象
//
// 返回:
//   - string: 用户 ID
func (*authentication) UserId(r *http.Request) string {
	query := r.URL.Query()
	if query != nil && query["userId"] != nil {
		return fmt.Sprintf("%v", query["userId"])
	}

	// 如果没有提供 userId，使用时间戳作为临时 ID
	return fmt.Sprintf("%v", time.Now().UnixMilli())
}