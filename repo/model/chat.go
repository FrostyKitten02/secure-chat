package model

import "github.com/google/uuid"

type Chat struct {
	ID      uuid.UUID
	User1ID uuid.UUID
	User2ID uuid.UUID
}
