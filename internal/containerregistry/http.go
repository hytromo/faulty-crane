package containerregistry

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const baseURL = "https://eu.gcr.io/v2"

func (gcrClient GCRClient) newHTTPRequest(urlSuffix string) *http.Request {
	req, _ := http.NewRequest("GET", baseURL+urlSuffix, nil)
	req.SetBasicAuth("_token", gcrClient.AccessKey)
	return req
}

// getRequestTo does a GET request to the container registry and retries a few times on error
func (gcrClient GCRClient) getRequestTo(urlSuffix string) []byte {
	triesCount := 1

	sleepOrExitOnError := func(err error) {
		if triesCount > 3 {
			log.Fatalf("HTTP request failed many times, fatal error %v\n", err.Error())
		}

		log.Printf("HTTP request failed with %v, retrying...\n", err.Error())

		triesCount++

		time.Sleep(1000 * time.Millisecond)
	}

	for {
		resp, err := gcrClient.client.Do(
			gcrClient.newHTTPRequest(urlSuffix),
		)

		if err != nil {
			sleepOrExitOnError(err)
			continue
		}

		bodyBytes, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			sleepOrExitOnError(err)
			continue
		}

		return bodyBytes
	}
}
