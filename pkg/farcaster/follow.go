package farcaster

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type User struct {
	Fid         int    `json:"fid"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Pfp         struct {
		Url      string `json:"url"`
		Verified bool   `json:"verified"`
	} `json:"pfp"`
	Profile struct {
		Bio struct {
			Text     string   `json:"text"`
			Mentions []string `json:"mentions"`
		} `json:"bio"`
	} `json:"profile"`
	FollowerCount  int `json:"followerCount"`
	FollowingCount int `json:"followingCount"`
	ViewerContext  struct {
		Following  bool `json:"following"`
		FollowedBy bool `json:"followedBy"`
	} `json:"viewerContext"`
}

type usersResponse struct {
	Result struct {
		Users []User `json:"users"`
	} `json:"result"`
	Next struct {
		Cursor string `json:"cursor"`
	} `json:"next"`
}

func (client *FCRequestClient) GetFollowersPaginated(
	fid int,
	cursor string,
	limit int,
) ([]User, string, error) {
	followers := make([]User, 0, limit)

	query := make(url.Values, 3)
	query.Add("fid", fmt.Sprint(fid))
	query.Add("limit", fmt.Sprint(limit))
	if cursor != "" {
		query.Add("cursor", cursor)
	}

	requestUrl, err := url.JoinPath(client.BaseUrl, "/v2/followers")
	if err != nil {
		return followers, "", err
	}

	respBody, err := makeWarpcastRequest(
		requestUrl,
		query,
		client.appBearerToken,
		client.HTTPClient,
	)
	if err != nil {
		return followers, "", err
	}

	// Parsing the response
	var responseStruct usersResponse
	err = json.Unmarshal(respBody, &responseStruct)
	if err != nil {
		return followers, "", err
	}
	followers = responseStruct.Result.Users

	return followers, responseStruct.Next.Cursor, nil
}

func (client *FCRequestClient) GetFollowingPaginated(
	fid int,
	cursor string,
	limit int,
) ([]User, string, error) {
	following := make([]User, 0, limit)

	query := make(url.Values, 3)
	query.Add("fid", fmt.Sprint(fid))
	query.Add("limit", fmt.Sprint(limit))
	if cursor != "" {
		query.Add("cursor", cursor)
	}

	requestUrl, err := url.JoinPath(client.BaseUrl, "/v2/following")
	if err != nil {
		return following, "", err
	}

	respBody, err := makeWarpcastRequest(
		requestUrl,
		query,
		client.appBearerToken,
		client.HTTPClient,
	)
	if err != nil {
		return following, "", err
	}

	// Parsing the response
	var responseStruct usersResponse
	err = json.Unmarshal(respBody, &responseStruct)
	if err != nil {
		return following, "", err
	}
	following = responseStruct.Result.Users

	return following, responseStruct.Next.Cursor, nil
}
