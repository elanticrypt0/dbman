package dbman

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/k23dev/dbman/errors"
)

func LoadTomlFile[T any](file string, stru *T) {
	if ExitsFile(file) {
		tomlData := string(OpenFile(file))
		_, err := toml.Decode(tomlData, &stru)
		if err != nil {
			log.Fatalln(errors.FileNotLoaded(file))
		}
	} else {
		log.Println(errors.FileNotExistError(file))
	}
}

func ExitsFile(filepath string) bool {
	if _, err := os.Stat(filepath); err != nil {
		return false
	}
	return true
}

func OpenFile(file string) []byte {
	if ExitsFile(file) {
		filedata, err := os.ReadFile(file)
		if err != nil {
			log.Println(errors.FileNotOpened(file))
		}
		return filedata
	} else {
		log.Println(errors.FileNotExistError(file))
		return []byte{}
	}
}
