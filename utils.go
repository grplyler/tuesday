package main

import (
	"log"
	"math"
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

func ColorFromScaleProgress(progress string) uint8 {
	prog, err := strconv.Atoi(progress)
	CheckError("Could not convert progress to int", err)
	var OldMin float64 = 0.0
	var OldMax float64 = 100.0
	var NewMin float64 = 17.0
	var NewMax float64 = 21.0
	var OldValue float64 = float64(prog)

	OldRange := (OldMax - OldMin)
	NewRange := (NewMax - NewMin)
	NewValue := (((OldValue - OldMin) * NewRange) / OldRange) + NewMin
	return uint8(math.Round(NewValue))
}
