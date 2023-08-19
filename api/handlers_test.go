package api

import (
	"log"
	"testing"

	"github.com/tymbaca/kodenotes/database"
	"github.com/tymbaca/kodenotes/spellcheck"
	"github.com/tymbaca/kodenotes/util"
)

const (
	serverPortEnvVar = "SERVER_PORT"
	pgHostEnvVar     = "POSTGRES_HOST"
	pgPasswordEnvVar = "POSTGRES_PASSWORD"
)

func TestMain(m *testing.M) {
        m.Run()
}

func mustSetupServer() *Server {
	serverPort := util.MustGetenv(serverPortEnvVar)

	pgHost := util.MustGetenv(pgHostEnvVar)
	pgPassword := util.MustGetenv(pgPasswordEnvVar)

	postgres, err := database.NewPostgresDatabase(pgHost, pgPassword)
	if err != nil {
		log.Fatal(err)
	}
	err = postgres.Init()
	if err != nil {
		log.Fatal(err)
	}

	yandexSpeller := spellcheck.NewYandexSpeller()
	server := NewServer(":"+serverPort, postgres, yandexSpeller)
        go server.Start()
        return server
}
