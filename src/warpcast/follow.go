package warpcast

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type WarpcastUserInfo struct {
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

type WarpcastResponse struct {
	Result struct {
		Users []WarpcastUserInfo `json:"users"`
	} `json:"result"`
	Next struct {
		Cursor string `json:"cursor"`
	} `json:"next"`
}

func GetFollowersPaginated(
	fid int,
	appBearerToken string,
	client *http.Client,
	cursor string,
	limit int,
) ([]WarpcastUserInfo, string, error) {
	followers := make([]WarpcastUserInfo, 0, limit)

	query := make(url.Values, 3)
	query.Add("fid", fmt.Sprint(fid))
	query.Add("limit", fmt.Sprint(limit))
	if cursor != "" {
		query.Add("cursor", cursor)
	}

	respBody, err := makeWarpcastRequest(
		"https://api.warpcast.com/v2/followers",
		query,
		appBearerToken,
		client,
	)
	if err != nil {
		return followers, "", err
	}

	// Parsing the response
	var responseStruct WarpcastResponse
	err = json.Unmarshal(respBody, &responseStruct)
	if err != nil {
		return followers, "", err
	}
	followers = responseStruct.Result.Users

	return followers, responseStruct.Next.Cursor, nil
}

func GetFollowingPaginated(
	fid int,
	appBearerToken string,
	client *http.Client,
	cursor string,
	limit int,
) ([]WarpcastUserInfo, string, error) {
	following := make([]WarpcastUserInfo, 0, limit)

	query := make(url.Values, 3)
	query.Add("fid", fmt.Sprint(fid))
	query.Add("limit", fmt.Sprint(limit))
	if cursor != "" {
		query.Add("cursor", cursor)
	}

	respBody, err := makeWarpcastRequest(
		"https://api.warpcast.com/v2/following",
		query,
		appBearerToken,
		client,
	)
	if err != nil {
		return following, "", err
	}

	// Parsing the response
	var responseStruct WarpcastResponse
	err = json.Unmarshal(respBody, &responseStruct)
	if err != nil {
		return following, "", err
	}
	following = responseStruct.Result.Users

	return following, responseStruct.Next.Cursor, nil
}
