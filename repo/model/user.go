package model

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
}

type DirectMessage struct {
	ID                 uuid.UUID
	SenderID           uuid.UUID
	ReceiverID         uuid.UUID
	CipherText         []byte
	Nonce              []byte
	SenderIdentityID   uuid.UUID
	ReceiverIdentityID uuid.UUID
	CreatedAt          time.Time
}

type Identity struct {
	ID         uuid.UUID
	PubKey     []byte
	EncPrivKey []byte
	Active     bool
	UserID     uuid.UUID
}
