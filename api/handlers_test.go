package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/uuid"
	"github.com/tymbaca/kodenotes/database"
	"github.com/tymbaca/kodenotes/spellcheck"
	"github.com/tymbaca/kodenotes/util"
)

const (
	serverPortEnvVar = "SERVER_PORT"
	pgHostEnvVar     = "POSTGRES_HOST"
	pgPasswordEnvVar = "POSTGRES_PASSWORD"

	usersTable = "users"
	notesTable = "notes"
)

var (
	server, db  = mustSetupServerAndDb()
	defUsername = "username"
	defPassword = "password"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestRegister(t *testing.T) {
	mustClearDb(db)

	// Register
	rr := testRegister(defUsername, defPassword, http.MethodPost)
	if rr.Code != http.StatusCreated {
		t.Error()
	}
	if mustCountTable(usersTable) != 1 {
		t.Error()
	}

	// Try register same username
	rr = testRegister(defUsername, "anotherPassword", http.MethodPost)
	if rr.Code != http.StatusUnprocessableEntity {
		t.Error()
	}
	if mustCountTable(usersTable) != 1 {
		t.Error()
	}

	userId := db.GetUserIdIfAuthorized(database.NewUserSecureCredentials(defUsername, defPassword))
	if !userId.Valid {
		t.Error()
	}

	// Adding user with too long credentials. 382 charackers (Max is set to 250)
	longUsername := "Very long username gdisopgj 90sdguj0er89osgj r8iosgoerhnjgo34ijg80434w4j 890guerjgj 023ugjf 0298nub 07ut08e4u 027ut 0uje4029ytut0 n94eg7uj480gyu4 n3y2u0-y9ghui89p-e48y-n092y8i p-e4r92ui890ghn86i42 0vn68i4-2098bh49-2nw86b 8i0-246i696 666=2640-96ib-2486 0-24i9-n=b9 ih046i90oik6ygh03yyyyiykyyyy0py[y3y54iy3kmyp0i-34[yikh-08jm4yb,gm3 05ik6340=f69g3k0=g-6tm90-4r39646b846ryh4903y8i9034u"
	longPassword := "Very long password 02934tu0eu 89034t7gu89n74r89ftgyu894eyuhnvf9y5t894bb7y589 vbyu79ervgheu8iogyeiruhg8ier7gyh3794hiurh3g8943iouhjf58uio34yht89347hvtn98347ngbv89t34 nu9687340b69g 7u039476083gn739047h3h94gf6890347h9fg734906g7934867h3480g76h9083745h69803gb6h89347890h3g7n0347868034g76034hb6730g760376h90234g0v4n0673904g7b6024760247034670934786gbv84672390g73b034706g7b80347603760hb3g704"
	rr = testRegister(longUsername, longPassword, http.MethodPost)
	if rr.Code != http.StatusRequestEntityTooLarge {
		t.Error()
	}

	rr = testRegister("newname", "newpasswd", http.MethodGet)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Error()
	}
}

func TestPostNote(t *testing.T) {
	mustClearDb(db)
	text := `This is my text`
	textRus := `Этом мой текст`
	bigTest := url.QueryEscape(`Liberty is the state of being free within society from oppressive restrictions imposed by authority on one’s way of life, behavior, or political views.[1]

	In theology, liberty is freedom from the effects of "sin, spiritual servitude, [or] worldly ties".[2]
	
	In economics, liberty means free, fair, and open competition, often referred to as a free market.
	
	Sometimes liberty is differentiated from freedom by using the word "freedom" primarily, if not exclusively, to mean the ability to do as one wills and what one has the power to do; and using the word "liberty" to mean the absence of arbitrary restraints, taking into account the rights of all involved. In this sense, the exercise of liberty is subject to capability and limited by the rights of others. Thus liberty entails the responsible use of freedom under the rule of law without depriving anyone else of their freedom. Liberty can be taken away as a form of punishment. In many countries, people can be deprived of their liberty if they are convicted of criminal acts.
        `)

	bigTextRus := url.QueryEscape(`Существует множество различных определений свободы[⇨]. В этике понимание свободы связано с наличием свободы воли человека.

        Свобода в философии — универсалия культуры субъектного ряда, фиксирующая возможность деятельности и поведения в условиях отсутствия внешнего целеполагания[3].

        Свобода личности в праве — закреплённая в конституции или ином нормативном правовом акте возможность определённого поведения человека (например, свобода слова, свобода вероисповедания). Категория свободы близка к понятию права в субъективном смысле — субъективному праву, однако последнее предполагает наличие юридического механизма для реализации и обычно соответствующей обязанности государства или другого субъекта совершить какое-либо действие. Напротив, юридическая свобода не имеет чёткого механизма реализации, ей соответствует обязанность воздерживаться от совершения каких-либо нарушающих данную свободу действий.
        `)

	mustAddUserReturnId(defUsername, defPassword)

	rr := testPostNote(defUsername, defPassword, text)
	if rr.Code != http.StatusCreated {
		t.Error(rr.Body.String())
	}

	rr = testPostNote(defUsername, defPassword, textRus)
	if rr.Code != http.StatusCreated {
		t.Error(rr.Body.String())
	}

	rr = testPostNote(defUsername, defPassword, bigTest)
	if rr.Code != http.StatusCreated {
		t.Error(rr.Body.String())
	}

	rr = testPostNote(defUsername, defPassword, bigTextRus)
	if rr.Code != http.StatusCreated {
		t.Error(rr.Body.String())
	}
}

func TestPostNoteMistakesEng(t *testing.T) {
	mustClearDb(db)
	text := `This is mi textr`

	mustAddUserReturnId(defUsername, defPassword)
	if mustCountTable(usersTable) != 1 {
		t.Error("user not imported in database")
	}

	rr := testPostNote(defUsername, defPassword, text)
	if rr.Code != http.StatusBadRequest {
		msg := rr.Body.String()
		t.Error(msg)
	}

	if mustCountTable(notesTable) != 0 {
		t.Error("notes with mistakes were inserted in database")
	}
}

func TestPostNoteMistakesBigEng(t *testing.T) {
	mustClearDb(db)

	mustAddUserReturnId(defUsername, defPassword)
	if mustCountTable(usersTable) != 1 {
		t.Error("user not imported in database")
	}
	bigTest := url.QueryEscape(`Libertt is tejt state of being free within society from oppressive restrictions imposed by authority on one’s way of life, behavior, or political views.[1]

	In theolfy, liberty is freedom from the effects of "sin, spiritual servitude, [or] worldly ties".[2]
	
	In economics, liberty means free, fair, andhopen competition, often referred to as a free market.
	
	Sometimes liberty is differentiated from freedom by using the word "freedom" primarily, if not exclusively, to mean the ability to do as one wills and what one has the power to do; and using the word "liberty" to mean the absence of arbitrary restraints, taking into account the rights of all involved. In this sense, the exercise of liberty is subject to capability and limited by the rights of others. Thus liberty entails the responsible use of freedom under the rule of law without depriving anyone else of their freedom. Liberty can be taken away as a form of punishment. In many countries, people can be deprived of their liberty if they are convicted of criminal acts.
        `)

	rr := testPostNote(defUsername, defPassword, bigTest)
	if rr.Code != http.StatusBadRequest {
		msg := rr.Body.String()
		t.Error(msg)
	}

	if mustCountTable(notesTable) != 0 {
		t.Error("notes with mistakes were inserted in database")
	}
}

func TestPostNoteMistakesRus(t *testing.T) {
	mustClearDb(db)

	textRus := `Этом мрй тккст`

	mustAddUserReturnId(defUsername, defPassword)

	rr := testPostNote(defUsername, defPassword, textRus)
	if rr.Code != http.StatusBadRequest {
		msg := rr.Body.String()
		t.Error(msg)
	}

	if mustCountTable(notesTable) != 0 {
		t.Error("notes with mistakes were inserted in database")
	}
}

func TestPostNoteMistakesBigRus(t *testing.T) {
	mustClearDb(db)

	bigTextRus := url.QueryEscape(`В этм теусте очен мгого ошибак. Существует множество различных определений свободы[⇨]. В этике понимание свободы связано с наличием свободы воли человека.

        Свобода в философии — универсалия культуры субъектного ряда, фиксирующая возможность деятельности и поведения в условиях отсутствия внешнего целеполагания[3].

        Свобода личности в праве — закреплённая в конституции или ином нормативном правовом акте возможность определённого поведения человека (например, свобода слова, свобода вероисповедания). Категория свободы близка к понятию права в субъективном смысле — субъективному праву, однако последнее предполагает наличие юридического механизма для реализации и обычно соответствующей обязанности государства или другого субъекта совершить какое-либо действие. Напротив, юридическая свобода не имеет чёткого механизма реализации, ей соответствует обязанность воздерживаться от совершения каких-либо нарушающих данную свободу действий.
        `)

	mustAddUserReturnId(defUsername, defPassword)
	if mustCountTable(usersTable) != 1 {
		t.Error("user not imported in database")
	}
	rr := testPostNote(defUsername, defPassword, bigTextRus)
	if rr.Code != http.StatusBadRequest {
		msg := rr.Body.String()
		t.Error(msg)
	}

	if mustCountTable(notesTable) != 0 {
		t.Error("notes with mistakes were inserted in database")
	}
}

func TestGetNotes(t *testing.T) {
	mustClearDb(db)
	userId := mustAddUserReturnId(defUsername, defPassword)
	if mustCountTable(usersTable) != 1 {
		t.Error("user not imported in database")
	}
	note1ID := mustAddNoteReturnId(userId, "This is first note")
	note2ID := mustAddNoteReturnId(userId, "This is second note")
	note3ID := mustAddNoteReturnId(userId, "This is third note")

	if mustCountTable(notesTable) != 3 {
		t.Error("notes count mismatch")
	}

	expectedNotesID := []uuid.UUID{note1ID, note2ID, note3ID}
	rr := testGetNotes(defUsername, defPassword)
	if rr.Code != http.StatusOK {
		t.Error()
	}
	actualResult := parseGetNotes(t, rr)

	if len(actualResult.Notes) != mustCountTable(notesTable) {
		t.Error("notes count mismatch")
	}

	for _, note := range actualResult.Notes {
		if !util.Contains[uuid.UUID](expectedNotesID, note.Id) {
			t.Error()
		}
	}

	// Second user
	secondUsername := "newname"
	secondPassword := "password"
	secondUserId := mustAddUserReturnId(secondUsername, secondPassword)
	if mustCountTable(usersTable) != 2 {
		t.Error("user not imported in database")
	}
	note5ID := mustAddNoteReturnId(secondUserId, "New user. This is first note")
	note6ID := mustAddNoteReturnId(secondUserId, "New user. This is second note")

	if mustCountTable(notesTable) != 5 {
		t.Error("notes count mismatch")
	}

	expectedNotesID = []uuid.UUID{note5ID, note6ID}
	rr = testGetNotes(secondUsername, secondPassword)
	if rr.Code != http.StatusOK {
		t.Error()
	}
	actualResult = parseGetNotes(t, rr)

	// All note != User notes
	if len(actualResult.Notes) == mustCountTable(notesTable) {
		t.Error("notes count mismatch")
	}

	for _, note := range actualResult.Notes {
		if !util.Contains[uuid.UUID](expectedNotesID, note.Id) {
			t.Error()
		}
	}
}

func testRegister(username, password, method string) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	r := httptest.NewRequest(method, "/register", nil)
	r.SetBasicAuth(username, password)

	server.handleRegister(rr, r)
	return rr
}

func testPostNote(username, password, text string) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	text = fmt.Sprintf(`{"text": "%s"}`, text)
	body := bytes.NewBufferString(text)
	r := httptest.NewRequest(http.MethodPost, "/notes", body)
	r.SetBasicAuth(defUsername, defPassword)
	r.Header.Set("Content-Type", "application/json")

	server.handleNotes(rr, r)
	return rr
}

func testGetNotes(username, password string) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/notes", nil)
	r.SetBasicAuth(username, password)

	server.handleNotes(rr, r)
	return rr
}

func parseGetNotes(t *testing.T, resp *httptest.ResponseRecorder) database.NoteGetAll {
	var result database.NoteGetAll
	err := json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Error(err)
	}
	return result
}

func mustSetupServerAndDb() (*Server, *database.PostgresDatabase) {
	serverPort := util.MustGetenv(serverPortEnvVar)

	pgHost := util.MustGetenv(pgHostEnvVar)
	pgPassword := util.MustGetenv(pgPasswordEnvVar)

	postgres, err := database.NewPostgresDatabase(pgHost, pgPassword)
	if err != nil {
		log.Fatal(err)
	}
	err = postgres.Init()
	if err != nil {
		log.Fatal(err)
	}
	mustClearDb(postgres)

	yandexSpeller := spellcheck.NewYandexSpeller()
	server := NewServer(":"+serverPort, postgres, yandexSpeller)
	go server.Start()
	return server, postgres
}

func mustCountTable(table string) int {
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

func mustAddUserReturnId(username, password string) uuid.UUID {
	var id uuid.UUID

	creds := database.NewUserSecureCredentials(username, password)
	rows, err := db.Query("INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id;", creds.Username, creds.Password)
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

func mustAddNoteReturnId(user_id uuid.UUID, text string) uuid.UUID {
	var id uuid.UUID
	err := db.QueryRow("INSERT INTO notes (user_id, text) VALUES ($1, $2) RETURNING id", user_id, text).Scan(&id)
	if err != nil {
		panic(err)
	}
	return id
}

func mustClearDb(db *database.PostgresDatabase) {
	_, err := db.Exec(`
		DROP SCHEMA public CASCADE;
		CREATE SCHEMA public;
	`)
	if err != nil {
		panic(err)
	}
	db.Init()
}
