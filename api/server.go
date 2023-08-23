package api

import (
	"net/http"

	"github.com/tymbaca/kodenotes/database"
	"github.com/tymbaca/kodenotes/spellcheck"
)

type Server struct {
	addr          string
	db            database.Database
	spellschecker spellcheck.SpellChecker
}

func NewServer(addr string, db database.Database, spellchecker spellcheck.SpellChecker) *Server {
	server := &Server{
		addr:          addr,
		db:            db,
		spellschecker: spellchecker,
	}
	return server
}

func (s *Server) Start() error {
	http.HandleFunc("/register", s.handleRegister)
	http.HandleFunc("/notes", s.handleNotes)
	return http.ListenAndServe(s.addr, nil)
}
