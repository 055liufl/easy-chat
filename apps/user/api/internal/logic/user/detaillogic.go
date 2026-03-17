// =============================================================================
// 用户详情业务逻辑
// =============================================================================
// 处理获取用户详情的业务逻辑，包括:
//   - 从请求上下文中获取用户 ID（JWT token 解析）
//   - 调用 User RPC 服务获取用户详细信息
//   - 数据转换（RPC 层类型 -> API 层类型）
//   - 返回用户详细信息
//
// 数据来源:
//   用户 ID 从 JWT token 中解析，用户信息从 User RPC 服务获取
//
// 业务场景:
//   已登录用户获取自己的详细信息（昵称、头像、性别等）
//
// 业务流程:
//   1. 从请求上下文中获取用户 ID（JWT 中间件已解析）
//   2. 调用 User RPC 服务的 GetUserInfo 方法获取用户信息
//   3. 将 RPC 响应转换为 API 响应格式
//   4. 返回用户详细信息
//
// =============================================================================
package user

import (
	"context"
	"github.com/jinzhu/copier"
	"imooc.com/easy-chat/apps/user/rpc/user"
	"imooc.com/easy-chat/pkg/ctxdata"

	"imooc.com/easy-chat/apps/user/api/internal/svc"
	"imooc.com/easy-chat/apps/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// DetailLogic 用户详情业务逻辑结构
type DetailLogic struct {
	logx.Logger                    // 日志记录器
	ctx         context.Context    // 请求上下文
	svcCtx      *svc.ServiceContext // 服务上下文（包含配置和依赖）
}

// NewDetailLogic 创建用户详情业务逻辑实例
//
// 参数:
//   - ctx: 请求上下文
//   - svcCtx: 服务上下文
//
// 返回:
//   - *DetailLogic: 业务逻辑实例
func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Detail 执行获取用户详情业务逻辑
//
// 参数:
//   - req: 用户信息请求参数（实际无参数，用户 ID 从 context 中获取）
//
// 返回:
//   - resp: 用户信息响应（包含用户详细信息）
//   - err: 错误信息
//
// 业务流程:
//   1. 从请求上下文中获取用户 ID（JWT 中间件已将用户 ID 写入 context）
//   2. 调用 User RPC 服务的 GetUserInfo 方法获取用户信息
//   3. 使用 copier 将 RPC 响应转换为 API 响应格式
//   4. 返回用户详细信息
func (l *DetailLogic) Detail(req *types.UserInfoReq) (resp *types.UserInfoResp, err error) {
	// todo: add your logic here and delete this line

	// 从请求上下文中获取用户 ID（JWT 中间件已解析并写入 context）
	uid := ctxdata.GetUId(l.ctx)

	// 调用 User RPC 服务获取用户信息
	userInfoResp, err := l.svcCtx.User.GetUserInfo(l.ctx, &user.GetUserInfoReq{
		Id: uid,
	})
	if err != nil {
		return nil, err
	}

	// 将 RPC 响应转换为 API 响应格式
	var res types.User
	copier.Copy(&res, userInfoResp.User)

	return &types.UserInfoResp{Info: res}, nil
}
