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

// InjectAuthInRequest is a function to inject authorisation information on every request
type InjectAuthInRequest func(req *http.Request)

// Client is just a wrapper around the normal http client to provide some retry logic
type Client struct {
	BaseURL             string
	realClient          *http.Client
	InjectAuthInRequest InjectAuthInRequest
}

func (httpClient Client) newPOST(url string, jsonPayload []byte) *http.Request {
	req, _ := http.NewRequest("POST", httpClient.getFullURLFor(url), bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")

	if httpClient.InjectAuthInRequest != nil {
		httpClient.InjectAuthInRequest(req)
	}

	return req
}

func (httpClient *Client) getFullURLFor(url string) string {
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return url
	}

	return httpClient.BaseURL + url
}

func (httpClient Client) newGET(url string) *http.Request {
	req, _ := http.NewRequest("GET", httpClient.getFullURLFor(url), nil)

	if httpClient.InjectAuthInRequest != nil {
		httpClient.InjectAuthInRequest(req)
	}

	return req
}

func (httpClient Client) newDELETE(url string) *http.Request {
	req, _ := http.NewRequest("DELETE", httpClient.getFullURLFor(url), nil)

	if httpClient.InjectAuthInRequest != nil {
		httpClient.InjectAuthInRequest(req)
	}

	return req
}

// GetRequestTo does a GET request and retries a few times on error
func (httpClient Client) GetRequestTo(url string) ([]byte, error) {
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
func (httpClient Client) PostRequestTo(url string, jsonPayload []byte, allowCompleteFailure bool, silentErrors bool) ([]byte, error) {
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
func (httpClient Client) DeleteRequestTo(url string, allowCompleteFailure bool, silentErrors bool) error {
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

// NewClientParams is the parameters required to build a new client
type NewClientParams struct {
	BaseURL             string
	InjectAuthInRequest InjectAuthInRequest
}

// NewClient is building a new client
func NewClient(params NewClientParams) Client {
	return Client{
		BaseURL:             params.BaseURL,
		realClient:          &http.Client{},
		InjectAuthInRequest: params.InjectAuthInRequest,
	}
}
