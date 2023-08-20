package database

import "github.com/google/uuid"

type Database interface {
	RegisterUser(creds UserSecureCredentials) (uuid.UUID, error)
	GetUserId(username string) uuid.NullUUID
	GetUserIdIfAuthorized(creds UserSecureCredentials) uuid.NullUUID
	GetNotes(userId uuid.UUID) (NoteGetAll, error)
	PostNote(userId uuid.UUID, note NoteCreate) error
	MaxCredsLength() int
}
