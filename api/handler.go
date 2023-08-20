package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/tymbaca/kodenotes/log"
	"github.com/tymbaca/kodenotes/spellcheck"
)

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		msg, status := makeJsonError("method not allowed"), http.StatusMethodNotAllowed
		http.Error(w, msg, status)
		log.RequestInfo(r, msg, status)
		return
	}
	_, err := s.getUserId(r)
	if err == ErrCredsTooLong {
		msg, status := makeJsonError(err.Error()), http.StatusRequestEntityTooLarge
		http.Error(w, msg, status)
		log.RequestInfo(r, msg, status)
		return
	}
	if err == ErrParseBasicAuth {
		msg, status := makeJsonError(err.Error()), http.StatusBadRequest
		http.Error(w, msg, status)
		log.RequestInfo(r, msg, status)
		return
	}
	if err == nil {
		msg, status := makeJsonError("username already registred"), http.StatusUnprocessableEntity
		http.Error(w, msg, status)
		log.RequestInfo(r, msg, status)
		return
	}
	if err == ErrCantFindUser {
		secureCreds, err := getUserSecureCredentials(r)
		if err != nil {
			msg, status := makeJsonError(err.Error()), http.StatusBadRequest
			http.Error(w, msg, status)
			log.RequestInfo(r, msg, status)
			return
		}
		_, err = s.db.RegisterUser(secureCreds)
		if err != nil {
			// This case should not happen
			msg, status := makeJsonError(err.Error()), http.StatusInternalServerError
			http.Error(w, msg, status)
			log.RequestInfo(r, msg, status)
			return
		}
	}
	status := http.StatusCreated
	w.WriteHeader(status)
	log.RequestInfo(r, "--", status)
}

func (s *Server) handleNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userId, err := s.getUserIdIfAuthorized(r)
	if err != nil {
		msg, status := makeJsonError("not authorized"), http.StatusUnauthorized
		http.Error(w, msg, status)
		log.RequestInfo(r, msg, status)
		return
	}

	switch r.Method {
	case http.MethodGet:
		log.Info("%s\t | Redirecting to %s", r.URL, http.MethodGet)
		s.handleGetNotes(w, r, userId)
	case http.MethodPost:
		log.Info("%s\t | Redirecting to %s", r.URL, http.MethodPost)
		s.handlePostNote(w, r, userId)
	default:
		msg, status := makeJsonError(fmt.Sprintf("method %s not allowed. Use eather GET or POST", r.Method)), http.StatusMethodNotAllowed
		http.Error(w, msg, status)
		log.RequestInfo(r, msg, status)
		return
	}
}

func (s *Server) handleGetNotes(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	notes, err := s.db.GetNotes(userId)
	if err != nil {
		msg, status := makeJsonError(err.Error()), http.StatusInternalServerError
		http.Error(w, msg, status)
		log.RequestInfo(r, msg, status)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(notes)
	if err != nil {
		msg, status := makeJsonError(err.Error()), http.StatusInternalServerError
		http.Error(w, msg, status)
		log.RequestInfo(r, msg, status)
		return
	}
}

func (s *Server) handlePostNote(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	if r.Header.Get("Content-Type") != "application/json" {
		msg, status := makeJsonError("Content-Type must be application/json"), http.StatusBadRequest
		http.Error(w, msg, status)
		log.RequestInfo(r, msg, status)
		return
	}

	note, err := unmarshalNoteRequest(r)
	if err != nil {
		msg, status := makeJsonError(err.Error()), http.StatusBadRequest
		http.Error(w, msg, status)
		log.RequestInfo(r, msg, status)
		return
	}

	decodedText, err := url.QueryUnescape(note.Text)
	if err != nil {
		msg, status := makeJsonError("pass correctly URL Encoded note text inside of 'text' field"), http.StatusBadRequest
		http.Error(w, msg, status)
		log.RequestInfo(r, msg, status)
		return
	}
	note.Text = decodedText

	// Spell checker
	checkData, err := s.spellschecker.Check(note.Text)
	if err == spellcheck.ErrYandexTooBigText {
		msg, status := makeJsonError(err.Error()), http.StatusRequestEntityTooLarge
		http.Error(w, msg, status)
		log.RequestInfo(r, msg, status)
		return
	}
	if err != nil {
		msg, status := makeJsonError(err.Error()), http.StatusBadGateway
		http.Error(w, msg, status)
		log.RequestInfo(r, msg, status)
		return
	}

	if len(checkData) > 0 {
		// Bad text
		body, err := makeMisspellErrorBytes("spelling mistakes found", checkData)
		if err != nil {
			msg, status := makeJsonError("error while parsing spallchecker data"), http.StatusInternalServerError
			http.Error(w, msg, status)
			log.RequestInfo(r, msg, status)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(body)
		return
	} else if len(checkData) == 0 {
		// Good text
		err = s.db.PostNote(userId, note)
		if err != nil {
			msg, status := makeJsonError("bad credentials, make sure your username and password are less than 250 chars"), http.StatusBadRequest
			http.Error(w, msg, status)
			log.RequestInfo(r, msg, status)
			return
		}
		status := http.StatusCreated
		w.WriteHeader(status)
		log.RequestInfo(r, "--", status)
		return
	} else {
		msg, status := makeJsonError("spellchecker is broken"), http.StatusTeapot
		http.Error(w, msg, status)
		log.RequestInfo(r, msg, status)
		return
	}
}
