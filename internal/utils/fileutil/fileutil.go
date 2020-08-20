package fileutil

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
)

// SaveJSON saves a struct as pretty-formatted JSON data to a specific path; it can optionally compress the final result
func SaveJSON(path string, dataToWrite interface{}, doCompress bool) error {
	bytesToWrite, _ := json.MarshalIndent(dataToWrite, "", "\t")

	if doCompress {
		var compressedData bytes.Buffer
		gz, err := gzip.NewWriterLevel(&compressedData, gzip.BestCompression)

		if err != nil {
			return err
		}

		_, err = gz.Write(bytesToWrite)

		if err != nil {
			return err
		}

		gz.Close()
		return ioutil.WriteFile(path, compressedData.Bytes(), 0644)
	}

	return ioutil.WriteFile(path, bytesToWrite, 0644)
}
