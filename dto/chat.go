package dto

import "github.com/google/uuid"

type ChatDto struct {
	ID   string      `json:"id"`
	User ChatUserDto `json:"user"`
}

// TODO: in future add identity used for text, rn it is optimistic using only current identity
type DirectMessage struct {
	FromUserId string `json:"fromUserId"`
	CipherText string `json:"cipherText"`
	Nonce      string `json:"nonce"`
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

type GetChatHistoryRequest struct {
	UserId string `path:"userId"`
}

type GetChatHistoryResponse struct {
	Body GetChatHistoryResponseBody
}

type GetChatHistoryResponseBody struct {
	DirectMessages []DirectMessage `json:"directMessages"`
	FromUser       []ChatUserDto   `json:"fromUser"`
}
