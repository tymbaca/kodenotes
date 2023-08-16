package database

import "github.com/google/uuid"

type Database interface {
        GetNotes(userId uuid.UUID) NoteGetAll
        PostNote(userId uuid.UUID, note NoteCreate)
}
