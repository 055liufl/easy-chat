// =============================================================================
// 建立用户会话业务逻辑
// =============================================================================
// 提供会话建立功能，支持:
//   - 单聊会话建立（双向建立发送者和接收者的会话）
//   - 群聊会话建立（为用户添加群聊会话）
//
// 数据操作:
//   - Conversation 集合：存储会话基本信息
//   - Conversations 集合：存储用户的会话列表
//
// 业务场景:
//   - 用户首次发送消息时建立会话
//   - 用户加入群聊时建立群聊会话
//   - 确保会话的幂等性（重复建立不会出错）
//
// =============================================================================
package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"imooc.com/easy-chat/apps/im/immodels"
	"imooc.com/easy-chat/apps/im/rpc/im"
	"imooc.com/easy-chat/apps/im/rpc/internal/svc"
	"imooc.com/easy-chat/pkg/constants"
	"imooc.com/easy-chat/pkg/wuid"
	"imooc.com/easy-chat/pkg/xerr"
)

// SetUpUserConversationLogic 建立用户会话业务逻辑处理器
type SetUpUserConversationLogic struct {
	ctx    context.Context        // 请求上下文，用于超时控制和链路追踪
	svcCtx *svc.ServiceContext    // 服务上下文，提供数据模型访问
	logx.Logger                   // 日志记录器
}

// NewSetUpUserConversationLogic 创建建立用户会话业务逻辑实例
//
// 参数:
//   - ctx: 请求上下文
//   - svcCtx: 服务上下文，包含数据模型等依赖
//
// 返回:
//   - *SetUpUserConversationLogic: 业务逻辑处理器实例
func NewSetUpUserConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetUpUserConversationLogic {
	return &SetUpUserConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SetUpUserConversation 建立会话: 群聊, 私聊
// 根据聊天类型建立不同的会话关系
//
// 参数:
//   - in: 建立会话请求
//     - SendId: 发送者用户 ID
//     - RecvId: 接收者 ID（单聊为用户 ID，群聊为群组 ID）
//     - ChatType: 聊天类型（1=单聊，2=群聊）
//
// 返回:
//   - *im.SetUpUserConversationResp: 建立会话响应
//   - error: 错误信息
//
// 业务流程:
//   单聊模式:
//     1. 生成会话 ID（组合发送者和接收者 ID）
//     2. 检查会话是否已存在，不存在则创建
//     3. 为发送者和接收者双方建立会话记录
//   群聊模式:
//     1. 使用群组 ID 作为会话 ID
//     2. 为发送者建立群聊会话记录
func (l *SetUpUserConversationLogic) SetUpUserConversation(in *im.SetUpUserConversationReq) (*im.SetUpUserConversationResp, error) {
	// todo: add your logic here and delete this line

	var res im.SetUpUserConversationResp
	switch constants.ChatType(in.ChatType) {
	case constants.SingleChatType:
		// 生成会话的id
		conversationId := wuid.CombineId(in.SendId, in.RecvId)
		// 验证是否建立过会话
		conversationRes, err := l.svcCtx.ConversationModel.FindOne(l.ctx, conversationId)
		if err != nil {
			// 建立会话
			if err == immodels.ErrNotFound {
				err = l.svcCtx.ConversationModel.Insert(l.ctx, &immodels.Conversation{
					ConversationId: conversationId,
					ChatType:       constants.SingleChatType,
				})

				if err != nil {
					return nil, errors.Wrapf(xerr.NewDBErr(), "ConversationsModel.Insert err %v", err)
				}
			} else {
				return nil, errors.Wrapf(xerr.NewDBErr(), "ConversationsModel.FindOne err %v, req %v", err, conversationId)
			}
		} else if conversationRes != nil {
			return &res, nil
		}
		// 建立两者的会话
		err = l.setUpUserConversation(conversationId, in.SendId, in.RecvId, constants.SingleChatType, true)
		if err != nil {
			return nil, err
		}
		err = l.setUpUserConversation(conversationId, in.RecvId, in.SendId, constants.SingleChatType, false)
		if err != nil {
			return nil, err
		}
	case constants.GroupChatType:
		err := l.setUpUserConversation(in.RecvId, in.SendId, in.RecvId, constants.GroupChatType, true)
		if err != nil {
			return nil, err
		}
	}

	return &res, nil
}

// setUpUserConversation 为用户建立会话记录
// 在用户的会话列表中添加新的会话
//
// 参数:
//   - conversationId: 会话 ID
//   - userId: 用户 ID
//   - recvId: 接收者 ID（用于标识会话对象）
//   - chatType: 聊天类型
//   - isShow: 是否在会话列表中显示
//
// 返回:
//   - error: 错误信息
//
// 业务流程:
//   1. 查询用户的会话列表，不存在则创建
//   2. 检查会话是否已在列表中，存在则直接返回
//   3. 添加新会话到用户的会话列表
//   4. 更新用户的会话列表到数据库
func (l *SetUpUserConversationLogic) setUpUserConversation(conversationId, userId, recvId string,
	chatType constants.ChatType, isShow bool) error {
	// 用户的会话列表
	conversations, err := l.svcCtx.ConversationsModel.FindByUserId(l.ctx, userId)
	if err != nil {
		if err == immodels.ErrNotFound {
			conversations = &immodels.Conversations{
				ID:               primitive.NewObjectID(),
				UserId:           userId,
				ConversationList: make(map[string]*immodels.Conversation),
			}
		} else {
			return errors.Wrapf(xerr.NewDBErr(), "ConversationsModel.FindOne err %v, req %v", err, userId)
		}
	}

	// 更新会话记录
	if _, ok := conversations.ConversationList[conversationId]; ok {
		return nil
	}

	// 添加会话记录
	conversations.ConversationList[conversationId] = &immodels.Conversation{
		ConversationId: conversationId,
		ChatType:       constants.SingleChatType,
		IsShow:         isShow,
	}

	// 更新
	err = l.svcCtx.ConversationsModel.Update(l.ctx, conversations)
	if err != nil {
		return errors.Wrapf(xerr.NewDBErr(), "ConversationsModel.Insert err %v, req %v", err, conversations)
	}
	return nil
}
