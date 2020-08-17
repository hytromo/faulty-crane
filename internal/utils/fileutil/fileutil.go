package fileutil

import (
	"encoding/json"
	"io/ioutil"
)

// SaveJSON saves a struct as pretty-formatted JSON data to a specific path
func SaveJSON(path string, dataToWrite interface{}) error {
	bytesToWrite, _ := json.MarshalIndent(dataToWrite, "", "\t")

	return ioutil.WriteFile(path, bytesToWrite, 0644)
}
