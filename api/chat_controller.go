package api

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"secure-chat/dto"
	"secure-chat/env"
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
}
