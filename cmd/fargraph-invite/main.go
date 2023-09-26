package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"

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
	maxFid := 20271
	MAX_CONCURRENT_JOBS := 100

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

	var wg sync.WaitGroup
	ch := make(chan error)
	waitCh := make(chan struct{}, MAX_CONCURRENT_JOBS)

	for fid := maxFid; fid > 0; fid-- {
		waitCh <- struct{}{}
		wg.Add(1)

		go func(fid int) {
			defer wg.Done()

			userInfo, inviterInfo, err := fcRequestClient.GetUserByFid(fid)
			if err != nil {
				ch <- err
			}

			err = db.InsertUserNodeToDB(ctx, driver, userInfo.Fid, userInfo.Username)
			if err != nil {
				ch <- err
			}

			if inviterInfo.Fid != 0 {
				err = db.InsertUserNodeToDB(ctx, driver, inviterInfo.Fid, inviterInfo.Username)
				if err != nil {
					ch <- err
				}

				err = db.CreateInvitedEdge(ctx, driver, userInfo.Fid, inviterInfo.Fid)
				if err != nil {
					ch <- err
				}
			}

			bar.Add(1)
			<-waitCh
		}(fid)
	}

	go func() {
		wg.Wait()
		close(ch)
		close(waitCh)
		httpClient.CloseIdleConnections()
	}()

	for err := range ch {
		if err != nil {
			panic(err)
		}
	}
}
