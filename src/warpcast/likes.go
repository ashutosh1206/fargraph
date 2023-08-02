package warpcast

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type UserCastInfo struct {
	Hash       string `json:"hash"`
	ParentHash string `json:"parentHash"`
	Author     struct {
		Fid      int    `json:"fid"`
		Username string `json:"username"`
	} `json:"author"`
	ParentAuthor struct {
		Fid      int    `json:"fid"`
		Username string `json:"username"`
	} `json:"parentAuthor"`
	Recast bool `json:"recast"`
}

type userLikedCastsResponse struct {
	Result struct {
		Casts []UserCastInfo `json:"casts"`
	} `json:"result"`
	Next struct {
		Cursor string `json:"cursor"`
	} `json:"next"`
}

func GetUserLikedCasts(fid int, appBearerToken string, client *http.Client, cursor string, limit int) ([]UserCastInfo, string, error) {
	userLikedCasts := make([]UserCastInfo, 0, limit)

	// Preparing the request
	request, err := http.NewRequest(
		"GET",
		"https://api.warpcast.com/v2/user-liked-casts",
		nil,
	)
	if err != nil {
		return userLikedCasts, "", err
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
		return userLikedCasts, "", err
	}
	defer response.Body.Close()

	// Validating status code of response
	if response.StatusCode > 400 || response.StatusCode < 200 {
		return userLikedCasts, "", fmt.Errorf("invalid status code: %d", response.StatusCode)
	}

	// Reading response
	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return userLikedCasts, "", err
	}

	// Parsing the response
	var responseStruct userLikedCastsResponse
	err = json.Unmarshal(respBody, &responseStruct)
	if err != nil {
		return userLikedCasts, "", err
	}
	userLikedCasts = responseStruct.Result.Casts

	return userLikedCasts, responseStruct.Next.Cursor, nil
}
