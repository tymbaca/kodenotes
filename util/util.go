package util

import (
	"log"
	"os"
)

func MustGetenv(key string) string {
        val := os.Getenv(key)
        if val == "" {
                log.Fatalf("please set %s env variable", key)
        }
        return val
}

func MustGetenvWithMessage(key, message string) string {
        val := os.Getenv(key)
        if val == "" {
                log.Fatal(message)
        }
        return val
}
