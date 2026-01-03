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

type Identity struct {
	ID         uuid.UUID
	PubKey     []byte
	EncPrivKey []byte
	Active     bool
	UserID     uuid.UUID
}
