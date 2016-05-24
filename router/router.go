package router

import (
	"../utils"
	"fmt"
	"github.com/tyler-sommer/stick"
	"net/http"
	"strings"
)

//*************************************************************************************************************************************
// Exported functions

func MakeRoutes(readFromFile func() (map[string]interface{}, error)) {
	// Recover the data from the file to build the API
	jsonData, err := readFromFile()
	utils.CheckError(err)

	// Greet the user and show the resources available to them
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var jsonData, err = readFromFile()
		if !utils.CheckError(err) {
			return
		}

		tmp := []string{}

		for key, value := range jsonData {
			switch value.(type) {
			case []interface{}:
				tmp = append(tmp, key)
			}
		}

		stick.NewEnv(stick.NewFilesystemLoader(".")).Execute("listing.twig", w, map[string]stick.Value{"roots": tmp})
	})

	// Serve the static files, such as the stylesheets, scripts, etc.
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		replacer := strings.NewReplacer("/static", "static")
		splittedPath := strings.Split(r.URL.Path, ".")

		w.Header().Set("Content-type", utils.GetMime(splittedPath[len(splittedPath)-1]))

		filePath := replacer.Replace(r.URL.Path)
		http.ServeFile(w, r, filePath)
	})

	// Show the entire database
	http.HandleFunc("/db/", func(w http.ResponseWriter, r *http.Request) {
		var jsonData, err = readFromFile()
		if !utils.CheckError(err) {
			return
		}
		prepareJSONResponse(&w)
		fmt.Fprint(w, utils.EncodeJSON(jsonData))
	})

	// For each root key, expose a URL path
	for key, value := range jsonData {
		switch value.(type) {
		case []interface{}:
			http.HandleFunc("/"+key+"/", getNewPathHandler(key, readFromFile))
		}
	}
}

func StartServer(p string) {
	// Set up the server
	utils.LogMessage("Server startup")
	http.ListenAndServe(p, nil)
}

//*************************************************************************************************************************************
// Package functions

func getNewPathHandler(basePath string, readFromFile func() (map[string]interface{}, error)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		jsonData, err := readFromFile()

		if utils.CheckError(err) {
			dataBit := jsonData[basePath].([]interface{})
			utils.LogMessage(fmt.Sprintf("Call to /%s/ - %d results\n", basePath, len(dataBit)))
			prepareJSONResponse(&w)
			fmt.Fprint(w, utils.EncodeJSON(dataBit))
		} else {
			http.Error(w, "Data source file could not be found", http.StatusInternalServerError)
		}
	}
}

func prepareJSONResponse(w *http.ResponseWriter) {
	(*w).Header().Set("Content-type", utils.GetMime("json"))
}
