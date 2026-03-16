// =============================================================================
// 群组创建逻辑
// =============================================================================
// 实现群组创建功能，包括:
//   - 生成唯一群组 ID
//   - 创建群组信息记录
//   - 添加创建者为群主
//
// 业务流程:
//  1. 生成全局唯一的群组 ID
//  2. 创建群组信息记录（名称、图标、创建者等）
//  3. 添加创建者为群成员，角色为群主
//  4. 使用事务保证数据一致性
//
// 数据流:
//   创建请求 -> 生成群组 ID -> 创建群组记录 -> 添加群主成员 -> 响应群组 ID
//
// 事务保证:
//   群组创建和群主成员添加在同一事务中，保证数据一致性
// =============================================================================
package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"imooc.com/easy-chat/apps/social/socialmodels"
	"imooc.com/easy-chat/pkg/constants"
	"imooc.com/easy-chat/pkg/wuid"
	"imooc.com/easy-chat/pkg/xerr"
	"time"

	"imooc.com/easy-chat/apps/social/rpc/internal/svc"
	"imooc.com/easy-chat/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

// GroupCreateLogic 群组创建逻辑处理器
type GroupCreateLogic struct {
	ctx    context.Context        // 请求上下文
	svcCtx *svc.ServiceContext     // 服务上下文
	logx.Logger                    // 日志记录器
}

// NewGroupCreateLogic 创建群组创建逻辑处理器实例
//
// 参数:
//   - ctx: 请求上下文
//   - svcCtx: 服务上下文
//
// 返回:
//   - *GroupCreateLogic: 初始化完成的逻辑处理器实例
func NewGroupCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupCreateLogic {
	return &GroupCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GroupCreate 创建群组
//
// 参数:
//   - in: 群组创建请求，包含群名称、图标、创建者 ID
//
// 返回:
//   - *social.GroupCreateResp: 群组创建响应，包含群组 ID
//   - error: 错误信息（数据库错误等）
//
// 业务逻辑:
//  1. 生成全局唯一的群组 ID（基于数据源的分布式 ID）
//  2. 创建群组信息记录
//  3. 使用事务添加创建者为群成员，角色为群主
func (l *GroupCreateLogic) GroupCreate(in *social.GroupCreateReq) (*social.GroupCreateResp, error) {
	// todo: add your logic here and delete this line

	// 构建群组信息
	groups := &socialmodels.Groups{
		Id:         wuid.GenUid(l.svcCtx.Config.Mysql.DataSource), // 生成全局唯一 ID
		Name:       in.Name,                                        // 群名称
		Icon:       in.Icon,                                        // 群图标
		CreatorUid: in.CreatorUid,                                  // 创建者 ID
		//IsVerify:   true,                                         // 是否需要验证（已注释）
		IsVerify: false, // 是否需要验证（当前设置为不需要验证）
	}

	// 使用事务：创建群组 -> 添加群主成员
	err := l.svcCtx.GroupsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// 插入群组记录
		_, err := l.svcCtx.GroupsModel.Insert(l.ctx, session, groups)

		if err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "insert group err %v req %v", err, in)
		}

		// 添加创建者为群成员，角色为群主
		_, err = l.svcCtx.GroupMembersModel.Insert(l.ctx, session, &socialmodels.GroupMembers{
			GroupId:   groups.Id,                                 // 群组 ID
			UserId:    in.CreatorUid,                             // 用户 ID（创建者）
			RoleLevel: int(constants.CreatorGroupRoleLevel),      // 角色等级（群主）
		})
		if err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "insert group member err %v req %v", err, in)
		}
		return nil
	})

	// 延迟 2 秒（可能用于测试或其他目的）
	time.Sleep(2 * time.Second)

	return &social.GroupCreateResp{
		Id: groups.Id, // 返回群组 ID
	}, err
}
