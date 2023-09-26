package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/ashutosh1206/fargraph/internal/db"
	"github.com/ashutosh1206/fargraph/pkg/farcaster"
	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/schollz/progressbar/v3"
)

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	appBearerToken := os.Getenv("APP_BEARER_TOKEN")

	// maxFid := 20271
	maxFid := 20

	driver, err := neo4j.NewDriverWithContext(
		os.Getenv("DB_URI"),
		neo4j.BasicAuth(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), ""),
	)
	if err != nil {
		panic(err)
	}
	defer driver.Close(ctx)

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("Database connected!")

	httpClient := http.DefaultClient
	fcRequestClient := farcaster.GetFCRequestClient("https://api.warpcast.com", appBearerToken, httpClient)

	bar := progressbar.Default(int64(maxFid), "Users retrieved")
	defer bar.Close()

	for fid := maxFid; fid > 0; fid-- {
		userInfo, inviterInfo, err := fcRequestClient.GetUserByFid(fid)
		if err != nil {
			panic(err)
		}

		err = db.InsertUserNodeToDB(ctx, driver, userInfo.Fid, userInfo.Username)
		if err != nil {
			panic(err)
		}

		if inviterInfo.Fid != 0 {
			err = db.InsertUserNodeToDB(ctx, driver, inviterInfo.Fid, inviterInfo.Username)
			if err != nil {
				panic(err)
			}

			err = db.CreateInvitedEdge(ctx, driver, userInfo.Fid, inviterInfo.Fid)
			if err != nil {
				panic(err)
			}
		}

		bar.Add(1)
	}
}
