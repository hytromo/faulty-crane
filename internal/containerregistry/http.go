package containerregistry

import (
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// getBaseUrl returns the base url for the api requests, e.g. where the container registry root is
func (gcrClient GCRClient) getBaseURL() string {
	return "https://" + gcrClient.Host + "/v2"
}

func (gcrClient GCRClient) newHTTPRequest(urlSuffix string) *http.Request {
	req, _ := http.NewRequest("GET", gcrClient.getBaseURL()+urlSuffix, nil)
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

		log.Infof("HTTP request failed with %v, retrying...\n", err.Error())

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
