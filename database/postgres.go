package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// "github.com/jackc/pgx/v5"

// PostgresDatabase implements Database interface. Connects to
type PostgresDatabase struct {
        *sql.DB
}

// NewPostgresDatabase creates new PostgresDatabase object and connects to PostgreSQL server.
// Please DON'T use whitespaces and backslashes in credentials (it is possible but unwanted).
// NEVER user tailing backslashes.
func NewPostgresDatabase(addr, user, password string) (*PostgresDatabase, error) {
        // Details: 34.1.2 https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING
        connStr := fmt.Sprintf("host='%s' port=5432 user='%s' password='%s'",
                addr, user, password)
        db, err := sql.Open("postgres", connStr)
        if err != nil {
                return nil, err
        }

        err = db.Ping()
        if err != nil {
                return nil, err
        }

        pg := &PostgresDatabase{db}
        return pg, nil
}

func (s *PostgresDatabase) GetNotes(userId int) NoteGetAll {
        var notes NoteGetAll 
        return notes
}

func (s *PostgresDatabase) PostNote(userId int, note NoteCreate) {}
