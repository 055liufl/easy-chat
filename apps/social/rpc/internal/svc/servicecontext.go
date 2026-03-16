// =============================================================================
// Social RPC 服务上下文模块
// =============================================================================
// 管理社交服务的全局依赖和资源，包括:
//   - 配置信息
//   - 数据库模型实例（好友、群组等）
//   - 数据库连接池
//
// 职责:
//   作为依赖注入容器，为各个业务逻辑层提供统一的资源访问入口
//
// 生命周期:
//   在服务启动时创建，服务运行期间保持单例，服务关闭时释放资源
// =============================================================================
package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"imooc.com/easy-chat/apps/social/rpc/internal/config"
	"imooc.com/easy-chat/apps/social/socialmodels"
)

// ServiceContext 服务上下文结构
// 包含服务运行所需的所有依赖资源
type ServiceContext struct {
	Config config.Config // 服务配置

	// 数据模型层（嵌入式字段，可直接访问模型方法）
	socialmodels.FriendsModel        // 好友关系模型
	socialmodels.FriendRequestsModel // 好友申请模型
	socialmodels.GroupsModel         // 群组信息模型
	socialmodels.GroupRequestsModel  // 群组申请模型
	socialmodels.GroupMembersModel   // 群组成员模型
}

// NewServiceContext 创建服务上下文实例
// 初始化数据库连接和所有数据模型
//
// 参数:
//   - c: 服务配置对象
//
// 返回:
//   - *ServiceContext: 初始化完成的服务上下文实例
//
// 初始化流程:
//  1. 创建 MySQL 数据库连接池
//  2. 初始化各个数据模型（带缓存）
//  3. 返回服务上下文实例
func NewServiceContext(c config.Config) *ServiceContext {

	// 创建 MySQL 连接池
	sqlConn := sqlx.NewMysql(c.Mysql.DataSource)

	return &ServiceContext{
		Config: c,

		// 初始化数据模型（带 Redis 缓存）
		FriendsModel:        socialmodels.NewFriendsModel(sqlConn, c.Cache),
		FriendRequestsModel: socialmodels.NewFriendRequestsModel(sqlConn, c.Cache),
		GroupsModel:         socialmodels.NewGroupsModel(sqlConn, c.Cache),
		GroupRequestsModel:  socialmodels.NewGroupRequestsModel(sqlConn, c.Cache),
		GroupMembersModel:   socialmodels.NewGroupMembersModel(sqlConn, c.Cache),
	}
}
