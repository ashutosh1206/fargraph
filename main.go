package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/ashutosh1206/fargraph/pkg/farcaster"
	"github.com/ashutosh1206/fargraph/src/db"
	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func retrieveUserFollowers(
	ctx context.Context,
	driver neo4j.DriverWithContext,
	fid int,
	pageLimit int,
	fcClient *farcaster.FCRequestClient,
) error {
	followersPaginated, cursor, err := fcClient.GetFollowersPaginated(fid, "", pageLimit)
	if err != nil {
		return err
	}
	err = db.InsertFollowersToDB(ctx, driver, followersPaginated, fid)
	if err != nil {
		return err
	}
	for cursor != "" {
		followersPaginated, cursor, err = fcClient.GetFollowersPaginated(fid, cursor, pageLimit)
		if err != nil {
			return err
		}
		err = db.InsertFollowersToDB(ctx, driver, followersPaginated, fid)
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
	pageLimit int,
	fcClient *farcaster.FCRequestClient,
) error {
	followingPaginated, cursor, err := fcClient.GetFollowingPaginated(fid, "", pageLimit)
	if err != nil {
		return err
	}
	err = db.InsertFollowingToDB(ctx, driver, followingPaginated, fid)
	if err != nil {
		return err
	}
	for cursor != "" {
		followingPaginated, cursor, err = fcClient.GetFollowingPaginated(fid, cursor, pageLimit)
		if err != nil {
			return err
		}
		err = db.InsertFollowingToDB(ctx, driver, followingPaginated, fid)
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
	pageLimit int,
	fcClient *farcaster.FCRequestClient,
) error {
	likedCastsPaginated, cursor, err := fcClient.GetUserLikedPaginated(fid, "", pageLimit)
	if err != nil {
		return err
	}
	err = db.InsertUserLikesToDB(ctx, driver, likedCastsPaginated, fid)
	if err != nil {
		return err
	}
	for cursor != "" {
		likedCastsPaginated, cursor, err = fcClient.GetUserLikedPaginated(fid, cursor, pageLimit)
		if err != nil {
			return err
		}
		err = db.InsertUserLikesToDB(ctx, driver, likedCastsPaginated, fid)
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
	pageLimit int,
	fcClient *farcaster.FCRequestClient,
) error {
	userCastsPaginated, cursor, err := fcClient.GetUserCastsPaginated(fid, "", pageLimit)
	if err != nil {
		return err
	}
	err = db.InsertUserPostsToDB(ctx, driver, userCastsPaginated, fid)
	if err != nil {
		return err
	}
	for cursor != "" {
		userCastsPaginated, cursor, err = fcClient.GetUserCastsPaginated(fid, cursor, pageLimit)
		if err != nil {
			return err
		}
		err = db.InsertUserPostsToDB(ctx, driver, userCastsPaginated, fid)
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
	err = retrieveUserFollowers(ctx, driver, fid, pageLimit, fcRequestClient)
	if err != nil {
		panic(err)
	}
	fmt.Println("Inserted followers")

	fmt.Println("Getting following")
	err = retrieveUserFollowing(ctx, driver, fid, pageLimit, fcRequestClient)
	if err != nil {
		panic(err)
	}
	fmt.Println("Inserted following")

	fmt.Println("Getting list of posts liked by user")
	err = retrieveUserLikedCasts(ctx, driver, fid, pageLimit, fcRequestClient)
	if err != nil {
		panic(err)
	}
	fmt.Println("Inserted liked posts")

	fmt.Println("Getting user casts, recasts and replies")
	err = retrieveUserCasts(ctx, driver, fid, pageLimit, fcRequestClient)
	if err != nil {
		panic(err)
	}
	fmt.Println("Inserted user casts, recasts and replies")

	httpClient.CloseIdleConnections()
}
