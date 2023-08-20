package main

import (
	"errors"
	"io"
	"os"

	"github.com/tymbaca/kodenotes/api"
	"github.com/tymbaca/kodenotes/database"
	"github.com/tymbaca/kodenotes/log"
	"github.com/tymbaca/kodenotes/spellcheck"
	"github.com/tymbaca/kodenotes/util"
)

const (
	serverPortEnvVar = "SERVER_PORT"
	pgHostEnvVar     = "POSTGRES_HOST"
	pgPasswordEnvVar = "POSTGRES_PASSWORD"
)

func createLogWriter() io.Writer {
	path := "logs/server.log"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		_, err2 := os.Create(path)
		if err2 != nil {
			panic(err2)
		}
	}
	logFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	multiOut := io.MultiWriter(os.Stdout, logFile)
	return multiOut
}

func main() {
	logOutput := createLogWriter()
	log.SetOutput(logOutput)
	serverPort := util.MustGetenv(serverPortEnvVar)

	pgHost := util.MustGetenv(pgHostEnvVar)
	pgPassword := util.MustGetenv(pgPasswordEnvVar)

	postgres, err := database.NewPostgresDatabase(pgHost, pgPassword)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Connected to PostgreSQL addr: %s", pgHost)
	err = postgres.Init()
	if err != nil {
		log.Fatal(err)
	}
	log.Info("PosgreSQL successfully initialized")

	yandexSpeller := spellcheck.NewYandexSpeller()
	server := api.NewServer(":"+serverPort, postgres, yandexSpeller)
	log.Info("Created main server")

	log.Info("Server running on port: %s", serverPort)
	log.Fatal(server.Start())
}
