package helper

import (
	"errors"
	"os"
	"path/filepath"
)

func FindGengoLocalRoot() (string, error) {
	var dir string
	// verify if gengo initialized is in .gengo directory
	_, err := os.Stat(".gengo/gengo.json")
	if err != nil && os.IsNotExist(err) {
		// verify if gengo initialized is in current directory
		_, err = os.Stat("gengo.json")
		if err != nil && os.IsNotExist(err) {
			return "", errors.New("Gengo not initialized")
		} else if err != nil {
			return "", err
		} else {
			dir = filepath.Join("./")
		}
	} else if err != nil {
		return "", err
	} else {
		dir = filepath.Join("./", ".gengo")
	}
	return dir, nil
}
