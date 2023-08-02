package warpcast

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type WarpcastUserInfoResponse struct {
	Result struct {
		User WarpcastUserInfo
	} `json:"result"`
}

func GetUserInfoByUsername(username string, appBearerToken string, client *http.Client) (WarpcastUserInfo, error) {
	var userInfo WarpcastUserInfo

	query := make(url.Values, 1)
	query.Add("username", username)

	respBody, err := makeWarpcastRequest(
		"https://api.warpcast.com/v2/user-by-username",
		query,
		appBearerToken,
		client,
	)
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
