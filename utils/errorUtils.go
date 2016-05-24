package utils

import (
	"errors"
	"fmt"
)

func CheckError(e error) bool {
	if e != nil {
		LogError(fmt.Sprint(e))
	}

	return e == nil
}

func GetNewError(m string) (err error) {
	err = errors.New(m)
	return
}
