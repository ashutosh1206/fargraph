package farcaster

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type FCRequestClient struct {
	BaseUrl        string
	appBearerToken string
	HTTPClient     *http.Client
}

func GetFCRequestClient(baseUrl string, appBearerToken string, httpClient *http.Client) *FCRequestClient {
	return &FCRequestClient{BaseUrl: baseUrl, appBearerToken: appBearerToken, HTTPClient: httpClient}
}

func makeWarpcastRequest(
	url string,
	query url.Values,
	appBearerToken string,
	client *http.Client,
) ([]byte, error) {
	// All requests to the API are GET requests right now
	var responseBytes []byte

	// Preparing the request
	request, err := http.NewRequest(
		"GET",
		url,
		nil,
	)
	if err != nil {
		return responseBytes, err
	}
	// Will overwrite existing query parameters in the URL
	request.URL.RawQuery = query.Encode()
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", appBearerToken))

	// Sending the request
	response, err := client.Do(request)
	if err != nil {
		return responseBytes, err
	}
	defer response.Body.Close()

	// Validating status code of response
	if response.StatusCode > 400 || response.StatusCode < 200 {
		return responseBytes, fmt.Errorf("invalid status code: %d", response.StatusCode)
	}

	// Reading response
	responseBytes, err = io.ReadAll(response.Body)
	if err != nil {
		return responseBytes, err
	}

	return responseBytes, err
}
