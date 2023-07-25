package warpcast

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

	// Preparing the request
	request, err := http.NewRequest(
		"GET",
		"https://api.warpcast.com/v2/followers",
		nil,
	)
	if err != nil {
		return followers, "", err
	}
	query := request.URL.Query()
	query.Add("fid", fmt.Sprint(fid))
	query.Add("limit", fmt.Sprint(limit))
	if cursor != "" {
		query.Add("cursor", cursor)
	}
	request.URL.RawQuery = query.Encode()
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", appBearerToken))

	// Sending the request
	response, err := client.Do(request)
	if err != nil {
		return followers, "", err
	}
	defer response.Body.Close()

	// Validating status code of response
	if response.StatusCode > 300 || response.StatusCode < 100 {
		return followers, "", fmt.Errorf("invalid status code: %d", response.StatusCode)
	}

	// Reading response
	respBody, err := io.ReadAll(response.Body)
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

	// Preparing the request
	request, err := http.NewRequest(
		"GET",
		"https://api.warpcast.com/v2/following",
		nil,
	)
	if err != nil {
		return following, "", err
	}
	query := request.URL.Query()
	query.Add("fid", fmt.Sprint(fid))
	query.Add("limit", fmt.Sprint(limit))
	if cursor != "" {
		query.Add("cursor", cursor)
	}
	request.URL.RawQuery = query.Encode()
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", appBearerToken))

	// Sending the request
	response, err := client.Do(request)
	if err != nil {
		return following, "", err
	}
	defer response.Body.Close()

	// Validating status code of response
	if response.StatusCode > 300 || response.StatusCode < 100 {
		return following, "", fmt.Errorf("invalid status code: %d", response.StatusCode)
	}

	// Reading response
	respBody, err := io.ReadAll(response.Body)
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
