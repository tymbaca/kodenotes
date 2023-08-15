package storage

type Database interface {
        GetNotes() []Note
        PostNote(Note)
}
