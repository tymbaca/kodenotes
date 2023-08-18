package util

import (
	"os"
	"testing"
)

func TestMustGetenv(t *testing.T) {
        os.Setenv("TEST", "1234")

        val := MustGetenv("TEST")
        if val != "1234" {
                t.FailNow()
        }

        defer recoverMustPanic(t)

        val = MustGetenv("NOT_EXIST")
}

func TestMustGetenvWithMessage(t *testing.T) {
        os.Setenv("TEST", "1234")

        val := MustGetenvWithMessage("TEST", "Please set TEST env var")
        if val != "1234" {
                t.FailNow()
        }

        defer recoverMustPanic(t)

        val = MustGetenvWithMessage("NOT_EXIST", "Please set NOT_EXIST env var")
}

func TestGetenvOrDefault(t *testing.T) {
        os.Setenv("TEST", "1234")

        val := GetenvOrDefault("TEST", "default")
        if val != "1234" {
                t.FailNow()
        }

        val = GetenvOrDefault("NOT_EXIST", "default")
        if val != "default" {
                t.FailNow()
        }
}

func TestGetenvIntOrDefault(t *testing.T) {
        os.Setenv("TEST", "10")
        os.Setenv("NOT_INT", "myname")

        val := GetenvIntOrDefault("TEST", 20)
        if val != 10 {
                t.FailNow()
        }
        
        val = GetenvIntOrDefault("NOT_INT", 20)
        if val != 20 {
                t.FailNow()
        }

        val = GetenvIntOrDefault("NOT_EXIST", 20)
        if val != 20 {
                t.FailNow()
        }
}

func recoverMustPanic(t *testing.T) {
        if r := recover(); r == nil {
                // Panic not found -> Fail
                t.FailNow()
        }
}
