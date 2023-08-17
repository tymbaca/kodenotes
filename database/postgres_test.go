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
