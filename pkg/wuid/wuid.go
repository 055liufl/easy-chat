// =============================================================================
// WUID - 全局唯一 ID 生成器
// =============================================================================
// 基于 WUID 库封装的全局唯一 ID 生成工具，使用 MySQL 存储 ID 段。
//
// 功能特性:
//   - 生成全局唯一的用户 ID（16 位十六进制格式）
//   - 组合两个用户 ID 生成会话 ID（用于私聊）
//   - 基于 MySQL 的分布式 ID 段分配，保证唯一性
//   - 懒加载初始化，首次调用时自动连接数据库
//
// 使用场景:
//   - 用户注册：生成唯一的用户 ID
//   - 私聊会话：组合两个用户 ID 生成会话标识
//
// 设计思路:
//   - 使用 WUID 的号段模式，从 MySQL 批量获取 ID 段
//   - 高 28 位从 MySQL 获取，低 36 位本地递增，减少数据库访问
//   - ID 格式为 16 位十六进制字符串（如 0x0000000100000001）
//   - CombineId 将两个 ID 排序后拼接，确保 A_B 和 B_A 生成相同的会话 ID
//
// 项目中的应用:
//   - User RPC 服务：用户注册时调用 GenUid 生成用户 ID
//   - IM 服务：创建私聊会话时调用 CombineId 生成会话 ID
//
// 第三方依赖:
//   - github.com/edwingeng/wuid: 高性能唯一 ID 生成库
//
// 数据库要求:
//   - MySQL 中需要有 wuid 表，用于存储 ID 段信息
//
// =============================================================================
package wuid

import (
	"database/sql"
	"fmt"
	"github.com/edwingeng/wuid/mysql/wuid"
	"sort"
	"strconv"
)

// w 全局 WUID 实例（单例模式）
// 整个应用共享一个 WUID 实例，保证 ID 的全局唯一性
var w *wuid.WUID

// Init 初始化 WUID 实例
// 连接 MySQL 数据库，加载 ID 段信息
//
// 参数:
//   - dsn: MySQL 数据源名称（如 "user:password@tcp(host:port)/dbname"）
//
// 工作流程:
//   1. 创建 MySQL 连接工厂函数
//   2. 创建 WUID 实例
//   3. 从 MySQL 的 wuid 表加载高 28 位 ID 段
//
// 注意:
//   - 通常不需要手动调用，GenUid 会自动初始化
//   - MySQL 中需要预先创建 wuid 表
func Init(dsn string) {

	// 创建数据库连接工厂函数
	// WUID 库需要一个返回 *sql.DB 的函数
	newDB := func() (*sql.DB, bool, error) {
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			return nil, false, err
		}
		// 返回 db 实例，true 表示需要关闭连接
		return db, true, nil
	}

	// 创建 WUID 实例并从 MySQL 加载 ID 段
	w = wuid.NewWUID("default", nil)
	_ = w.LoadH28FromMysql(newDB, "wuid")
}

// GenUid 生成全局唯一的用户 ID
// 返回 16 位十六进制格式的唯一 ID 字符串
//
// 参数:
//   - dsn: MySQL 数据源名称，用于首次初始化
//
// 返回:
//   - string: 16 位十六进制 ID（如 "0x0000000100000001"）
//
// 使用场景:
//   - 用户注册时生成唯一用户 ID
//
// 示例:
//   uid := wuid.GenUid("user:pass@tcp(localhost:3306)/easy_chat")
//   // 输出: "0x0000000100000001"
//
// 注意:
//   - 首次调用会自动初始化 WUID（懒加载）
//   - 生成的 ID 全局唯一，单调递增
func GenUid(dsn string) string {
	if w == nil {
		Init(dsn)
	}

	return fmt.Sprintf("%#016x", w.Next())
}

// CombineId 组合两个用户 ID 生成会话 ID
// 将两个用户 ID 排序后用下划线拼接，确保相同的两个用户始终生成相同的会话 ID
//
// 参数:
//   - aid: 用户 A 的 ID
//   - bid: 用户 B 的 ID
//
// 返回:
//   - string: 组合后的会话 ID（格式: "较小ID_较大ID"）
//
// 使用场景:
//   - 创建私聊会话时生成唯一的会话标识
//   - 查询两个用户之间的私聊记录
//
// 示例:
//   id := wuid.CombineId("0x0000000100000002", "0x0000000100000001")
//   // 输出: "0x0000000100000001_0x0000000100000002"（较小的在前）
//
//   // 交换顺序结果相同
//   id2 := wuid.CombineId("0x0000000100000001", "0x0000000100000002")
//   // 输出: "0x0000000100000001_0x0000000100000002"
//
// 设计思路:
//   - 先将两个 ID 转为无符号整数进行比较
//   - 较小的 ID 放在前面，确保 CombineId(A,B) == CombineId(B,A)
//   - 这样无论谁发起私聊，都能找到同一个会话
func CombineId(aid, bid string) string {
	ids := []string{aid, bid}

	// 按数值大小排序（将十六进制字符串转为无符号整数比较）
	sort.Slice(ids, func(i, j int) bool {
		a, _ := strconv.ParseUint(ids[i], 0, 64)
		b, _ := strconv.ParseUint(ids[j], 0, 64)
		return a < b
	})

	// 拼接为 "较小ID_较大ID" 格式
	return fmt.Sprintf("%s_%s", ids[0], ids[1])
}
