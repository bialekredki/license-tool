package common

import (
	"errors"
	"os"

	log "github.com/sirupsen/logrus"
)

func IsExistingDirectory(directoryName string) bool {
	if stat, err := os.Stat(directoryName); errors.Is(err, os.ErrNotExist) || !stat.IsDir() {
		return false
	}
	return true
}

func IsExisitingFile(filename string) bool {
	stat, err := os.Stat(filename)
	return !(errors.Is(err, os.ErrNotExist) || stat.IsDir())
}

func ListContentOfDirectory(directoryName string) ([]string, []string) {
	entries, err := os.ReadDir(directoryName)
	if err != nil {
		log.Fatal(err)
	}

	var files []string
	var directories []string

	for _, e := range entries {
		name := e.Name()
		if e.IsDir() {
			directories = append(directories, name)
		} else {
			files = append(files, name)
		}
	}
	return files, directories
}
