package router

import (
	"../utils"
	"fmt"
	"net/http"
)

func getNewHandler(basePath string, readFromFile func() (map[string]interface{}, error)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		jsonData, err := readFromFile()

		if utils.CheckError(err) {
			dataBit := jsonData[basePath].([]interface{})
			utils.LogMessage(fmt.Sprint("Call to /%s/ - %d results\n", basePath, len(dataBit)))
			w.Header().Set("Content-type", "application/json")
			fmt.Fprint(w, utils.EncodeJSON(dataBit))
		} else {
			http.Error(w, "Data source file could not be found", http.StatusInternalServerError)
		}
	}
}

func MakeRoutes(readFromFile func() (map[string]interface{}, error)) {
	// Recover the data from the file to build the API
	jsonData, err := readFromFile()
	utils.CheckError(err)

	// For each root key, expose a URL path
	for key, value := range jsonData {
		switch value.(type) {
		case []interface{}:
			http.HandleFunc("/"+key+"/", getNewHandler(key, readFromFile))
		}
	}

	// Test URL path
	http.HandleFunc("/test/", func(w http.ResponseWriter, r *http.Request) {
		var jsonData, err = readFromFile()
		utils.CheckError(err)

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
}

func StartServer(p string) {
	// Set up the server
	utils.LogMessage("Server startup")
	http.ListenAndServe(p, nil)
}
