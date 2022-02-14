package http

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type InjectAuthInRequest func(req *http.Request)

// HttpClient is just a wrapper around the normal http client to provide some retry logic
type HttpClient struct {
	BaseUrl             string
	realClient          *http.Client
	InjectAuthInRequest InjectAuthInRequest
}

func (httpClient HttpClient) newPOST(url string, jsonPayload []byte) *http.Request {
	req, _ := http.NewRequest("POST", httpClient.getFullUrlFor(url), bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")

	if httpClient.InjectAuthInRequest != nil {
		httpClient.InjectAuthInRequest(req)
	}

	return req
}

func (httpClient *HttpClient) getFullUrlFor(url string) string {
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return url
	}

	return httpClient.BaseUrl + url
}

func (httpClient HttpClient) newGET(url string) *http.Request {
	req, _ := http.NewRequest("GET", httpClient.getFullUrlFor(url), nil)

	if httpClient.InjectAuthInRequest != nil {
		httpClient.InjectAuthInRequest(req)
	}

	return req
}

func (httpClient HttpClient) newDELETE(url string) *http.Request {
	req, _ := http.NewRequest("DELETE", httpClient.getFullUrlFor(url), nil)

	if httpClient.InjectAuthInRequest != nil {
		httpClient.InjectAuthInRequest(req)
	}

	return req
}

// GetRequestTo does a GET request and retries a few times on error
func (httpClient HttpClient) GetRequestTo(url string) ([]byte, error) {
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
			httpClient.newGET(url),
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

		return bodyBytes, err
	}
}

// PostRequestTo does a POST request and retries a few times on error
func (httpClient HttpClient) PostRequestTo(url string, jsonPayload []byte, allowCompleteFailure bool, silentErrors bool) ([]byte, error) {
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
			return nil, nil // request retried too many times but we don't care anymore
		}

		resp, err := httpClient.realClient.Do(
			httpClient.newPOST(url, jsonPayload),
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

		return bodyBytes, err
	}
}

// DeleteRequestTo does a DELETE request and retries a few times on error
func (httpClient HttpClient) DeleteRequestTo(url string, allowCompleteFailure bool, silentErrors bool) error {
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
			httpClient.newDELETE(url),
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
		InjectAuthInRequest: params.InjectAuthInRequest,
	}
}
