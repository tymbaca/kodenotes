package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/tymbaca/kodenotes/database"
	"github.com/tymbaca/kodenotes/spellcheck"
)

func (s *Server) handleNotes(w http.ResponseWriter, r *http.Request) {
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
                http.Error(w, err.Error(), http.StatusInternalServerError)
        }

        result, err := s.spellschecker.Check(note.Text)
        if err == spellcheck.ErrYandexTooBigText {
                http.Error(w, err.Error(), http.StatusRequestEntityTooLarge)
        }
        if err != nil {
                
        }
        
        if len(result) > 0 {
                // Bad text
                // TODO: return spellschecker response
        } else if len(result) == 0 {
                // Good text
                err = s.db.PostNote(userId, note)
                if err != nil {
                        http.Error(w, "bad credentials, make sure your username and password are less than 250 chars", http.StatusBadRequest)
                }
        } else {
                http.Error(w, "i am a teapot", http.StatusTeapot)
        }
}
