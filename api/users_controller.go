package api

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"secure-chat/dto"
	"secure-chat/env"
	"secure-chat/service"
)

func usersController(api huma.API, options *env.Options) {
	basePath := "/users"

	huma.Get(api, basePath+"/list", func(ctx context.Context, params *struct {
		Query string `query:"query"`
	}) (*dto.UserListResponse, error) {
		userID := service.GetUserID(ctx)
		if userID == nil {
			return nil, huma.Error401Unauthorized("Unauthorized")
		}

		users, err := service.SearchUsersExcludingChats(ctx, *userID, params.Query)
		if err != nil {
			return nil, err
		}

		return &dto.UserListResponse{
			Body: dto.UserListResponseBody{
				Users: users,
			},
		}, nil
	})

	huma.Post(api, basePath+"{userId}/add-to-chat", func(ctx context.Context, req *dto.AddUserToChatRequest) (*dto.AddUserToChatResponse, error) {
		userID := service.GetUserID(ctx)
		if userID == nil {
			return nil, huma.Error401Unauthorized("Unauthorized")
		}

		err := service.CreateChatWithUser(ctx, userID.String(), req.UserID)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})
}
