package database

import (
        "os"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

var (
	host     = os.Getenv("POSTGRES_HOST")
	password = os.Getenv("POSTGRES_PASSWORD")
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

func TestReqisterUser(t *testing.T) {
        db := mustSetupCleanDb()

        semenId, err := db.RegisterUser(UserSecureCredentials{Username: "Semen", Password: "hashedpasswd"})
        if err != nil || mustCountTable(db, "users") != 1 {
                t.FailNow()
        }
        var semen User
        err = db.QueryRow("SELECT id, username, password FROM users WHERE id = $1", semenId).Scan(&semen.Id, &semen.Username, &semen.Password)
        if err != nil {
                t.FailNow()
        }

        // Adding user with same username. Must return error
        _, err = db.RegisterUser(UserSecureCredentials{Username: "Semen", Password: "anotherpasswd"})
        if err != nil {
                // OK
        } else {
                t.FailNow()
        }

        // Adding user with too long credentials. 382 charackers (Max is set to 250)
        longUsername := "Very long username gdisopgj 90sdguj0er89osgj r8iosgoerhnjgo34ijg80434w4j 890guerjgj 023ugjf 0298nub 07ut08e4u 027ut 0uje4029ytut0 n94eg7uj480gyu4 n3y2u0-y9ghui89p-e48y-n092y8i p-e4r92ui890ghn86i42 0vn68i4-2098bh49-2nw86b 8i0-246i696 666=2640-96ib-2486 0-24i9-n=b9 ih046i90oik6ygh03yyyyiykyyyy0py[y3y54iy3kmyp0i-34[yikh-08jm4yb,gm3 05ik6340=f69g3k0=g-6tm90-4r39646b846ryh4903y8i9034u"
        longPassword := "Very long password 02934tu0eu 89034t7gu89n74r89ftgyu894eyuhnvf9y5t894bb7y589 vbyu79ervgheu8iogyeiruhg8ier7gyh3794hiurh3g8943iouhjf58uio34yht89347hvtn98347ngbv89t34 nu9687340b69g 7u039476083gn739047h3h94gf6890347h9fg734906g7934867h3480g76h9083745h69803gb6h89347890h3g7n0347868034g76034hb6730g760376h90234g0v4n0673904g7b6024760247034670934786gbv84672390g73b034706g7b80347603760hb3g704"
        longUserCreds := UserSecureCredentials{Username: longUsername, Password: longPassword}
        _, err = db.RegisterUser(longUserCreds)
        if err != nil {
                // OK
        } else {
                t.FailNow()
        }
}

func TestGetNotes(t *testing.T) {
        db := mustSetupCleanDb()

	userId := mustAddUserReturnId(db, "Semen", "hashedpasswd")
	if mustCountTable(db, "users") != 1 {
		t.FailNow()
	}

	// ACTUAL TEST
	var note NoteCreate
	note.Text = `Hello! This is my text. My name is Semen and only I can see this tho.`

        err := db.PostNote(userId, note)
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
        db := mustSetupCleanDb()

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

func mustSetupCleanDb() *PostgresDatabase {
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
		panic("tables are not clean")
	}
        return db
}

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
