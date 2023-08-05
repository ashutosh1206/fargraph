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

func retrieveUserFollowers(
	ctx context.Context,
	driver neo4j.DriverWithContext,
	fid int,
	username string,
	pageLimit int,
	fcClient *warpcast.FCRequestClient,
) error {
	followersPaginated, cursor, err := fcClient.GetFollowersPaginated(fid, "", pageLimit)
	if err != nil {
		return err
	}
	err = db.InsertFollowersToDB(followersPaginated, username, fid, ctx, driver)
	if err != nil {
		return err
	}
	for cursor != "" {
		followersPaginated, cursor, err = fcClient.GetFollowersPaginated(fid, cursor, pageLimit)
		if err != nil {
			return err
		}
		err = db.InsertFollowersToDB(followersPaginated, username, fid, ctx, driver)
		if err != nil {
			return err
		}
	}
	return nil
}

func retrieveUserFollowing(
	ctx context.Context,
	driver neo4j.DriverWithContext,
	fid int,
	username string,
	pageLimit int,
	fcClient *warpcast.FCRequestClient,
) error {
	followingPaginated, cursor, err := fcClient.GetFollowingPaginated(fid, "", pageLimit)
	if err != nil {
		return err
	}
	err = db.InsertFollowingToDB(followingPaginated, username, fid, ctx, driver)
	if err != nil {
		return err
	}
	for cursor != "" {
		followingPaginated, cursor, err = fcClient.GetFollowingPaginated(fid, cursor, pageLimit)
		if err != nil {
			return err
		}
		err = db.InsertFollowingToDB(followingPaginated, username, fid, ctx, driver)
		if err != nil {
			return err
		}
	}
	return nil
}

func retrieveUserLikedCasts(
	ctx context.Context,
	driver neo4j.DriverWithContext,
	fid int,
	username string,
	pageLimit int,
	fcClient *warpcast.FCRequestClient,
) error {
	likedPostsPaginated, cursor, err := fcClient.GetUserLikedCasts(fid, "", pageLimit)
	if err != nil {
		return err
	}
	err = db.InsertUserLikesToDB(likedPostsPaginated, fid, username, ctx, driver)
	if err != nil {
		return err
	}
	for cursor != "" {
		likedPostsPaginated, cursor, err = fcClient.GetUserLikedCasts(fid, cursor, pageLimit)
		if err != nil {
			return err
		}
		err = db.InsertUserLikesToDB(likedPostsPaginated, fid, username, ctx, driver)
		if err != nil {
			return err
		}
	}
	return nil
}

func retrieveUserCasts(
	ctx context.Context,
	driver neo4j.DriverWithContext,
	fid int,
	username string,
	pageLimit int,
	fcClient *warpcast.FCRequestClient,
) error {
	userCastsPaginated, cursor, err := fcClient.GetUserCasts(fid, "", pageLimit)
	if err != nil {
		return err
	}
	err = db.InsertUserPostsToDB(userCastsPaginated, fid, username, ctx, driver)
	if err != nil {
		return err
	}
	for cursor != "" {
		userCastsPaginated, cursor, err = fcClient.GetUserCasts(fid, cursor, pageLimit)
		if err != nil {
			return err
		}
		err = db.InsertUserPostsToDB(userCastsPaginated, fid, username, ctx, driver)
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

	// TODO: get app bearer token programmatically
	appBearerToken := os.Getenv("APP_BEARER_TOKEN")

	pageLimit := 100

	httpClient := http.DefaultClient
	fcRequestClient := warpcast.GetFCRequestClient("https://api.warpcast.com", appBearerToken, httpClient)

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
	err = retrieveUserFollowers(ctx, driver, fid, username, pageLimit, fcRequestClient)
	if err != nil {
		panic(err)
	}
	fmt.Println("Inserted followers")

	fmt.Println("Getting following")
	err = retrieveUserFollowing(ctx, driver, fid, username, pageLimit, fcRequestClient)
	if err != nil {
		panic(err)
	}
	fmt.Println("Inserted following")

	fmt.Println("Getting list of posts liked by user")
	err = retrieveUserLikedCasts(ctx, driver, fid, username, pageLimit, fcRequestClient)
	if err != nil {
		panic(err)
	}
	fmt.Println("Inserted liked posts")

	fmt.Println("Getting user casts, recasts and replies")
	err = retrieveUserCasts(ctx, driver, fid, username, pageLimit, fcRequestClient)
	if err != nil {
		panic(err)
	}
	fmt.Println("Inserted user casts, recasts and replies")

	httpClient.CloseIdleConnections()
}
