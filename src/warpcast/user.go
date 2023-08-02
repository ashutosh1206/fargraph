package warpcast

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type WarpcastUserInfoResponse struct {
	Result struct {
		User WarpcastUserInfo
	} `json:"result"`
}

func GetUserInfoByUsername(username string, appBearerToken string, client *http.Client) (WarpcastUserInfo, error) {
	var userInfo WarpcastUserInfo

	// Preparing the request
	request, err := http.NewRequest(
		"GET",
		"https://api.warpcast.com/v2/user-by-username",
		nil,
	)
	if err != nil {
		return userInfo, err
	}
	query := request.URL.Query()
	query.Add("username", username)
	request.URL.RawQuery = query.Encode()
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", appBearerToken))

	// Sending the request
	response, err := client.Do(request)
	if err != nil {
		return userInfo, err
	}
	defer response.Body.Close()

	// Validating status code of response
	if response.StatusCode > 400 || response.StatusCode < 200 {
		return userInfo, fmt.Errorf("invalid status code: %d", response.StatusCode)
	}

	// Reading response
	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return userInfo, err
	}

	// Parsing the response
	var responseStruct WarpcastUserInfoResponse
	err = json.Unmarshal(respBody, &responseStruct)
	if err != nil {
		return userInfo, err
	}
	userInfo = responseStruct.Result.User

	return userInfo, nil
}
