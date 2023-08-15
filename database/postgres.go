package database

import (
        // "github.com/jackc/pgx/v5"
)

// PostgresDatabase implements Database interface. Connects to 
type PostgresDatabase struct {
        addr            string
        user            string
        password        string
}

func NewPostgresDatabase(addr, user, password string) *PostgresDatabase {
        db := &PostgresDatabase{addr: addr, user: user, password: password}
        // Maybe need to add ping and return error if db is down
        return db
}

func (s *PostgresDatabase) GetNotes(userId int) NoteGetAll {
        var notes NoteGetAll 
        return notes
}

func (s *PostgresDatabase) PostNote(userId int, note NoteCreate) {}
