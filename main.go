package main

import (
	"./router"
	"./utils"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	FILE string
)

func init() {
	if len(os.Args) < 2 {
		fmt.Println("No source file given...\n\n\n")
		os.Exit(1)
	}

	FILE = os.Args[1]
}

// Recover JSON from the loaded file
func readFromFile() (map[string]interface{}, error) {
	fileBytes, err := ioutil.ReadFile(FILE)
	if err != nil {
		return nil, utils.GetNewError(fmt.Sprintf("File could not be found: %v", FILE))
	}
	return utils.DecodeJSON(string(fileBytes)).(map[string]interface{}), nil
}

func main() {
	router.MakeRoutes(readFromFile)
	router.StartServer(":8080")
}
