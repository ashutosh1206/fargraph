package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/ashutosh1206/fargraph/internal/db"
	"github.com/ashutosh1206/fargraph/internal/fargraph"
	"github.com/ashutosh1206/fargraph/pkg/farcaster"
	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	if len(os.Args) <= 1 {
		fmt.Println("Usage: fargraph <username>")
		return
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

	username := os.Args[1]

	appBearerToken := os.Getenv("APP_BEARER_TOKEN")

	pageLimit := 100

	httpClient := http.DefaultClient
	fcRequestClient := farcaster.GetFCRequestClient("https://api.warpcast.com", appBearerToken, httpClient)

	userInfo, err := fcRequestClient.GetUserInfoByUsername(username)
	if err != nil {
		panic(err)
	}
	fid := userInfo.Fid

	// Insert source node
	err = db.InsertUserNodeToDB(ctx, driver, fid, username)
	if err != nil {
		panic(err)
	}

	fmt.Println("Getting followers")
	err = fargraph.RetrieveUserFollowers(ctx, driver, fid, pageLimit, fcRequestClient)
	if err != nil {
		panic(err)
	}
	fmt.Println("Inserted followers")

	fmt.Println("Getting following")
	err = fargraph.RetrieveUserFollowing(ctx, driver, fid, pageLimit, fcRequestClient)
	if err != nil {
		panic(err)
	}
	fmt.Println("Inserted following")

	fmt.Println("Getting list of posts liked by user")
	err = fargraph.RetrieveUserLikedCasts(ctx, driver, fid, pageLimit, fcRequestClient)
	if err != nil {
		panic(err)
	}
	fmt.Println("Inserted liked posts")

	fmt.Println("Getting user casts, recasts and replies")
	err = fargraph.RetrieveUserCasts(ctx, driver, fid, pageLimit, fcRequestClient)
	if err != nil {
		panic(err)
	}
	fmt.Println("Inserted user casts, recasts and replies")

	httpClient.CloseIdleConnections()
}
