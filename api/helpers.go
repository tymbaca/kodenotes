package api

import (
	"crypto/sha256"
	"net/http"

	"github.com/google/uuid"
	"github.com/tymbaca/kodenotes/database"
)

// func authorized(r *http.Request) bool {
//         sessionKey := r.Header.Get("X-Session-Key")
//         if sessionKey == "" {
//                 return false
//         }
        
//         userId := sessions[sessionKey]
//         if userId == 0 {
//                 return false // 0 means key is not in map
//         } else {
//                 return true
//         }
// }

func (s *Server) getUserId(r *http.Request) uuid.NullUUID {
        username, password, ok := r.BasicAuth()
        if !ok {
                return uuid.NullUUID{Valid: false}
        }
        hashedPassword := hashString(password)
        creds := database.UserSecureCredentials{Username: username, Password: hashedPassword}

        id := s.db.GetUserId(creds)
        return id
}

func hashString(s string) string {
        h := sha256.New()
        h.Write([]byte(s))
        return string(h.Sum(nil))
}
