package api

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"log/slog"
	"secure-chat/dto"
	"secure-chat/env"
	"secure-chat/repo"
	"secure-chat/service"
)

func chatController(api huma.API, options *env.Options) {
	basePath := "/chats"

	huma.Get(api, basePath, func(ctx context.Context, _ *struct{}) (*dto.GetChatsResponse, error) {
		userID := service.GetUserID(ctx)
		if userID == nil {
			return nil, huma.Error401Unauthorized("Unauthorized")
		}

		chats, err := service.GetChatsForUser(ctx, *userID)
		if err != nil {
			return nil, err
		}

		return &dto.GetChatsResponse{
			Body: dto.GetChatsResponseBody{
				Chats: chats,
			},
		}, nil
	})

	huma.Get(api, basePath+"/history/by-user/{userId}", func(ctx context.Context, i *dto.GetChatHistoryRequest) (*dto.GetChatHistoryResponse, error) {
		userID := service.GetUserID(ctx)
		if userID == nil {
			return nil, huma.Error401Unauthorized("Unauthorized")
		}

		//this should never happen
		if i.UserId == "" {
			return nil, huma.Error400BadRequest("UserId is required")
		}

		fromUserId, fromUserIdErr := uuid.Parse(i.UserId)
		if fromUserIdErr != nil {
			return nil, huma.Error400BadRequest("UserId must be a valid UUID")
		}

		msgs, msgsErr := service.GetDirectMessages(context.Background(), userID.String(), fromUserId.String())
		if msgsErr != nil {
			slog.Error("Error getting direct messages", "error", msgsErr)
			return nil, huma.Error500InternalServerError("Internal Server Error")
		}

		fromUser, fromUserErr := repo.FindUserByID(context.Background(), fromUserId.String())
		if fromUserErr != nil {
			return nil, huma.Error500InternalServerError("Internal Server Error")
		}

		return &dto.GetChatHistoryResponse{
			Body: dto.GetChatHistoryResponseBody{
				DirectMessages: msgs,
				FromUser: []dto.ChatUserDto{
					{
						UserId:   fromUser.ID.String(),
						Username: fromUser.Username,
					},
				},
			},
		}, nil
	})
}
