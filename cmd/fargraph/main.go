package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"

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

	var wg sync.WaitGroup
	ch := make(chan error, 4)

	fmt.Println("Getting followers")
	wg.Add(1)
	go fargraph.RetrieveUserFollowers(ctx, driver, fid, pageLimit, fcRequestClient, ch, &wg)

	fmt.Println("Getting following")
	wg.Add(1)
	go fargraph.RetrieveUserFollowing(ctx, driver, fid, pageLimit, fcRequestClient, ch, &wg)

	fmt.Println("Getting list of posts liked by user")
	wg.Add(1)
	go fargraph.RetrieveUserLikedCasts(ctx, driver, fid, pageLimit, fcRequestClient, ch, &wg)

	fmt.Println("Getting user casts, recasts and replies")
	wg.Add(1)
	go fargraph.RetrieveUserCasts(ctx, driver, fid, pageLimit, fcRequestClient, ch, &wg)

	go func() {
		wg.Wait()
		close(ch)
		httpClient.CloseIdleConnections()
	}()

	for err := range ch {
		if err != nil {
			panic(err)
		}
	}
}
