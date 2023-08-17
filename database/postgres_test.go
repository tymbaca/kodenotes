package database

import (
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

var (
	host     = os.Getenv("POSTGRES_HOST")
	password = os.Getenv("POSTGRES_PASSWORD")
)

func TestMain(m *testing.M) {
	m.Run()
}

// func setup() {

// }

func TestExpiredSession(t *testing.T) {
	db, err := NewPostgresDatabase(host, password)
	if err != nil {
		panic(err)
	}

	db.addUuidExtension()
	db.createUsersTable()
	db.createSessionsTable()

	id := addUserReturnId(db, "Misha", "lalala")
	_, err = db.Exec("INSERT INTO sessions (user_id) VALUES ($1)", id)
	if err != nil {
		t.FailNow()
	}

	if countTable(db, "sessions") != 1 {
		t.FailNow()
	}

	go db.loopCleanExpiredSessions(100*time.Millisecond, 50*time.Millisecond)

	// Check before cleaning
	time.Sleep(51 * time.Millisecond)
	if countTable(db, "sessions") != 1 {
		t.FailNow()
	}

	// Check after cleaning
	time.Sleep(51 * time.Millisecond)
	if countTable(db, "sessions") != 0 {
		t.FailNow()
	}

}

func countTable(db *PostgresDatabase, table string) int {
	var count int

	rows, err := db.Query("SELECT count(*) FROM $1;", table)
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

func addUserReturnId(db *PostgresDatabase, username, password string) string {
	var id string // BAD. ID is not a string
	row, err := db.Query("INSERT INTO users (username, password) VALUES ($1, $1) RETURNING id;", username, password)
	if err != nil {
		panic(err)
	}
	err = row.Scan(&id)
	if err != nil {
		panic(err)
	}

	return id
}
