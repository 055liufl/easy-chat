// =============================================================================
// 服务上下文 - 依赖注入容器
// =============================================================================
// 管理 IM API 服务的所有依赖项，包括:
//   - 配置对象
//   - RPC 客户端（IM、User、Social）
//
// 作用:
//   - 统一管理服务依赖
//   - 提供给所有 Handler 和 Logic 使用
//   - 实现依赖注入模式
//
// =============================================================================
package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"imooc.com/easy-chat/apps/im/api/internal/config"
	"imooc.com/easy-chat/apps/im/rpc/imclient"
	"imooc.com/easy-chat/apps/social/rpc/socialclient"
	"imooc.com/easy-chat/apps/user/rpc/userclient"
)

// ServiceContext 服务上下文结构
// 包含所有业务逻辑需要的依赖项
type ServiceContext struct {
	// Config 服务配置对象
	Config config.Config

	// Im IM RPC 客户端
	// 提供聊天记录、会话管理等功能
	imclient.Im

	// User User RPC 客户端
	// 提供用户信息查询等功能
	userclient.User

	// Social Social RPC 客户端
	// 提供群组信息、好友关系等功能
	socialclient.Social
}

// NewServiceContext 创建服务上下文实例
// 初始化所有 RPC 客户端连接
//
// 参数:
//   - c: 服务配置对象，包含 RPC 服务地址等信息
//
// 返回:
//   - *ServiceContext: 初始化完成的服务上下文实例
//
// 初始化流程:
//  1. 创建 IM RPC 客户端连接
//  2. 创建 User RPC 客户端连接
//  3. 创建 Social RPC 客户端连接
//  4. 返回服务上下文对象
func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		Im:     imclient.NewIm(zrpc.MustNewClient(c.ImRpc)),
		User:   userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		Social: socialclient.NewSocial(zrpc.MustNewClient(c.SocialRpc)),
	}
}
