package database

type Database interface {
        GetNotes(userId int) NoteGetAll
        PostNote(userId int, note NoteCreate)
}
