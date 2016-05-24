package utils

import (
	"bytes"
	"encoding/json"
	"strings"
)

func EncodeJSON(elem interface{}) string {
	var output = new(bytes.Buffer)

	encoder := json.NewEncoder(output)

	encoder.Encode(elem)

	return output.String()
}

func DecodeJSON(str string) interface{} {
	var output interface{}
	decoder := json.NewDecoder(strings.NewReader(str))

	for decoder.More() {
		err := decoder.Decode(&output)
		CheckError(err)
	}

	return output
}
