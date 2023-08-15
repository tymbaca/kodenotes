package main

import (
	"log"
	"os"

	"github.com/tymbaca/kodenotes/api"
	"github.com/tymbaca/kodenotes/database"
	"github.com/tymbaca/kodenotes/spellcheck"
)

const (
        serverAddressEnvVar = "SERVER_ADDRESS"
        pgAddressEnvVar = "POSTGRES_URL"
        pgUserEnvVar = "POSTGRES_USER"
        pgPasswordEnvVar = "POSTGRES_PASSWORD"
)

func main() {
        serverAddress := os.Getenv(serverAddressEnvVar)
        if serverAddress == "" {
                serverAddress = ":8080"
        }

        pgAddress := os.Getenv(pgAddressEnvVar)
        if pgAddress == "" {
                log.Fatalf("set PostgreSQL address in '%s' environment variable", pgAddressEnvVar)
        }

        pgUser := os.Getenv(pgUserEnvVar)
        if pgAddress == "" {
                log.Fatalf("set PostgreSQL user in '%s' environment variable", pgUserEnvVar)
        }

        pgPassword := os.Getenv(pgPasswordEnvVar)
        if pgAddress == "" {
                log.Fatalf("set PostgreSQL password in '%s' environment variable", pgPasswordEnvVar)
        }


        postgres := database.NewPostgresDatabase(pgAddress, pgUser, pgPassword)
        yandexScpeller := spellcheck.NewYandexSpeller()
        server := api.NewServer(serverAddress, postgres, yandexScpeller)

        log.Fatal(server.Start())
}
