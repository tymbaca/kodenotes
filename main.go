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
        pgHostEnvVar = "POSTGRES_HOST"
        pgUserEnvVar = "POSTGRES_USER"
        pgPasswordEnvVar = "POSTGRES_PASSWORD"
)

func main() {
        serverAddress := util.GetenvOrDefault(serverAddressEnvVar, ":8080")

        pgHost := util.MustGetenv(pgHostEnvVar)
        pgUser := util.MustGetenv(pgUserEnvVar)
        pgPassword := util.MustGetenv(pgPasswordEnvVar)


        postgres, err := database.NewPostgresDatabase(pgHost, pgUser, pgPassword)
        if err != nil {
                log.Fatal(err)
        }

        yandexScpeller := spellcheck.NewYandexSpeller()
        server := api.NewServer(serverAddress, postgres, yandexScpeller)

        log.Fatal(server.Start())
}
