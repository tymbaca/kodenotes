package api

import (
	"errors"
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

