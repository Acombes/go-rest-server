package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	FILE string
)

func init() {
	if len(os.Args) < 1 {
		return
	}

	FILE = os.Args[1]
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Recover JSON from the loaded file
func readFromFile() map[string]interface{} {
	fileBytes, err := ioutil.ReadFile(FILE)
	if err != nil {
		fmt.Println("File could not be found: ", FILE)
		os.Exit(1)
	}
	return decodeJSON(string(fileBytes)).(map[string]interface{})
}

func encodeJSON(elem interface{}) string {
	var output = new(bytes.Buffer)

	encoder := json.NewEncoder(output)

	encoder.Encode(elem)

	return output.String()
}

func decodeJSON(str string) interface{} {
	var (
		output interface{}
		err    error
	)
	decoder := json.NewDecoder(strings.NewReader(str))

	for decoder.More() {
		err = decoder.Decode(&output)
		check(err)
	}

	return output
}

func getNewHandler(basePath string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var jsonData = readFromFile()[basePath].([]interface{})

		w.Header().Set("Content-type", "application/json")
		fmt.Fprint(w, encodeJSON(jsonData))
	}
}

func main() {
	// Recover the data from the file to build the API
	jsonData := readFromFile()

	// For each root key, expose a URL path
	for key, value := range jsonData {
		switch value.(type) {
		case []interface{}:
			http.HandleFunc("/"+key+"/", getNewHandler(key))
		}
	}

	// Test URL path
	http.HandleFunc("/test/", func(w http.ResponseWriter, r *http.Request) {
		var jsonData = readFromFile()
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
	log.Fatal(http.ListenAndServe(":8080", nil))
}
