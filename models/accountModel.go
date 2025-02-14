package model

import (
	"regexp"
	"time"

	"github.com/demkowo/booking/utils/errs"
	"github.com/golang-jwt/jwt"
	uuid "github.com/google/uuid"
)

type Account struct {
	ID       uuid.UUID      `json:"id"`
	Email    string         `json:"email"`
	Password string         `json:"password"`
	Roles    []AccountRoles `json:"roles"`
	APIKeys  []APIKey       `json:"api_keys"`
	Created  time.Time      `json:"created"`
	Updated  time.Time      `json:"updated"`
	Blocked  time.Time      `json:"blocked"`
	Deleted  bool           `json:"deleted"`
	jwt.StandardClaims
}

type AccountRoles struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type UpdateRoles struct {
	Roles []RoleUpdate `json:"roles"`
}

type RoleUpdate struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

type APIKey struct {
	ID        uuid.UUID `json:"id"`
	Key       string    `json:"key"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (a *Account) Validate() *errs.Error {
	if a.Email == "" || a.Password == "" {
		return errs.NewError("Name and Email can't be empty", 400, "Bad Request", nil)
	}

	emailRegex := regexp.MustCompile(`^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$`)
	if !emailRegex.MatchString(a.Email) {
		return errs.NewError("Invalid email address", 400, "Bad Request", nil)
	}

	password := a.Password
	if len(password) < 8 || !containsCapitalLetter(password) || !containsSpecialCharacter(password) || !containsDigit(password) {
		return errs.NewError("Password must contain at least 8 characters, 1 capital letter, 1 special character, and 1 digit", 400, "Bad Request", nil)
	}
	return nil
}

func containsCapitalLetter(password string) bool {
	match, _ := regexp.MatchString("[A-Z]", password)
	return match
}

func containsSpecialCharacter(password string) bool {
	match, _ := regexp.MatchString(`[!@#$%^&*()-_+=\[\]{}|;:'",<>.?/~]`, password)
	return match
}

func containsDigit(password string) bool {
	match, _ := regexp.MatchString("[0-9]", password)
	return match
}
