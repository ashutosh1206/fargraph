package fargraph

import (
	"context"
	"sync"

	"github.com/ashutosh1206/fargraph/internal/db"
	"github.com/ashutosh1206/fargraph/pkg/farcaster"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func RetrieveUserFollowers(
	ctx context.Context,
	driver neo4j.DriverWithContext,
	fid int,
	pageLimit int,
	fcClient *farcaster.FCRequestClient,
	ch chan error,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	followersPaginated, cursor, err := fcClient.GetFollowersPaginated(fid, "", pageLimit)
	if err != nil {
		ch <- err
	}
	err = db.InsertFollowersToDB(ctx, driver, followersPaginated, fid)
	if err != nil {
		ch <- err
	}
	for cursor != "" {
		followersPaginated, cursor, err = fcClient.GetFollowersPaginated(fid, cursor, pageLimit)
		if err != nil {
			ch <- err
		}
		err = db.InsertFollowersToDB(ctx, driver, followersPaginated, fid)
		if err != nil {
			ch <- err
		}
	}
	ch <- nil
}

func RetrieveUserFollowing(
	ctx context.Context,
	driver neo4j.DriverWithContext,
	fid int,
	pageLimit int,
	fcClient *farcaster.FCRequestClient,
	ch chan error,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	followingPaginated, cursor, err := fcClient.GetFollowingPaginated(fid, "", pageLimit)
	if err != nil {
		ch <- err
	}
	err = db.InsertFollowingToDB(ctx, driver, followingPaginated, fid)
	if err != nil {
		ch <- err
	}
	for cursor != "" {
		followingPaginated, cursor, err = fcClient.GetFollowingPaginated(fid, cursor, pageLimit)
		if err != nil {
			ch <- err
		}
		err = db.InsertFollowingToDB(ctx, driver, followingPaginated, fid)
		if err != nil {
			ch <- err
		}
	}
	ch <- nil
}

func RetrieveUserLikedCasts(
	ctx context.Context,
	driver neo4j.DriverWithContext,
	fid int,
	pageLimit int,
	fcClient *farcaster.FCRequestClient,
	ch chan error,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	likedCastsPaginated, cursor, err := fcClient.GetUserLikedPaginated(fid, "", pageLimit)
	if err != nil {
		ch <- err
	}
	err = db.InsertUserLikesToDB(ctx, driver, likedCastsPaginated, fid)
	if err != nil {
		ch <- err
	}
	for cursor != "" {
		likedCastsPaginated, cursor, err = fcClient.GetUserLikedPaginated(fid, cursor, pageLimit)
		if err != nil {
			ch <- err
		}
		err = db.InsertUserLikesToDB(ctx, driver, likedCastsPaginated, fid)
		if err != nil {
			ch <- err
		}
	}
	ch <- nil
}

func RetrieveUserCasts(
	ctx context.Context,
	driver neo4j.DriverWithContext,
	fid int,
	pageLimit int,
	fcClient *farcaster.FCRequestClient,
	ch chan error,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	userCastsPaginated, cursor, err := fcClient.GetUserCastsPaginated(fid, "", pageLimit)
	if err != nil {
		ch <- err
	}
	err = db.InsertUserPostsToDB(ctx, driver, userCastsPaginated, fid)
	if err != nil {
		ch <- err
	}
	for cursor != "" {
		userCastsPaginated, cursor, err = fcClient.GetUserCastsPaginated(fid, cursor, pageLimit)
		if err != nil {
			ch <- err
		}
		err = db.InsertUserPostsToDB(ctx, driver, userCastsPaginated, fid)
		if err != nil {
			ch <- err
		}
	}
	ch <- nil
}
