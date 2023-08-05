package warpcast

import (
	"encoding/json"
	"net/url"
)

type WarpcastUserInfoResponse struct {
	Result struct {
		User WarpcastUserInfo
	} `json:"result"`
}

func (client *FCRequestClient) GetUserInfoByUsername(username string) (WarpcastUserInfo, error) {
	var userInfo WarpcastUserInfo

	query := make(url.Values, 1)
	query.Add("username", username)

	respBody, err := makeWarpcastRequest(
		"https://api.warpcast.com/v2/user-by-username",
		query,
		client.appBearerToken,
		client.HTTPClient,
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
