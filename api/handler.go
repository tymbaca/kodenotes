package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/tymbaca/kodenotes/spellcheck"
)

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, makeJsonError("method not allowed"), http.StatusMethodNotAllowed)
		return
	}
	userId := s.getUserId(r)
	if userId.Valid {
		http.Error(w, makeJsonError("username already registred"), http.StatusUnprocessableEntity)
		return
	} else {
		secureCreds, err := getUserSecureCredentials(r)
		if err != nil {
			http.Error(w, makeJsonError(err.Error()), http.StatusBadRequest)
			return
		}
		_, err = s.db.RegisterUser(secureCreds)
		if err != nil {
			// This case should not happen
			http.Error(w, makeJsonError(err.Error()), http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) handleNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// if !authorized(r) {
	//         http.Error(w, makeJsonError("not authorized"), http.StatusUnauthorized)
	// }

	result := s.getUserId(r)
	if !result.Valid {
		http.Error(w, makeJsonError("not authorized"), http.StatusUnauthorized)
		return
	}

	userId := result.UUID

	switch r.Method {
	case http.MethodGet:
		s.handleGetNotes(w, r, userId)
	case http.MethodPost:
		s.handlePostNote(w, r, userId)
	default:
		http.Error(w, fmt.Sprintf("method %s not allowed. Use eather GET or POST", r.Method), http.StatusMethodNotAllowed)
		return
	}
}

func (s *Server) handleGetNotes(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	notes, err := s.db.GetNotes(userId)
	if err != nil {
		http.Error(w, makeJsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(notes)
	if err != nil {
		http.Error(w, makeJsonError(err.Error()), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handlePostNote(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, makeJsonError("Content-Type must be application/json"), http.StatusBadRequest)
		return
	}

	note, err := unmarshalNoteRequest(r)
	if err != nil {
		http.Error(w, makeJsonError(err.Error()), http.StatusBadRequest)
		return
	}

	decodedText, err := url.QueryUnescape(note.Text)
	if err != nil {
		http.Error(w, makeJsonError("pass correctly URL Encoded note text inside of 'text' field"), http.StatusBadRequest)
		return
	}
	note.Text = decodedText

	result, err := s.spellschecker.Check(note.Text)
	if err == spellcheck.ErrYandexTooBigText {
		http.Error(w, makeJsonError(err.Error()), http.StatusRequestEntityTooLarge)
		return
	}
	if err != nil {
		http.Error(w, makeJsonError(err.Error()), http.StatusBadGateway)
		return
	}

	if len(result) > 0 {
		// Bad text
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&result)
		return
	} else if len(result) == 0 {
		// Good text
		err = s.db.PostNote(userId, note)
		if err != nil {
			http.Error(w, makeJsonError("bad credentials, make sure your username and password are less than 250 chars"), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
	} else {
		http.Error(w, makeJsonError("spellchecker is broken"), http.StatusTeapot)
		return
	}
}
