package api

import (
        "fmt"
	"encoding/json"
	"io"
	"net/http"
)

func (s *Server) handleNotes(w http.ResponseWriter, r *http.Request) {
        if !authorized(r) {
                http.Error(w, "not authorized", http.StatusUnauthorized)
        }

        userId := getUserId(r)

        switch r.Method {
        case http.MethodGet:
                s.handleGetNotes(w, r, userId)
        case http.MethodPost:
                s.handlePostNote(w, r, userId)
        default:
                http.Error(w, fmt.Sprintf("method %s not allowed. Use eather GET or POST", r.Method), http.StatusMethodNotAllowed)
        }
}


func (s *Server) handleGetNotes(w http.ResponseWriter, r *http.Request, userId int) {
        // notes := s.db.GetNotes(userId)
        json.NewEncoder(w)
}

func (s *Server) handlePostNote(w http.ResponseWriter, r *http.Request, userId int) {
        textBytes, err := io.ReadAll(r.Body)
        text := string(textBytes)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
        }

        s.spellschecker.Check(text)
}
