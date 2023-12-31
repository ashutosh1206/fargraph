package farcaster

import (
	"encoding/json"
	"fmt"
	"net/url"
)

func (client *FCRequestClient) GetUserCastsPaginated(
	fid int,
	cursor string,
	limit int,
) ([]Cast, string, error) {
	userCasts := make([]Cast, 0, limit)

	query := make(url.Values, 3)
	query.Add("fid", fmt.Sprint(fid))
	query.Add("limit", fmt.Sprint(limit))
	if cursor != "" {
		query.Add("cursor", cursor)
	}

	requestUrl, err := url.JoinPath(client.BaseUrl, "/v2/casts")
	if err != nil {
		return userCasts, "", err
	}

	respBody, err := makeWarpcastRequest(
		requestUrl,
		query,
		client.appBearerToken,
		client.HTTPClient,
	)
	if err != nil {
		return userCasts, "", err
	}

	// Parsing the response
	var responseStruct castsResponse
	err = json.Unmarshal(respBody, &responseStruct)
	if err != nil {
		return userCasts, "", err
	}
	userCasts = responseStruct.Result.Casts

	return userCasts, responseStruct.Next.Cursor, nil
}
