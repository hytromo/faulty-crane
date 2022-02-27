package fileutil

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

type toWriteType struct {
	Name     string
	Username string
}

func TestReadWrite(t *testing.T) {
	file, err := ioutil.TempFile("", "unit_test")
	file.Name()

	if err != nil {
		t.Error("Could not create temp file")
	}

	defer os.Remove(file.Name())

	err = SaveJSON(file.Name(), toWriteType{
		Name:     "alex",
		Username: "hytromo",
	}, true)

	if err != nil {
		t.Error("Could not save JSON")
	}

	fileBytes, err := ReadFile(file.Name(), true)
	if err != nil {
		t.Errorf("Could not read file '%v': %v\n", file.Name(), err.Error())
	}

	parsedData := toWriteType{}

	err = json.Unmarshal([]byte(fileBytes), &parsedData)

	if err != nil {
		log.Fatalf("Cannot parse json of plan file %v: %v\n", file.Name(), err.Error())
	}

	if parsedData.Name != "alex" || parsedData.Username != "hytromo" {
		t.Error("Wrong parsed data values")
	}
}
