package database

import "github.com/google/uuid"

type Database interface {
        GetUserId(creds UserSecureCredentials) uuid.NullUUID
        GetNotes(userId uuid.UUID) (NoteGetAll, error)
        PostNote(userId uuid.UUID, note NoteCreate) error
}
