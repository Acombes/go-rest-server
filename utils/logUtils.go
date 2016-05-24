package utils

import (
	"log"
	"os"
)

var (
	logFileName      = "go-rest-server-log"
	iLog, eLog, wLog *log.Logger
)

func init() {
	// Initialize the logfile
	logFile, err := os.OpenFile(logFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	iLog = log.New(logFile, "         ", log.Ldate|log.Ltime)
	eLog = log.New(logFile, "ERROR:   ", log.Ldate|log.Ltime)
	wLog = log.New(logFile, "Warning: ", log.Ldate|log.Ltime)

	log.SetOutput(logFile)
}

func LogMessage(m string) {
	iLog.Println(m)
}

func LogError(m string) {
	eLog.Println(m)
}

func LogWarning(m string) {
	wLog.Println(m)
}
