// =============================================================================
// 获取消息已读记录业务逻辑
// =============================================================================
// 实现消息已读/未读状态查询的核心业务逻辑
//
// 功能说明:
//   - 查询指定消息的已读和未读用户列表
//   - 支持单聊和群聊两种场景
//   - 单聊: 判断接收者是否已读
//   - 群聊: 通过位图（bitmap）判断每个群成员的已读状态
//
// 业务流程:
//  1. 根据消息 ID 查询聊天记录
//  2. 根据聊天类型（单聊/群聊）处理已读状态
//  3. 单聊: 检查 ReadRecords 字段判断是否已读
//  4. 群聊: 获取群成员列表，通过位图判断每个成员的已读状态
//  5. 查询用户信息，将用户 ID 转换为手机号
//  6. 返回已读和未读用户的手机号列表
//
// =============================================================================
package logic

import (
	"context"
	"imooc.com/easy-chat/apps/im/rpc/im"
	"imooc.com/easy-chat/apps/social/rpc/socialclient"
	"imooc.com/easy-chat/apps/user/rpc/user"
	"imooc.com/easy-chat/pkg/bitmap"
	"imooc.com/easy-chat/pkg/constants"

	"imooc.com/easy-chat/apps/im/api/internal/svc"
	"imooc.com/easy-chat/apps/im/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// GetChatLogReadRecordsLogic 获取消息已读记录业务逻辑结构
type GetChatLogReadRecordsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetChatLogReadRecordsLogic 创建获取消息已读记录业务逻辑实例
//
// 参数:
//   - ctx: 上下文对象，用于传递请求信息和控制超时
//   - svcCtx: 服务上下文，包含所有依赖项（RPC 客户端等）
//
// 返回:
//   - *GetChatLogReadRecordsLogic: 业务逻辑实例
func NewGetChatLogReadRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatLogReadRecordsLogic {
	return &GetChatLogReadRecordsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetChatLogReadRecords 获取消息已读记录
// 查询指定消息的已读和未读用户列表
//
// 参数:
//   - req: 请求参数，包含消息 ID
//
// 返回:
//   - resp: 响应数据，包含已读和未读用户的手机号列表
//   - err: 错误信息
//
// 业务逻辑:
//  1. 调用 IM RPC 服务查询消息详情
//  2. 根据聊天类型处理已读状态:
//     - 单聊: 检查 ReadRecords 字段，判断接收者是否已读
//     - 群聊: 获取群成员列表，通过位图判断每个成员的已读状态
//  3. 调用 User RPC 服务查询用户信息
//  4. 将用户 ID 转换为手机号返回
func (l *GetChatLogReadRecordsLogic) GetChatLogReadRecords(req *types.GetChatLogReadRecordsReq) (resp *types.GetChatLogReadRecordsResp, err error) {
	// 步骤 1: 调用 IM RPC 服务查询消息详情
	chatlogs, err := l.svcCtx.Im.GetChatLog(l.ctx, &im.GetChatLogReq{
		MsgId: req.MsgId,
	})

	if err != nil || len(chatlogs.List) == 0 {
		return nil, err
	}

	var (
		chatlog = chatlogs.List[0]
		reads   = []string{chatlog.SendId} // 发送者默认已读
		unreads []string
		ids     []string // 所有相关用户 ID
	)

	// 步骤 2: 根据聊天类型处理已读状态
	switch constants.ChatType(chatlog.ChatType) {
	case constants.SingleChatType:
		// 单聊场景: 只需判断接收者是否已读
		// ReadRecords 为空或第一个元素为 0 表示未读
		if len(chatlog.ReadRecords) == 0 || chatlog.ReadRecords[0] == 0 {
			unreads = []string{chatlog.RecvId}
		} else {
			reads = append(reads, chatlog.RecvId)
		}
		ids = []string{chatlog.RecvId, chatlog.SendId}

	case constants.GroupChatType:
		// 群聊场景: 需要判断每个群成员的已读状态
		// 步骤 2.1: 获取群组成员列表
		groupUsers, err := l.svcCtx.Social.GroupUsers(l.ctx, &socialclient.GroupUsersReq{
			GroupId: chatlog.RecvId,
		})
		if err != nil {
			return nil, err
		}

		// 步骤 2.2: 加载已读记录位图
		// 位图中每个用户 ID 对应一个位，1 表示已读，0 表示未读
		bitmaps := bitmap.Load(chatlog.ReadRecords)
		for _, members := range groupUsers.List {
			ids = append(ids, members.UserId)

			// 发送者默认已读，跳过
			if members.UserId == chatlog.SendId {
				continue
			}

			// 步骤 2.3: 检查位图判断是否已读
			if bitmaps.IsSet(members.UserId) {
				reads = append(reads, members.UserId)
			} else {
				unreads = append(unreads, members.UserId)
			}
		}
	}

	// 步骤 3: 查询用户信息
	userEntitys, err := l.svcCtx.User.FindUser(l.ctx, &user.FindUserReq{
		Ids: ids,
	})
	if err != nil {
		return nil, err
	}

	// 步骤 4: 构建用户 ID 到用户实体的映射
	userEntitySet := make(map[string]*user.UserEntity, len(userEntitys.User))
	for i, entity := range userEntitys.User {
		userEntitySet[entity.Id] = userEntitys.User[i]
	}

	// 步骤 5: 将用户 ID 转换为手机号
	// 已读用户列表
	for i, read := range reads {
		if u := userEntitySet[read]; u != nil {
			reads[i] = u.Phone
		}
	}
	// 未读用户列表
	for i, unread := range unreads {
		if u := userEntitySet[unread]; u != nil {
			unreads[i] = u.Phone
		}
	}

	// 步骤 6: 返回结果
	return &types.GetChatLogReadRecordsResp{
		Reads:   reads,
		UnReads: unreads,
	}, nil
}
