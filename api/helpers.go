package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/tymbaca/kodenotes/database"
)

func (s *Server) getUserId(r *http.Request) uuid.NullUUID {
	username, password, ok := r.BasicAuth()
	if !ok {
		return uuid.NullUUID{Valid: false}
	}
	creds := database.NewUserSecureCredentials(username, password)

	id := s.db.GetUserId(creds)
	return id
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
		return database.NoteCreate{}, errors.New("wrong json format")
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
