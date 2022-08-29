package json

import (
	"encoding/json"
	"os"
)

func LoadFile[R any](filePath string) (fileContent R, err error) {
	var configFile *os.File
	configFile, err = os.Open(filePath)
	defer configFile.Close()
	if err != nil {
		return
	}

	err = json.NewDecoder(configFile).Decode(&fileContent)
	return
}
