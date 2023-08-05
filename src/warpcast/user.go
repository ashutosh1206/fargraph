package warpcast

import (
	"encoding/json"
	"net/url"
)

type userResponse struct {
	Result struct {
		User User
	} `json:"result"`
}

func (client *FCRequestClient) GetUserInfoByUsername(username string) (User, error) {
	var userInfo User

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
	var responseStruct userResponse
	err = json.Unmarshal(respBody, &responseStruct)
	if err != nil {
		return userInfo, err
	}
	userInfo = responseStruct.Result.User

	return userInfo, nil
}
