package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/tymbaca/kodenotes/database"
	"github.com/tymbaca/kodenotes/spellcheck"
)

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userId := s.getUserId(r)
	if userId.Valid {
		http.Error(w, "username already registred", http.StatusBadRequest)
		return
	} else {
		secureCreds, err := getUserSecureCredentials(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, err = s.db.RegisterUser(secureCreds)
		if err != nil {
			// This case should not happen
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) handleNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// if !authorized(r) {
	//         http.Error(w, "not authorized", http.StatusUnauthorized)
	// }

	result := s.getUserId(r)
	if !result.Valid {
		http.Error(w, "not authorized", http.StatusUnauthorized)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(notes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) handlePostNote(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	var note database.NoteCreate

	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		http.Error(w, "pass json with 'text' field containing text of note", http.StatusBadRequest)
		return
	}

	result, err := s.spellschecker.Check(note.Text)
	if err == spellcheck.ErrYandexTooBigText {
		http.Error(w, err.Error(), http.StatusRequestEntityTooLarge)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	if len(result) > 0 {
		// Bad text
		// TODO: return spellschecker response
		json.NewEncoder(w).Encode(&result)
	} else if len(result) == 0 {
		// Good text
		err = s.db.PostNote(userId, note)
		if err != nil {
			http.Error(w, "bad credentials, make sure your username and password are less than 250 chars", http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, "spellchecker is broken", http.StatusTeapot)
	}
}
