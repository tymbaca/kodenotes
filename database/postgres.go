package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

var ErrUsernameExists = errors.New("username already exists")

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
                username        VARCHAR(250) NOT NULL UNIQUE,
                password        VARCHAR(250) NOT NULL,

                PRIMARY KEY (id)
        );`

	_, err := d.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

// func (d *PostgresDatabase) CreateUser(username, password string) error {}
func (d *PostgresDatabase) RegisterUser(creds UserSecureCredentials) (uuid.UUID, error) {
        var userId uuid.UUID
	err := d.QueryRow("INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id;", 
                          creds.Username, creds.Password).Scan(&userId)
        if err != nil {
                return uuid.UUID{}, err
        }
        return userId, nil
}

func (d *PostgresDatabase) GetUserId(creds UserSecureCredentials) uuid.NullUUID {
	var result uuid.NullUUID
	err := d.QueryRow("SELECT id FROM users WHERE username = $1 AND password = $2").Scan(&result)
	if err != nil || !result.Valid {
		return uuid.NullUUID{Valid: false}
	} else {
		return result
	}
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

func (d *PostgresDatabase) GetNotes(userId uuid.UUID) (NoteGetAll, error) {
	var result NoteGetAll

	rows, err := d.Query(`SELECT id, user_id, text FROM notes WHERE user_id = $1;`, userId)
	if err != nil {
		return NoteGetAll{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var note NoteGet
		err = rows.Scan(&note.Id,
			&note.UserId,
			&note.Text)
		if err != nil {
			return NoteGetAll{}, err
		}

		result.Notes = append(result.Notes, note)
	}

	return result, nil
}

func (d *PostgresDatabase) PostNote(userId uuid.UUID, note NoteCreate) error {
	_, err := d.Exec(`INSERT INTO notes (user_id, text) VALUES ($1, $2);`,
		userId, note.Text)
	// It also detects if userId not present in table users
	if err != nil {
		return err
	}
	return nil
}
