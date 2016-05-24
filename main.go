package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	FILE             string
	logFileName      = "go-rest-server-log"
	iLog, eLog, wLog *log.Logger
)

func init() {
	if len(os.Args) < 1 {
		return
	}

	FILE = os.Args[1]

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

func checkError(e error) bool {
	if e != nil {
		eLog.Println(e)
	}

	return e == nil
}

// Recover JSON from the loaded file
func readFromFile() (map[string]interface{}, error) {
	fileBytes, err := ioutil.ReadFile(FILE)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("File could not be found: %v", FILE))
	}
	return decodeJSON(string(fileBytes)).(map[string]interface{}), nil
}

func encodeJSON(elem interface{}) string {
	var output = new(bytes.Buffer)

	encoder := json.NewEncoder(output)

	encoder.Encode(elem)

	return output.String()
}

func decodeJSON(str string) interface{} {
	var output interface{}
	decoder := json.NewDecoder(strings.NewReader(str))

	for decoder.More() {
		err := decoder.Decode(&output)
		checkError(err)
	}

	return output
}

func getNewHandler(basePath string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		jsonData, err := readFromFile()

		if checkError(err) {
			dataBit := jsonData[basePath].([]interface{})
			iLog.Printf("Call to /%s/ - %d results\n", basePath, len(dataBit))
			w.Header().Set("Content-type", "application/json")
			fmt.Fprint(w, encodeJSON(dataBit))
		} else {
			// TODO : 500 error does not work
			http.Error(w, "Data source file could not be found", http.StatusInternalServerError)
		}
	}
}

func main() {
	// Recover the data from the file to build the API
	jsonData, err := readFromFile()
	checkError(err)

	// For each root key, expose a URL path
	for key, value := range jsonData {
		switch value.(type) {
		case []interface{}:
			http.HandleFunc("/"+key+"/", getNewHandler(key))
		}
	}

	// Test URL path
	http.HandleFunc("/test/", func(w http.ResponseWriter, r *http.Request) {
		var jsonData, err = readFromFile()
		checkError(err)

		for k, v := range jsonData {
			switch vv := v.(type) {
			case string:
				fmt.Fprintln(w, k, "is string", vv)
			case float64:
				fmt.Fprintln(w, k, "is float64", vv)
			case []interface{}:
				fmt.Fprintln(w, k, "is an array:")
				for i, u := range vv {
					fmt.Fprintln(w, i, u)
				}
			default:
				fmt.Fprintln(w, k, "is of a type I don't know how to handle: ", vv)
			}
		}
	})

	// Set up the server
	iLog.Println("Server startup")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
