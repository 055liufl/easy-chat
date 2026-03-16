// =============================================================================
// IM Models - 错误定义模块
// =============================================================================
// 定义 IM 模块中使用的通用错误类型
//
// 错误类型:
//   - ErrNotFound: 数据未找到错误，来自 MongoDB 驱动
//   - ErrInvalidObjectId: 无效的 MongoDB ObjectID 格式错误
//
// =============================================================================
package immodels

import (
	"errors"

	"github.com/zeromicro/go-zero/core/stores/mon"
)

var (
	// ErrNotFound 数据未找到错误
	// 当查询的数据在 MongoDB 中不存在时返回此错误
	ErrNotFound = mon.ErrNotFound

	// ErrInvalidObjectId 无效的 ObjectID 错误
	// 当传入的字符串无法转换为有效的 MongoDB ObjectID 时返回此错误
	ErrInvalidObjectId = errors.New("invalid objectId")
)
