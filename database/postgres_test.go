package database

import (
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

var (
	// host     = os.Getenv("POSTGRES_HOST")
	// password = os.Getenv("POSTGRES_PASSWORD")
	host     = "localhost" // WARNING: not production code
	password = "mypassword"
)

func TestMain(m *testing.M) {
	clearDb()
	m.Run()
}

// func setup() {

// }

func clearDb() {
	db, err := NewPostgresDatabase(host, password)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
		DROP SCHEMA public CASCADE;
		CREATE SCHEMA public;
	`)
	if err != nil {
		panic(err)
	}
}

func TestGetNotes(t *testing.T) {
	clearDb()
	// FOREPLAY
	db, err := NewPostgresDatabase(host, password)
	if err != nil {
		panic(err)
	}
	err = db.Init()
	if err != nil {
		panic(err)
	}

	if mustCountTable(db, "users") != 0 || mustCountTable(db, "notes") != 0 {
		t.FailNow()
	}
	userId := mustAddUserReturnId(db, "Semen", "hashedpasswd")

	if mustCountTable(db, "users") != 1 {
		t.FailNow()
	}

	// ACTUAL TEST
	var note NoteCreate
	note.Text = `Hello! This is my text. My name is Semen and only 
				I can see this tho.`

	err = db.PostNote(userId, note)
	if err != nil {
		t.FailNow()
	}
	// Must be increase to 1
	if mustCountTable(db, "notes") != 1 {
		t.FailNow()
	}

	var checkText string
	err = db.QueryRow("SELECT text FROM notes WHERE user_id = $1", userId).Scan(&checkText)
	if err != nil || checkText != note.Text {
		t.FailNow()
	}

	// Trying with random userId
	userIdNotExist := uuid.New()
	note.Text = `Hello! This is my text. I don't know my name and I even don't exist lol.`

	err = db.PostNote(userIdNotExist, note)
	if err != nil {
		// OK
	} else {
		t.FailNow()
	}
	// Must still the same 1 Semen's note
	if mustCountTable(db, "notes") != 1 {
		t.FailNow()
	}
}

func TestGetAllNotes(t *testing.T) {
	clearDb()
	// FOREPLAY
	db, err := NewPostgresDatabase(host, password)
	if err != nil {
		panic(err)
	}
	err = db.Init()
	if err != nil {
		panic(err)
	}

	if mustCountTable(db, "users") != 0 || mustCountTable(db, "notes") != 0 {
		t.FailNow()
	}
	semenId := mustAddUserReturnId(db, "Semen", "hashedpasswd")
	// Added 3 notes
	mustAddNoteReturnId(db, semenId, NoteCreate{Text: "First note"})
	mustAddNoteReturnId(db, semenId, NoteCreate{Text: "Second note"})
	mustAddNoteReturnId(db, semenId, NoteCreate{Text: "Third note"})

	result, err := db.GetNotes(semenId)
	if err != nil || len(result.Notes) != 3 {
		t.FailNow()
	}

	nicolId := mustAddUserReturnId(db, "Nikol", "hashedpasswd")
	// Only 2 notes
	mustAddNoteReturnId(db, nicolId, NoteCreate{Text: "First note"})
	mustAddNoteReturnId(db, nicolId, NoteCreate{Text: "Second note"})

	result, err = db.GetNotes(nicolId)
	if err != nil || len(result.Notes) != 2 {
		t.FailNow()
	}

}

// func TestExpiredSession(t *testing.T) {
// 	db, err := NewPostgresDatabase(host, password)
// 	if err != nil {
// 		panic(err)
// 	}
// 	err = db.addUuidExtension()
// 	if err != nil {
// 		panic(err)
// 	}
// 	err = db.createUsersTable()
// 	if err != nil {
// 		panic(err)
// 	}
// 	err = db.createSessionsTable()
// 	if err != nil {
// 		panic(err)
// 	}

// 	id := mustAddUserReturnId(db, "Misha", "lalala")
// 	_, err = db.Exec("INSERT INTO sessions (user_id) VALUES ($1)", id)
// 	if err != nil {
// 		t.FailNow()
// 	}

// 	if mustCountTable(db, "sessions") != 1 {
// 		t.FailNow()
// 	}

// 	go db.loopCleanExpiredSessions(50000*time.Millisecond, 50*time.Millisecond)

// 	// Check before cleaning
// 	time.Sleep(51 * time.Millisecond)
// 	if mustCountTable(db, "sessions") != 1 {
// 		t.FailNow()
// 	}
// 	t.Log("Session in db")

// 	// Check after cleaning
// 	time.Sleep(51 * time.Millisecond)
// 	if mustCountTable(db, "sessions") != 0 {
// 		t.FailNow()
// 	}
// 	t.Log("Session NOT in db")

// }

func mustCountTable(db *PostgresDatabase, table string) int {
	var count int

	rows, err := db.Query("SELECT count(*) FROM " + table + ";")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			panic(err)
		}
	}
	return count
}

func mustAddUserReturnId(db *PostgresDatabase, username, password string) uuid.UUID {
	var id uuid.UUID // BAD. ID is not a string
	rows, err := db.Query("INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id;", username, password)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			panic(err)
		}
	}

	return id
}

func mustAddNoteReturnId(db *PostgresDatabase, user_id uuid.UUID, note NoteCreate) uuid.UUID {
	var id uuid.UUID
	err := db.QueryRow("INSERT INTO notes (user_id, text) VALUES ($1, $2) RETURNING id", user_id, note.Text).Scan(&id)
	if err != nil {
		panic(err)
	}
	return id
}
