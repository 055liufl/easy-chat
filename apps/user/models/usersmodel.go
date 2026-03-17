// =============================================================================
// 用户数据模型
// =============================================================================
// 定义用户数据模型的接口和实现，包括:
//   - 用户表的 CRUD 操作
//   - 自定义查询方法（可扩展）
//
// 数据来源:
//   MySQL 数据库 users 表
//
// 业务场景:
//   为用户相关的业务逻辑提供数据访问层，支持缓存
//
// 注意:
//   - UsersModel 接口用于定义自定义方法
//   - customUsersModel 结构用于实现自定义方法
//   - 基础 CRUD 方法由 defaultUsersModel 提供（自动生成）
//
// =============================================================================
package models

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UsersModel = (*customUsersModel)(nil)

type (
	// UsersModel 用户数据模型接口
	// 可在此接口中添加自定义方法，并在 customUsersModel 中实现
	UsersModel interface {
		usersModel // 继承自动生成的基础接口（包含 CRUD 方法）
	}

	// customUsersModel 自定义用户数据模型实现
	// 用于实现 UsersModel 接口中的自定义方法
	customUsersModel struct {
		*defaultUsersModel // 继承自动生成的默认实现（包含基础 CRUD 方法）
	}
)

// NewUsersModel 创建用户数据模型实例
// 返回带缓存的用户数据模型
//
// 参数:
//   - conn: 数据库连接
//   - c: 缓存配置
//
// 返回:
//   - UsersModel: 用户数据模型实例
//
// 特性:
//   - 支持 Redis 缓存（提高查询性能）
//   - 自动缓存失效（数据更新时）
func NewUsersModel(conn sqlx.SqlConn, c cache.CacheConf) UsersModel {
	return &customUsersModel{
		defaultUsersModel: newUsersModel(conn, c),
	}
}
