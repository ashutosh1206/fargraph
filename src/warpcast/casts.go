package warpcast

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetUserCasts(fid int, appBearerToken string, client *http.Client, cursor string, limit int) ([]UserCastInfo, string, error) {
	userCasts := make([]UserCastInfo, 0, limit)

	// Preparing the request
	request, err := http.NewRequest(
		"GET",
		"https://api.warpcast.com/v2/casts",
		nil,
	)
	if err != nil {
		return userCasts, "", err
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
		return userCasts, "", err
	}
	defer response.Body.Close()

	// Validating status code of response
	if response.StatusCode > 300 || response.StatusCode < 100 {
		return userCasts, "", fmt.Errorf("invalid status code: %d", response.StatusCode)
	}

	// Reading response
	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return userCasts, "", err
	}

	// Parsing the response
	var responseStruct userLikedCastsResponse
	err = json.Unmarshal(respBody, &responseStruct)
	if err != nil {
		return userCasts, "", err
	}
	userCasts = responseStruct.Result.Casts

	return userCasts, responseStruct.Next.Cursor, nil
}
