// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package database

import (
	"time"

	"github.com/google/uuid"
)

type Item struct {
	ID    uuid.UUID
	Name  string
	Price int32
}

type Transaction struct {
	ID         uuid.UUID
	SenderID   uuid.UUID
	ReceiverID uuid.UUID
	Amount     int32
	CreatedAt  time.Time
}

type User struct {
	ID             uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Username       string
	HashedPassword string
	Balance        int32
}

type UsersItem struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	ItemID   uuid.UUID
	Quantity int32
}
