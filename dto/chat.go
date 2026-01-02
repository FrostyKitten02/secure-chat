package dto

import "github.com/google/uuid"

type ChatDto struct {
	ID   string      `json:"id"`
	User ChatUserDto `json:"user"`
}

type ChatUserDto struct {
	UserId     string    `json:"userId"`
	Username   string    `json:"username"`
	IdentityId uuid.UUID `json:"identityId"`
	PubKey     string    `json:"pubKey"`
}

type GetChatsResponse struct {
	Body GetChatsResponseBody
}

type GetChatsResponseBody struct {
	Chats []ChatDto `json:"chats"`
}
