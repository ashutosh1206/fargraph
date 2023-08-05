package warpcast

import (
	"encoding/json"
	"fmt"
	"net/url"
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

func (client *FCRequestClient) GetUserLikedCasts(fid int, cursor string, limit int) ([]UserCastInfo, string, error) {
	userLikedCasts := make([]UserCastInfo, 0, limit)

	query := make(url.Values, 3)
	query.Add("fid", fmt.Sprint(fid))
	query.Add("limit", fmt.Sprint(limit))
	if cursor != "" {
		query.Add("cursor", cursor)
	}

	respBody, err := makeWarpcastRequest(
		"https://api.warpcast.com/v2/user-liked-casts",
		query,
		client.appBearerToken,
		client.HTTPClient,
	)
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
