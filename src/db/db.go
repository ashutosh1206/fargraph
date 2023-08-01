package db

import (
	"context"

	"github.com/farcaster-graph/src/warpcast"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// TODO: optimise insertions with bulk-insertions

func InsertFollowersToDB(followers []warpcast.WarpcastUserInfo, username string, fid int, ctx context.Context, driver neo4j.DriverWithContext) error {
	for _, follower := range followers {
		err := InsertUserNodeToDB(ctx, driver, follower.Fid, follower.Username)
		if err != nil {
			return err
		}
		err = CreateFollowsEdge(ctx, driver, follower.Fid, fid)
		if err != nil {
			return err
		}
	}
	return nil
}

func InsertFollowingToDB(following []warpcast.WarpcastUserInfo, username string, fid int, ctx context.Context, driver neo4j.DriverWithContext) error {
	for _, followee := range following {
		err := InsertUserNodeToDB(ctx, driver, followee.Fid, followee.Username)
		if err != nil {
			return err
		}
		err = CreateFollowsEdge(ctx, driver, fid, followee.Fid)
		if err != nil {
			return err
		}
	}
	return nil
}

func InsertUserLikesToDB(likedCasts []warpcast.UserCastInfo, fid int, username string, ctx context.Context, driver neo4j.DriverWithContext) error {
	for _, cast := range likedCasts {
		// Create User node
		err := InsertUserNodeToDB(ctx, driver, cast.Author.Fid, cast.Author.Username)
		if err != nil {
			return err
		}
		// Create Cast node
		err = InsertCastNodeToDB(ctx, driver, cast.Hash)
		if err != nil {
			return err
		}
		// Create a PUBLISHED edge between User (author) and Cast
		err = CreatePublishedEdge(ctx, driver, cast.Author.Fid, cast.Hash)
		if err != nil {
			return err
		}
		// Create a LIKED edge between User and Cast
		err = CreateLikedEdge(ctx, driver, fid, cast.Hash)
		if err != nil {
			return err
		}
	}
	return nil
}

func InsertUserPostsToDB(casts []warpcast.UserCastInfo, fid int, username string, ctx context.Context, driver neo4j.DriverWithContext) error {
	for _, cast := range casts {
		// Create User node for Author
		err := InsertUserNodeToDB(ctx, driver, cast.Author.Fid, cast.Author.Username)
		if err != nil {
			return err
		}

		// Create Cast node
		err = InsertCastNodeToDB(ctx, driver, cast.Hash)
		if err != nil {
			return err
		}

		// Create a PUBLISHED edge between User (author) and Cast
		err = CreatePublishedEdge(ctx, driver, cast.Author.Fid, cast.Hash)
		if err != nil {
			return err
		}

		if cast.Recast {
			// Create a RECASTED edge between fid and Cast
			err = CreateRecastedEdge(ctx, driver, fid, cast.Hash)
			if err != nil {
				return err
			}
			continue
		}

		// Not a recast

		// Simple cast (non-reply) by fid,username
		if cast.ParentAuthor.Fid == 0 || cast.ParentAuthor.Username == "" || cast.ParentHash == "" {
			continue
		}

		// Reply cast

		// Create User node for Parent author
		err = InsertUserNodeToDB(ctx, driver, cast.ParentAuthor.Fid, cast.ParentAuthor.Username)
		if err != nil {
			return err
		}

		// Create Cast node for parent cast
		err = InsertCastNodeToDB(ctx, driver, cast.ParentHash)
		if err != nil {
			return err
		}

		// Create a PUBLISHED edge between User node (Parent author)
		// and Cast node (Parent cast)
		err = CreatePublishedEdge(ctx, driver, cast.ParentAuthor.Fid, cast.ParentHash)
		if err != nil {
			return err
		}

		// Create a CHILD_OF edge between Cast node (Parent cast) and Cast node
		err = CreateChildOfEdge(ctx, driver, cast.ParentHash, cast.Hash)
		if err != nil {
			return err
		}
	}
	return nil
}
