package util

import (
	"log"
	"os"
)

func MustGetenv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Panicf("please set %s env variable", key)
	}
	return val
}

func MustGetenvWithMessage(key, message string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Panic(message)
	}
	return val
}

func GetenvOrDefault(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		val = def
	}
	return val
}
