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

func retrieveUserFollowers(
	ctx context.Context,
	driver neo4j.DriverWithContext,
	fid int,
	username string,
	pageLimit int,
	appBearerToken string,
	httpClient *http.Client,
) error {
	followersPaginated, cursor, err := warpcast.GetFollowersPaginated(fid, appBearerToken, httpClient, "", pageLimit)
	if err != nil {
		return err
	}
	err = db.InsertFollowersToDB(followersPaginated, username, fid, ctx, driver)
	if err != nil {
		return err
	}
	for cursor != "" {
		followersPaginated, cursor, err = warpcast.GetFollowersPaginated(fid, appBearerToken, httpClient, cursor, pageLimit)
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
	appBearerToken string,
	httpClient *http.Client,
) error {
	followingPaginated, cursor, err := warpcast.GetFollowingPaginated(fid, appBearerToken, httpClient, "", pageLimit)
	if err != nil {
		return err
	}
	err = db.InsertFollowingToDB(followingPaginated, username, fid, ctx, driver)
	if err != nil {
		return err
	}
	for cursor != "" {
		followingPaginated, cursor, err = warpcast.GetFollowingPaginated(fid, appBearerToken, httpClient, cursor, pageLimit)
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
	appBearerToken string,
	httpClient *http.Client,
) error {
	likedPostsPaginated, cursor, err := warpcast.GetUserLikedCasts(fid, appBearerToken, httpClient, "", pageLimit)
	if err != nil {
		return err
	}
	err = db.InsertUserLikesToDB(likedPostsPaginated, fid, username, ctx, driver)
	if err != nil {
		return err
	}
	for cursor != "" {
		likedPostsPaginated, cursor, err = warpcast.GetUserLikedCasts(fid, appBearerToken, httpClient, cursor, pageLimit)
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
	appBearerToken string,
	httpClient *http.Client,
) error {
	userCastsPaginated, cursor, err := warpcast.GetUserCasts(fid, appBearerToken, httpClient, "", pageLimit)
	if err != nil {
		return err
	}
	err = db.InsertUserPostsToDB(userCastsPaginated, fid, username, ctx, driver)
	if err != nil {
		return err
	}
	for cursor != "" {
		userCastsPaginated, cursor, err = warpcast.GetUserCasts(fid, appBearerToken, httpClient, cursor, pageLimit)
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

	userInfo, err := warpcast.GetUserInfoByUsername(username, appBearerToken, httpClient)
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
	err = retrieveUserFollowers(ctx, driver, fid, username, pageLimit, appBearerToken, httpClient)
	if err != nil {
		panic(err)
	}
	fmt.Println("Inserted followers")

	fmt.Println("Getting following")
	err = retrieveUserFollowing(ctx, driver, fid, username, pageLimit, appBearerToken, httpClient)
	if err != nil {
		panic(err)
	}
	fmt.Println("Inserted following")

	fmt.Println("Getting list of posts liked by user")
	err = retrieveUserLikedCasts(ctx, driver, fid, username, pageLimit, appBearerToken, httpClient)
	if err != nil {
		panic(err)
	}
	fmt.Println("Inserted liked posts")

	fmt.Println("Getting user casts, recasts and replies")
	err = retrieveUserCasts(ctx, driver, fid, username, pageLimit, appBearerToken, httpClient)
	if err != nil {
		panic(err)
	}
	fmt.Println("Inserted user casts, recasts and replies")

	httpClient.CloseIdleConnections()
}
