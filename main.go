package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/farcaster-graph/src/db"
	"github.com/farcaster-graph/src/warpcast"
	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func connectToDB(
	dbUri string,
	dbUsername string,
	dbPassword string,
	ctx context.Context,
) (neo4j.DriverWithContext, error) {
	var driver neo4j.DriverWithContext

	return driver, nil
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
	followersPaginated, cursor, err := warpcast.GetFollowersPaginated(fid, appBearerToken, httpClient, "", pageLimit)
	if err != nil {
		panic(err)
	}
	err = db.InsertFollowersToDB(followersPaginated, username, fid, ctx, driver)
	if err != nil {
		panic(err)
	}
	for cursor != "" {
		followersPaginated, cursor, err = warpcast.GetFollowersPaginated(fid, appBearerToken, httpClient, cursor, pageLimit)
		if err != nil {
			panic(err)
		}
		err = db.InsertFollowersToDB(followersPaginated, username, fid, ctx, driver)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Inserted followers")

	fmt.Println("Getting following")
	followingPaginated, cursor, err := warpcast.GetFollowingPaginated(fid, appBearerToken, httpClient, "", pageLimit)
	if err != nil {
		panic(err)
	}
	err = db.InsertFollowingToDB(followingPaginated, username, fid, ctx, driver)
	if err != nil {
		panic(err)
	}
	for cursor != "" {
		followingPaginated, cursor, err = warpcast.GetFollowingPaginated(fid, appBearerToken, httpClient, cursor, pageLimit)
		if err != nil {
			panic(err)
		}
		err = db.InsertFollowingToDB(followingPaginated, username, fid, ctx, driver)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Inserted following")

	httpClient.CloseIdleConnections()
}
