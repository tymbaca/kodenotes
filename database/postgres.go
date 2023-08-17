package database

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// "github.com/jackc/pgx/v5"

// PostgresDatabase implements Database interface. Connects to
type PostgresDatabase struct {
	*sql.DB
}

// NewPostgresDatabase creates new PostgresDatabase object and connects to PostgreSQL server.
// DON'T use whitespaces and backslashes in credentials.
func NewPostgresDatabase(host, password string) (*PostgresDatabase, error) {
	// Details: 34.1.2 https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING
	connStr := fmt.Sprintf("host=%s dbname=postgres user=postgres password=%s sslmode=disable", host, password)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	pg := &PostgresDatabase{db}

	err = pg.init()
	if err != nil {
		return nil, err
	}

	return pg, nil
}

func (s *PostgresDatabase) init() error {
        err := s.addUuidExtension()
	if err != nil {
		return err
	}

        err = s.createUsersTable()
        if err != nil {
                return err
        }

	err = s.createNotesTable()
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresDatabase) addUuidExtension() error {
	query := `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`

	_, err := s.Exec(query)
	if err != nil {
		return err
	}
	return nil
        
}

func (s *PostgresDatabase) createUsersTable() error {
	query := `
        CREATE TABLE IF NOT EXISTS users (
                id              UUID DEFAULT uuid_generate_v1(),
                username        VARCHAR(250) NOT NULL,
                password        VARCHAR(250) NOT NULL,

                PRIMARY KEY (id)
        );`

	_, err := s.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresDatabase) createNotesTable() error {
	query := `
        CREATE TABLE IF NOT EXISTS notes (
                id      UUID DEFAULT uuid_generate_v1(),
                user_id UUID,
                text    TEXT,

                PRIMARY KEY (id),
                CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES "users"(id)
        );`

	_, err := s.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresDatabase) GetNotes(userId uuid.UUID) NoteGetAll {
        s.Query(`SELECT * FROM notes`)
	var notes NoteGetAll
	return notes
}

func (s *PostgresDatabase) PostNote(userId uuid.UUID, note NoteCreate) {}
