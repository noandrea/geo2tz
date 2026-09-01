package helpers

import (
	"encoding/json"
	"fmt"
	"os"
)

func SaveJSON(data any, toFile string) (err error) {
	// write the version file as json, overwriting the old one
	jsonString, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return
	}
	if err = os.WriteFile(toFile, jsonString, 0o600); err != nil {
		return
	}
	return
}

func LoadJSON(fromFile string, data any) (err error) {
	// read the version file
	jsonFile, err := os.Open(fromFile)
	if err != nil {
		return
	}
	defer func() {
		if closeErr := jsonFile.Close(); closeErr != nil {
			fmt.Println("Error closing JSON file:", closeErr)
		}
	}()

	jsonParser := json.NewDecoder(jsonFile)
	if err = jsonParser.Decode(data); err != nil {
		return
	}
	return
}

func DeleteQuietly(filePath ...string) (err error) {
	for _, filePath := range filePath {
		if xErr := FileExists(filePath); xErr == nil {
			fmt.Println("removing file", filePath)
			if err = os.Remove(filePath); err != nil {
				return
			}
		}
	}
	return
}

func FileExists(filePath string) error {
	f, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("file does not exist: %s", filePath)
	}
	if f.IsDir() {
		return fmt.Errorf("file is a directory: %s", filePath)
	}
	return nil
}
