package util

import (
	"log"
	"os"
	"strconv"
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

func GetenvIntOrDefault(key string, def int) int {
        val := os.Getenv(key)
        if val == "" {
                return def
        }
        valInt, err := strconv.Atoi(val)
        if err != nil {
                return def
        } 
        return valInt 
}

func Contains[T comparable](slice []T, target T) bool {
        for _, item := range slice {
                if item == target {
                        return true
                }
        }
        return false
}
