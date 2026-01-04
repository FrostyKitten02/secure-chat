package service

import (
	"context"
	"encoding/base64"
	"secure-chat/dto"
	"secure-chat/repo"
)

func GetDirectMessages(ctx context.Context, user1ID, user2ID string) ([]dto.DirectMessage, error) {
	messages, err := repo.FindDirectMessagesBetweenUsers(ctx, user1ID, user2ID)
	if err != nil {
		return nil, err
	}

	result := make([]dto.DirectMessage, 0, len(messages))
	for _, m := range messages {
		cipherTextBase64 := base64.StdEncoding.EncodeToString(m.CipherText)
		nonceBase64 := base64.StdEncoding.EncodeToString(m.Nonce)

		result = append(result, dto.DirectMessage{
			FromUserId: m.SenderID.String(),
			CipherText: cipherTextBase64,
			Nonce:      nonceBase64,
		})
	}

	return result, nil
}
