package http

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type InjectAuthInRequest func(req *http.Request)

// HttpClient is just a wrapper around the normal http client to provide some retry logic
type HttpClient struct {
	BaseUrl             string
	realClient          *http.Client
	injectAuthInRequest InjectAuthInRequest
}

func (httpClient HttpClient) newGET(urlSuffix string) *http.Request {
	fmt.Println("Getting", httpClient.BaseUrl+urlSuffix)

	req, _ := http.NewRequest("GET", httpClient.BaseUrl+urlSuffix, nil)

	httpClient.injectAuthInRequest(req)

	return req
}

func (httpClient HttpClient) newDELETE(urlSuffix string) *http.Request {
	req, _ := http.NewRequest("DELETE", httpClient.BaseUrl+urlSuffix, nil)

	httpClient.injectAuthInRequest(req)

	return req
}

// GetRequestTo does a GET request to the container registry and retries a few times on error
func (httpClient HttpClient) GetRequestTo(urlSuffix string) []byte {
	triesCount := 1

	sleepOrExitOnError := func(err error) {
		if triesCount > 3 {
			log.Fatalf("HTTP request failed many times, fatal error: %v\n", err.Error())
		}

		log.Infof("HTTP request failed with %v, retrying...\n", err.Error())

		triesCount++

		time.Sleep(1000 * time.Millisecond)
	}

	for {
		resp, err := httpClient.realClient.Do(
			httpClient.newGET(urlSuffix),
		)

		if err != nil {
			sleepOrExitOnError(err)
			continue
		}

		bodyBytes, err := ioutil.ReadAll(resp.Body)

		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			sleepOrExitOnError(errors.New(string(bodyBytes)))
			continue
		}

		if err != nil {
			sleepOrExitOnError(err)
			continue
		}

		return bodyBytes
	}
}

// DeleteRequestTo does a DELETE request to the container registry and retries a few times on error
func (httpClient HttpClient) DeleteRequestTo(urlSuffix string, allowCompleteFailure bool, silentErrors bool) error {
	triesCount := 1

	sleepOrExitOnError := func(err error) {
		if triesCount > 3 && !allowCompleteFailure {
			log.Fatalf("HTTP request failed many times, fatal error: %v\n", err.Error())
		}

		if !silentErrors {
			log.Infof("HTTP request failed with %v, retrying...\n", err.Error())
		}

		triesCount++

		time.Sleep(1000 * time.Millisecond)
	}

	for {
		if triesCount >= 4 && allowCompleteFailure {
			return nil // request retried too many times but we don't care anymore
		}

		resp, err := httpClient.realClient.Do(
			httpClient.newDELETE(urlSuffix),
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

		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			sleepOrExitOnError(errors.New(string(bodyBytes)))
			continue
		}

		return err
	}
}

type NewHttpClientParams struct {
	BaseUrl             string
	InjectAuthInRequest InjectAuthInRequest
}

func NewHttpClient(params NewHttpClientParams) HttpClient {
	return HttpClient{
		BaseUrl:             params.BaseUrl,
		realClient:          &http.Client{},
		injectAuthInRequest: params.InjectAuthInRequest,
	}
}
