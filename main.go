package main

import (
	"log"

	"github.com/tymbaca/kodenotes/api"
	"github.com/tymbaca/kodenotes/database"
	"github.com/tymbaca/kodenotes/spellcheck"
	"github.com/tymbaca/kodenotes/util"
)

const (
	serverPortEnvVar = "SERVER_PORT"
	pgHostEnvVar     = "POSTGRES_HOST"
	pgPasswordEnvVar = "POSTGRES_PASSWORD"
)

func main() {
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
	server := api.NewServer(":"+serverPort, postgres, yandexSpeller)

	log.Printf("server running on port: %s", serverPort)
	log.Fatal(server.Start())
}
