package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/tymbaca/kodenotes/log"
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
		log.Error("PG: Cannot establish connection with PostgreSQL, addr: %s, user: postgres, error %s",
			host, err.Error())
		return nil, err
	}
	log.Info("PG: Opened PostgreSQL connection, addr: %s, user: %s, ssl: off", host, "postgres")

	err = db.Ping()
	if err != nil {
		log.Error("PG: Cannot establish connection with PostgreSQL, addr: %s, user: postgres, error: %s",
			host, err.Error())
		return nil, err
	}
	log.Info("PG: PostgreSQL connection is healthy")
	pg := &PostgresDatabase{db}

	return pg, nil
}

func (d *PostgresDatabase) Init() error {
	err := d.addUuidExtension()
	if err != nil {
		log.Error("PG: Error while 'uuid-ossp' extension: %s", err.Error())
		return err
	}
	log.Info("PG: Added 'uuid-ossp' extension")

	err = d.createUsersTable()
	if err != nil {
		log.Error("PG: error while creating 'users' table, error: %s", err.Error())
		return err
	}
	log.Info("PG: created 'users' table")

	err = d.createNotesTable()
	if err != nil {
		log.Error("PG: error while creating 'notes' table, error: %s", err.Error())
		return err
	}
	log.Info("PG: created 'notes' table")

	return nil
}

func (d *PostgresDatabase) addUuidExtension() error {
	_, err := d.Exec(`DROP EXTENSION IF EXISTS "uuid-ossp";`)
	if err != nil {
		return err
	}
	_, err = d.Exec(`CREATE EXTENSION "uuid-ossp";`)
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

func (d *PostgresDatabase) MaxCredsLength() int {
	return 250
}

// createNotesTable creates notes table in PostgreSQL.
// DEPENDS on d.addUuidExtension and d.createUsersTable.
func (d *PostgresDatabase) createNotesTable() error {
	query := `
        CREATE TABLE IF NOT EXISTS notes (
                id      UUID DEFAULT uuid_generate_v4(),
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

// func (d *PostgresDatabase) CreateUser(username, password string) error {}
func (d *PostgresDatabase) RegisterUser(creds UserSecureCredentials) (uuid.UUID, error) {
	var userId uuid.UUID
	err := d.QueryRow("INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id;",
		creds.Username, creds.Password).Scan(&userId)
	if err != nil {
		log.Warn("PG: Unsuccessful attemt to insert user '%s' in database", creds.Username)
		return uuid.UUID{}, err
	}
	log.Info("PG: inserted new user '%s' in database", creds.Username)
	return userId, nil
}

func (d *PostgresDatabase) GetUserId(username string) uuid.NullUUID {
	var result uuid.NullUUID
	err := d.QueryRow("SELECT id FROM users WHERE username = $1", username).Scan(&result)
	if err != nil || !result.Valid {
		log.Info("PG: user '%s' not found in database", username)
		return uuid.NullUUID{Valid: false}
	} else {
		log.Info("PG: found user '%s' in database", username)
		return result
	}
}

func (d *PostgresDatabase) GetUserIdIfAuthorized(creds UserSecureCredentials) uuid.NullUUID {
	var result uuid.NullUUID
	err := d.QueryRow("SELECT id FROM users WHERE username = $1 AND password = $2", creds.Username, creds.Password).Scan(&result)
	if err != nil || !result.Valid {
		log.Info("PG: cannot authorize user '%s'", creds.Username)
		return uuid.NullUUID{Valid: false}
	} else {
		log.Info("PG: user '%s' is authorized", creds.Username)
		return result
	}
}

func (d *PostgresDatabase) GetNotes(userId uuid.UUID) (NoteGetAll, error) {
	var result NoteGetAll

	rows, err := d.Query(`SELECT id, user_id, text FROM notes WHERE user_id = $1;`, userId)
	if err != nil {
		log.Error("PG: cannot get notes from database for user id: '%s', error: %s", userId.String(), err.Error())
		return NoteGetAll{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var note NoteGet
		err = rows.Scan(&note.Id,
			&note.UserId,
			&note.Text)
		if err != nil {
			log.Error("PG: cannot parse notes from database for user id: '%s', error: %s", userId.String(), err.Error())
			return NoteGetAll{}, err
		}

		result.Notes = append(result.Notes, note)
	}
	log.Info("PG: getted notes from database for user id: '%s'", userId.String())
	return result, nil
}

func (d *PostgresDatabase) PostNote(userId uuid.UUID, note NoteCreate) error {
	_, err := d.Exec(`INSERT INTO notes (user_id, text) VALUES ($1, $2);`,
		userId, note.Text)
	// It also detects if userId not present in table users
	if err != nil {
		log.Error("PG: cannot insert note for user id: '%s', error: %s", userId.String(), err.Error())
		return err
	}
	log.Info("PG: inserted note for user id: '%s'", userId.String())
	return nil
}
