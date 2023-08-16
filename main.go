package main

import (
	"log"
	"os"

	"github.com/tymbaca/kodenotes/api"
	"github.com/tymbaca/kodenotes/database"
	"github.com/tymbaca/kodenotes/spellcheck"
	"github.com/tymbaca/kodenotes/util"
)

const (
        serverAddressEnvVar = "SERVER_ADDRESS"
        pgHostEnvVar = "POSTGRES_HOST"
        pgDbNameEnvVar = "POSTGRES_DB"
        pgUserEnvVar = "POSTGRES_USER"
        pgPasswordEnvVar = "POSTGRES_PASSWORD"
)

func main() {
        serverAddress := util.GetenvOrDefault(serverAddressEnvVar, ":8080")
        pgDbName := util.GetenvOrDefault(pgDbNameEnvVar, "postgres")

        pgHost := util.MustGetenv(pgHostEnvVar)
        pgUser := util.MustGetenv(pgUserEnvVar)
        pgPassword := util.MustGetenv(pgPasswordEnvVar)


        postgres, err := database.NewPostgresDatabase(pgHost, pgDbName, pgUser, pgPassword)
        if err != nil {
                log.Fatal(err)
        }

        yandexScpeller := spellcheck.NewYandexSpeller()
        server := api.NewServer(serverAddress, postgres, yandexScpeller)

        log.Fatal(server.Start())
}
