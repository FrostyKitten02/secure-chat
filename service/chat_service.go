package service

import (
	"context"
	"github.com/google/uuid"
	"secure-chat/dto"
	"secure-chat/repo"
)

func GetChatsForUser(ctx context.Context, userID uuid.UUID) ([]dto.ChatDto, error) {
	chats, err := repo.FindChatsByUserID(ctx, userID.String())
	if err != nil {
		return nil, err
	}

	res := make([]dto.ChatDto, 0, len(chats))
	for _, c := range chats {
		receiverId := c.User1ID
		if receiverId == userID {
			receiverId = c.User2ID
		}

		recIdentity, identityErr := repo.FindActiveIdentityForUser(ctx, receiverId.String())
		if identityErr != nil {
			return nil, identityErr
		}
		res = append(res, dto.ChatDto{
			ID: c.ID.String(),
			User: dto.ChatUserDto{
				UserId:     receiverId.String(),
				IdentityId: recIdentity.ID,
				PubKey:     recIdentity.PubKey,
			},
		})
	}

	return res, nil
}
