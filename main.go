package main

import (
	"log"

	"github.com/tymbaca/kodenotes/api"
	"github.com/tymbaca/kodenotes/database"
	"github.com/tymbaca/kodenotes/spellcheck"
	"github.com/tymbaca/kodenotes/util"
)

const (
	serverAddressEnvVar = "SERVER_ADDRESS"
	pgHostEnvVar        = "POSTGRES_HOST"
	pgPasswordEnvVar    = "POSTGRES_PASSWORD"
)

func main() {
	serverAddress := util.GetenvOrDefault(serverAddressEnvVar, ":8080")

	pgHost := util.MustGetenv(pgHostEnvVar)
	pgPassword := util.MustGetenv(pgPasswordEnvVar)

	postgres, err := database.NewPostgresDatabase(pgHost, pgPassword)
	if err != nil {
		log.Fatal(err)
	}

	yandexSpeller := spellcheck.NewYandexSpeller()
	server := api.NewServer(serverAddress, postgres, yandexSpeller)

	log.Printf("server running on address: %s", serverAddress)
	log.Fatal(server.Start())
}
