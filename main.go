package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	FILE string
	mux  http.ServeMux
)

func init() {

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readFromFile() map[string]interface{} {
	var jsonData interface{}
	b, err := ioutil.ReadFile(FILE)

	check(err)

	json.Unmarshal(b, &jsonData)

	return jsonData.(map[string]interface{})
}

func getNewHandler(basePath string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var jsonData = readFromFile()[basePath].([]interface{})

		for i, u := range jsonData {
			fmt.Fprintln(w, i, u)
		}

		fmt.Fprintln(w, ...)
	}
}

func main() {
	if len(os.Args) < 1 {
		return
	}

	FILE = os.Args[1]

	jsonData := readFromFile()

	for key, value := range jsonData {
		fmt.Println(key)
		switch value.(type) {
		case []interface{}:
			http.HandleFunc("/"+key, getNewHandler(key))
		}
	}

	/*	http.HandleFunc("/plop", func(w http.ResponseWriter, r *http.Request) {
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
	*/
	log.Fatal(http.ListenAndServe(":8080", nil))
}
