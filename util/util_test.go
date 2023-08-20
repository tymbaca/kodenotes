package util

import (
	"os"
	"testing"

	"github.com/google/uuid"
)

func TestMustGetenv(t *testing.T) {
        os.Setenv("TEST", "1234")

        val := MustGetenv("TEST")
        if val != "1234" {
                t.Error()
        }

        defer recoverMustPanic(t)

        val = MustGetenv("NOT_EXIST")
}

func TestMustGetenvWithMessage(t *testing.T) {
        os.Setenv("TEST", "1234")

        val := MustGetenvWithMessage("TEST", "Please set TEST env var")
        if val != "1234" {
                t.Error()
        }

        defer recoverMustPanic(t)

        val = MustGetenvWithMessage("NOT_EXIST", "Please set NOT_EXIST env var")
}

func TestGetenvOrDefault(t *testing.T) {
        os.Setenv("TEST", "1234")

        val := GetenvOrDefault("TEST", "default")
        if val != "1234" {
                t.Error()
        }

        val = GetenvOrDefault("NOT_EXIST", "default")
        if val != "default" {
                t.Error()
        }
}

func TestGetenvIntOrDefault(t *testing.T) {
        os.Setenv("TEST", "10")
        os.Setenv("NOT_INT", "myname")

        val := GetenvIntOrDefault("TEST", 20)
        if val != 10 {
                t.Error()
        }
        
        val = GetenvIntOrDefault("NOT_INT", 20)
        if val != 20 {
                t.Error()
        }

        val = GetenvIntOrDefault("NOT_EXIST", 20)
        if val != 20 {
                t.Error()
        }
}

func TestContains(t *testing.T) {
        list := []string{"milk", "sugar", "bread"}

        goodTarget := "sugar"
        if !Contains[string](list, goodTarget) {
                t.Error()
        }

        badTarget := "water"
        if Contains[string](list, badTarget) {
                t.Error()
        }

        id1, _ := uuid.NewRandom()
        id2, _ := uuid.NewRandom()
        id3, _ := uuid.NewRandom()
        id4, _ := uuid.NewRandom()
        idList := []uuid.UUID{id1, id2, id3, id4}
        
        if !Contains[uuid.UUID](idList, id2) {
                t.Error()
        }
        wrongId, _ := uuid.NewRandom()
        if Contains[uuid.UUID](idList, wrongId) {
                t.Error()
        }
}

func recoverMustPanic(t *testing.T) {
        if r := recover(); r == nil {
                // Panic not found -> Fail
                t.Error()
        }
}
