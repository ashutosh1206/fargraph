package farcaster

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type userResponse struct {
	Result struct {
		User    User `json:"user"`
		Inviter User `json:"inviter"`
	} `json:"result"`
}

func (client *FCRequestClient) GetUserInfoByUsername(username string) (User, error) {
	var userInfo User

	query := make(url.Values, 1)
	query.Add("username", username)

	requestUrl, err := url.JoinPath(client.BaseUrl, "/v2/user-by-username")
	if err != nil {
		return userInfo, err
	}

	respBody, err := makeWarpcastRequest(
		requestUrl,
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

func (client *FCRequestClient) GetRecentUsersPaginated(cursor string, limit int) ([]User, string, error) {
	users := make([]User, 0, limit)

	query := make(url.Values, 2)
	query.Add("limit", fmt.Sprint(limit))
	if cursor != "" {
		query.Add("cursor", cursor)
	}

	requestUrl, err := url.JoinPath(client.BaseUrl, "/v2/user-liked-casts")
	if err != nil {
		return users, "", err
	}

	respBody, err := makeWarpcastRequest(requestUrl, query, client.appBearerToken, client.HTTPClient)
	if err != nil {
		return users, "", err
	}

	// Parsing the response
	var responseStruct usersResponse
	err = json.Unmarshal(respBody, &responseStruct)
	if err != nil {
		return users, "", err
	}
	users = responseStruct.Result.Users

	return users, responseStruct.Next.Cursor, nil
}

func (client *FCRequestClient) GetUserByFid(fid int) (User, User, error) {
	var userInfo User
	var inviterInfo User

	query := make(url.Values, 1)
	query.Add("fid", fmt.Sprint(fid))

	requestUrl, err := url.JoinPath(client.BaseUrl, "/v2/user")
	if err != nil {
		return userInfo, inviterInfo, err
	}

	respBody, err := makeWarpcastRequest(
		requestUrl,
		query,
		client.appBearerToken,
		client.HTTPClient,
	)
	if err != nil {
		return userInfo, inviterInfo, err
	}

	// Parsing the response
	var responseStruct userResponse
	err = json.Unmarshal(respBody, &responseStruct)
	if err != nil {
		return userInfo, inviterInfo, err
	}
	userInfo = responseStruct.Result.User
	inviterInfo = responseStruct.Result.Inviter

	return userInfo, inviterInfo, nil
}
