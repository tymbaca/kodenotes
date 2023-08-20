package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/tymbaca/kodenotes/database"
	"github.com/tymbaca/kodenotes/spellcheck"
)

var (
	ErrCredsTooLong = errors.New("credentials are too long")
	ErrCantFindUser = errors.New("username is not in database")
	ErrUnauthorized = errors.New("can't authorize")
)

func (s *Server) getUserId(r *http.Request) (uuid.UUID, error) {
	username, _, ok := r.BasicAuth()
	if !ok {
		return uuid.UUID{}, errors.New("can't parse basic auth")
	}

	if len(username) > s.db.MaxCredsLength() {
		return uuid.UUID{}, ErrCredsTooLong
	}
	id := s.db.GetUserId(username)
	if id.Valid {
		return id.UUID, nil
	} else {
		return uuid.UUID{}, ErrCantFindUser
	}
}

func (s *Server) getUserIdIfAuthorized(r *http.Request) (uuid.UUID, error) {
	username, password, ok := r.BasicAuth()
	if !ok {
		return uuid.UUID{}, errors.New("can't parse basic auth")
	}
	creds := database.NewUserSecureCredentials(username, password)

	if len(creds.Username) > s.db.MaxCredsLength() || len(creds.Password) > s.db.MaxCredsLength() {
		return uuid.UUID{}, ErrCredsTooLong
	}
	id := s.db.GetUserIdIfAuthorized(creds)
	if id.Valid {
		return id.UUID, nil
	} else {
		return uuid.UUID{}, ErrUnauthorized
	}
}

func getUserSecureCredentials(r *http.Request) (database.UserSecureCredentials, error) {
	username, password, ok := r.BasicAuth()
	if !ok {
		return database.UserSecureCredentials{}, errors.New("error while parsing credentials from request")
	}
	creds := database.NewUserSecureCredentials(username, password)
	return creds, nil
}

func unmarshalNoteRequest(r *http.Request) (database.NoteCreate, error) {
	var body map[string]string // smells
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return database.NoteCreate{}, errors.New("wrong json format, pass json object with 'text' field")
	}

	val, ok := body["text"]
	if ok {
		note := database.NoteCreate{Text: val}
		return note, nil
	} else {
		return database.NoteCreate{}, errors.New("pass json with 'text' field containing text of note")
	}
}

func makeJsonError(msg string) string {
	return fmt.Sprintf(`{"detail":"%s"}`, msg)
}

func makeMisspellErrorBytes(msg string, checkData spellcheck.CheckResponse) ([]byte, error) {
	body := map[string]interface{}{
		"detail":    msg,
		"misspells": checkData,
	}
	bodyEncoded, err := json.Marshal(body)
	if err != nil {
		return nil, err
	} else {
		return bodyEncoded, nil
	}
}
