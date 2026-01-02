package service

import (
	"context"
	"secure-chat/dto"
	"secure-chat/repo"
)

// Search users by query, exclude users already in my chats
func SearchUsersExcludingChats(ctx context.Context, myID string, query string) ([]dto.UserListItemDto, error) {
	users, err := repo.SearchUsersByUsernameExcludingChats(ctx, myID, query)
	if err != nil {
		return nil, err
	}

	res := make([]dto.UserListItemDto, 0, len(users))
	for _, u := range users {
		res = append(res, dto.UserListItemDto{
			ID:       u.ID.String(),
			Username: u.Username,
			Email:    u.Email,
		})
	}
	return res, nil
}
