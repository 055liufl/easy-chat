// =============================================================================
// Social Constants - 社交模块常量定义
// =============================================================================
// 定义社交模块中使用的各种常量，包括好友申请处理结果、群角色等级、进群方式等。
//
// 功能特性:
//   - 好友申请处理结果定义（未处理、通过、拒绝、取消）
//   - 群角色等级定义（创建者、管理员、普通成员）
//   - 进群方式定义（邀请、申请）
//
// 使用场景:
//   - 好友管理：处理好友申请
//   - 群组管理：管理群成员角色和权限
//   - 群组加入：区分邀请和申请两种方式
//
// 设计思路:
//   - 使用自定义类型而非直接使用 int，提高代码可读性
//   - 使用 iota 自动递增，便于扩展新类型
//   - 从 1 开始，避免零值带来的歧义
//
// 项目中的应用:
//   - Social 服务：好友申请、群组管理
//   - IM 服务：消息权限控制
//   - User 服务：用户关系管理
//
// =============================================================================
package constants

// HandlerResult 处理结果
// 定义好友申请或群申请的处理结果
type HandlerResult int

const (
	NoHandlerResult     HandlerResult = iota + 1 // 未处理（待审核）
	PassHandlerResult                            // 通过（同意申请）
	RefuseHandlerResult                          // 拒绝（拒绝申请）
	CancelHandlerResult                          // 取消（撤销申请）
)

// GroupRoleLevel 群角色等级
// 定义群成员的角色等级，用于权限控制
type GroupRoleLevel int

const (
	CreatorGroupRoleLevel GroupRoleLevel = iota + 1 // 创建者（群主，拥有最高权限）
	ManagerGroupRoleLevel                           // 管理员（协助群主管理群组）
	AtLargeGroupRoleLevel                           // 普通成员（无管理权限）
	// 注意: 从 1 开始是为了避免零值（0）带来的歧义
	// 零值可能表示未设置或无效状态，从 1 开始更明确
)

// GroupJoinSource 进群申请的方式
// 定义用户加入群组的方式，用于区分不同的加入场景
type GroupJoinSource int

const (
	InviteGroupJoinSource GroupJoinSource = iota + 1 // 邀请（被群成员邀请加入）
	PutInGroupJoinSource                             // 申请（主动申请加入群组）
)
