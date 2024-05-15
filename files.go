package dbman

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

func LoadTomlFile[T any](file string, stru *T) {
	tomlData := string(OpenFile(file))
	_, err := toml.Decode(tomlData, &stru)
	if err != nil {
		log.Fatalln("Cannot parse file: %s\n" + file)
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
			log.Println("Can not open file: " + file)
		}
		return filedata
	} else {
		log.Println("File doesn't exists: " + file)
		return []byte{}
	}
}
