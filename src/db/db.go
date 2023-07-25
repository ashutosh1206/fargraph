package db

import (
	"context"

	"github.com/farcaster-graph/src/warpcast"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func InsertFollowersToDB(followers []warpcast.WarpcastUserInfo, username string, fid int, ctx context.Context, driver neo4j.DriverWithContext) error {
	for _, follower := range followers {
		_, err := neo4j.ExecuteQuery(
			ctx,
			driver,
			"MERGE (u:User {fid: $fid, username: $username})",
			map[string]any{"fid": follower.Fid, "username": follower.Username},
			neo4j.EagerResultTransformer,
		)
		if err != nil {
			return err
		}
		_, err = neo4j.ExecuteQuery(
			ctx,
			driver,
			"MATCH (u1:User {fid: $fid1, username: $username1}), (u2:User {fid: $fid2, username: $username2}) MERGE (u1)-[r:FOLLOWS]->(u2)",
			map[string]any{"fid1": follower.Fid, "username1": follower.Username, "fid2": fid, "username2": username},
			neo4j.EagerResultTransformer,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func InsertFollowingToDB(following []warpcast.WarpcastUserInfo, username string, fid int, ctx context.Context, driver neo4j.DriverWithContext) error {
	for _, followee := range following {
		_, err := neo4j.ExecuteQuery(
			ctx,
			driver,
			"MERGE (u:User {fid: $fid, username: $username})",
			map[string]any{"fid": followee.Fid, "username": followee.Username},
			neo4j.EagerResultTransformer,
		)
		if err != nil {
			return err
		}
		_, err = neo4j.ExecuteQuery(
			ctx,
			driver,
			"MATCH (u1:User {fid: $fid1, username: $username1}), (u2:User {fid: $fid2, username: $username2}) MERGE (u1)-[r:FOLLOWS]->(u2)",
			map[string]any{"fid2": followee.Fid, "username2": followee.Username, "fid1": fid, "username1": username},
			neo4j.EagerResultTransformer,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
