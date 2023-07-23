package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type warpcastUserInfo struct {
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

type warpcastResponse struct {
	Result struct {
		Users []warpcastUserInfo `json:"users"`
	} `json:"result"`
	Next struct {
		Cursor string `json:"cursor"`
	} `json:"next"`
}

func getFollowersPaginated(
	fid int,
	appBearerToken string,
	client *http.Client,
	cursor string,
	limit int,
) ([]warpcastUserInfo, string, error) {
	followers := make([]warpcastUserInfo, 0, limit)

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
	var responseStruct warpcastResponse
	err = json.Unmarshal(respBody, &responseStruct)
	if err != nil {
		return followers, "", err
	}
	followers = responseStruct.Result.Users

	return followers, responseStruct.Next.Cursor, nil
}

func getFollowingPaginated(
	fid int,
	appBearerToken string,
	client *http.Client,
	cursor string,
	limit int,
) ([]warpcastUserInfo, string, error) {
	following := make([]warpcastUserInfo, 0, limit)

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
	var responseStruct warpcastResponse
	err = json.Unmarshal(respBody, &responseStruct)
	if err != nil {
		return following, "", err
	}
	following = responseStruct.Result.Users

	return following, responseStruct.Next.Cursor, nil
}

func connectToDB(
	dbUri string,
	dbUsername string,
	dbPassword string,
	ctx context.Context,
) (neo4j.DriverWithContext, error) {
	var driver neo4j.DriverWithContext

	return driver, nil
}

func insertFollowersToDB(followers []warpcastUserInfo, username string, fid int, ctx context.Context, driver neo4j.DriverWithContext) error {
	for _, follower := range followers {
		_, err := neo4j.ExecuteQuery(
			ctx,
			driver,
			"MERGE (u:User {fid: $fid, username: $username})",
			map[string]any{"fid": follower.Fid, "username": follower.Username},
			neo4j.EagerResultTransformer,
		)
		if err != nil {
			return err
		}
		_, err = neo4j.ExecuteQuery(
			ctx,
			driver,
			"MATCH (u1:User {fid: $fid1, username: $username1}), (u2:User {fid: $fid2, username: $username2}) MERGE (u1)-[r:FOLLOWS]->(u2)",
			map[string]any{"fid1": follower.Fid, "username1": follower.Username, "fid2": fid, "username2": username},
			neo4j.EagerResultTransformer,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func insertFollowingToDB(following []warpcastUserInfo, username string, fid int, ctx context.Context, driver neo4j.DriverWithContext) error {
	for _, followee := range following {
		_, err := neo4j.ExecuteQuery(
			ctx,
			driver,
			"MERGE (u:User {fid: $fid, username: $username})",
			map[string]any{"fid": followee.Fid, "username": followee.Username},
			neo4j.EagerResultTransformer,
		)
		if err != nil {
			return err
		}
		_, err = neo4j.ExecuteQuery(
			ctx,
			driver,
			"MATCH (u1:User {fid: $fid1, username: $username1}), (u2:User {fid: $fid2, username: $username2}) MERGE (u1)-[r:FOLLOWS]->(u2)",
			map[string]any{"fid2": followee.Fid, "username2": followee.Username, "fid1": fid, "username1": username},
			neo4j.EagerResultTransformer,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	driver, err := neo4j.NewDriverWithContext(
		os.Getenv("DB_URI"),
		neo4j.BasicAuth(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), ""))

	if err != nil {
		panic(err)
	}
	defer driver.Close(ctx)

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("Database connected!")

	// TODO: user input
	username := "ashutosh"

	// TODO: get fid programmatically
	fid := 5267

	// TODO: get app bearer token programmatically
	appBearerToken := os.Getenv("APP_BEARER_TOKEN")

	pageLimit := 100

	httpClient := http.DefaultClient

	// Insert source node
	_, err = neo4j.ExecuteQuery(
		ctx,
		driver,
		"MERGE (u:User {fid: $fid, username: $username})",
		map[string]any{"fid": fid, "username": username},
		neo4j.EagerResultTransformer,
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Getting followers")
	followersPaginated, cursor, err := getFollowersPaginated(fid, appBearerToken, httpClient, "", pageLimit)
	if err != nil {
		panic(err)
	}
	err = insertFollowersToDB(followersPaginated, username, fid, ctx, driver)
	if err != nil {
		panic(err)
	}
	for cursor != "" {
		followersPaginated, cursor, err = getFollowersPaginated(fid, appBearerToken, httpClient, cursor, pageLimit)
		if err != nil {
			panic(err)
		}
		err = insertFollowersToDB(followersPaginated, username, fid, ctx, driver)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Inserted followers")

	fmt.Println("Getting following")
	followingPaginated, cursor, err := getFollowingPaginated(fid, appBearerToken, httpClient, "", pageLimit)
	if err != nil {
		panic(err)
	}
	err = insertFollowingToDB(followingPaginated, username, fid, ctx, driver)
	if err != nil {
		panic(err)
	}
	for cursor != "" {
		followingPaginated, cursor, err = getFollowingPaginated(fid, appBearerToken, httpClient, cursor, pageLimit)
		if err != nil {
			panic(err)
		}
		err = insertFollowingToDB(followingPaginated, username, fid, ctx, driver)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Inserted following")

	httpClient.CloseIdleConnections()
}
