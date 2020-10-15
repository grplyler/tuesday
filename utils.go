package main

import (
	"log"
	"os/user"
	"path/filepath"
	"strconv"
)

func CheckError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

func DataPath() string {
	usr, err := user.Current()

	if err != nil {
		log.Fatal(err)
	}
	dataPath := filepath.Join(usr.HomeDir, ".tuesday/", "main.csv")
	return dataPath
}

func IsDigit(num string) bool {
	if _, err := strconv.Atoi(num); err == nil {
		return true
	} else {
		return false
	}
}
