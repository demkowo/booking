package model

import (
	"time"

	uuid "github.com/google/uuid"
)

const (
	AVAILABLE    = 0
	BLOCKED      = 1
	BOOK_REQUEST = 2
	RESERVATION  = 3
	RENT         = 4
)

type Reservation struct {
	Id        uuid.UUID
	UserId    uuid.UUID
	RoomID    uuid.UUID
	Status    int
	StartDate time.Time
	EndDate   time.Time
	Created   time.Time
	Updated   time.Time
	Deleted   bool
}
