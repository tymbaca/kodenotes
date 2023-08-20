package database

import (
	"crypto/sha256"
	"encoding/base64"

	"github.com/google/uuid"
)

type UserSecureCredentials struct {
	Username string
	Password string
}

func NewUserSecureCredentials(username, rawPassword string) UserSecureCredentials {
	password := hashPassword(rawPassword)
	result := UserSecureCredentials{Username: username, Password: password}
	return result
}

func hashPassword(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	sum := h.Sum(nil)
	sha := base64.URLEncoding.EncodeToString(sum)
	return sha
}

type User struct {
	Id       uuid.UUID
	Username string
	Password string
}

type NoteModel struct {
	Id     uuid.UUID
	UserId uuid.UUID
	Text   string
}

type NoteGetAll struct {
	Notes []NoteGet `json:"notes"`
}

type NoteGet struct {
	Id     uuid.UUID `json:"id"`
	UserId string `json:"userId"`
	Text   string `json:"text"`
}

type NoteCreate struct {
	Text string `json:"text"`
}
