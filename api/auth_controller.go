package api

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"secure-chat/dto"
	"secure-chat/env"
	"secure-chat/service"
)

func authController(api huma.API, options *env.Options) {
	basePath := "/auth"

	huma.Post(api, basePath+"/register", func(ctx context.Context, i *dto.RegisterRequest) (*dto.RegisterResponse, error) {
		err := service.RegisterNewUser(i)
		if err != nil {
			return nil, err
		}

		return &dto.RegisterResponse{
			Status: http.StatusCreated,
		}, nil
	})

	huma.Post(api, basePath+"/login", func(ctx context.Context, i *dto.LoginRequest) (*dto.LoginResponse, error) {
		return service.LoginUser(i)
	})

	//TODO: add session refresh endpoint!!!!
}
