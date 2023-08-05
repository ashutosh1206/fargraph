package fargraph

import (
	"context"

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

func RetrieveUserFollowing(
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

func RetrieveUserLikedCasts(
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

func RetrieveUserCasts(
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
