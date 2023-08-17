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

	return pg, nil
}

func (d *PostgresDatabase) Init() error {
	err := d.addUuidExtension()
	if err != nil {
		return err
	}

	err = d.createUsersTable()
	if err != nil {
		return err
	}

	err = d.createNotesTable()
	if err != nil {
		return err
	}

	return nil
}

func (d *PostgresDatabase) addUuidExtension() error {
	query := `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`

	_, err := d.Exec(query)
	if err != nil {
		return err
	}
	return nil

}

// createUsersTable creates users table in PostgreSQL.
// DEPENDS on d.addUuidExtension().
func (d *PostgresDatabase) createUsersTable() error {
	query := `
        CREATE TABLE IF NOT EXISTS users (
                id              UUID DEFAULT uuid_generate_v4(),
                username        VARCHAR(250) NOT NULL,
                password        VARCHAR(250) NOT NULL,

                PRIMARY KEY (id)
        );`

	_, err := d.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

// createSessionsTable creates users table in PostgreSQL.
// DEPENDS on d.addUuidExtension() and d.createUsersTable().
func (d *PostgresDatabase) createSessionsTable() error {
	query := `
        CREATE TABLE IF NOT EXISTS sessions (
                id              UUID DEFAULT uuid_generate_v1(),
                user_id         UUID NOT NULL,
                last_used_at    TIMESTAMP NOT NULL DEFAULT current_timestamp,

                PRIMARY KEY (id),
                CONSTRAINT fk_session_user FOREIGN KEY (user_id) REFERENCES users(id)
        );`

	_, err := d.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (d *PostgresDatabase) createNotesTable() error {
	query := `
        CREATE TABLE IF NOT EXISTS notes (
                id      UUID DEFAULT uuid_generate_v1(),
                user_id UUID,
                text    TEXT,

                PRIMARY KEY (id),
                CONSTRAINT fk_note_user FOREIGN KEY (user_id) REFERENCES "users"(id)
        );`

	_, err := d.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (d *PostgresDatabase) GetNotes(userId uuid.UUID) NoteGetAll {
	d.Query(`SELECT * FROM notes`)
	var notes NoteGetAll
	return notes
}

func (d *PostgresDatabase) PostNote(userId uuid.UUID, note NoteCreate) {}
