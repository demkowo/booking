package model

import (
	"time"

	"github.com/google/uuid"
)

type Room struct {
	Id      uuid.UUID
	Name    string
	Created time.Time
	Updated time.Time
}
