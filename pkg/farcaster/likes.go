package farcaster

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type Cast struct {
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
	Text   string `json:"text"`
	Recast bool   `json:"recast"`
}

type castsResponse struct {
	Result struct {
		Casts []Cast `json:"casts"`
	} `json:"result"`
	Next struct {
		Cursor string `json:"cursor"`
	} `json:"next"`
}

func (client *FCRequestClient) GetUserLikedPaginated(fid int, cursor string, limit int) ([]Cast, string, error) {
	userLikedCasts := make([]Cast, 0, limit)

	query := make(url.Values, 3)
	query.Add("fid", fmt.Sprint(fid))
	query.Add("limit", fmt.Sprint(limit))
	if cursor != "" {
		query.Add("cursor", cursor)
	}

	requestUrl, err := url.JoinPath(client.BaseUrl, "/v2/user-liked-casts")
	if err != nil {
		return userLikedCasts, "", err
	}

	respBody, err := makeWarpcastRequest(
		requestUrl,
		query,
		client.appBearerToken,
		client.HTTPClient,
	)
	if err != nil {
		return userLikedCasts, "", err
	}

	// Parsing the response
	var responseStruct castsResponse
	err = json.Unmarshal(respBody, &responseStruct)
	if err != nil {
		return userLikedCasts, "", err
	}
	userLikedCasts = responseStruct.Result.Casts

	return userLikedCasts, responseStruct.Next.Cursor, nil
}
