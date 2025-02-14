package model

import (
	"time"

	uuid "github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID
	FirstName string
	LastName  string
	Phone     string
	Email     string
	Country   string
	Created   time.Time
	Updated   time.Time
	Deleted   bool
}
