package db

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func InsertUserNodeToDB(ctx context.Context, driver neo4j.DriverWithContext, fid int, username string) error {
	_, err := neo4j.ExecuteQuery(
		ctx,
		driver,
		"MERGE (u:User {fid: $fid, username: $username})",
		map[string]any{"fid": fid, "username": username},
		neo4j.EagerResultTransformer,
	)
	return err
}

func InsertCastNodeToDB(ctx context.Context, driver neo4j.DriverWithContext, hash string, text string) error {
	if text == "" {
		_, err := neo4j.ExecuteQuery(
			ctx,
			driver,
			"MERGE (c:Cast {hash: $hash})",
			map[string]any{"hash": hash},
			neo4j.EagerResultTransformer,
		)
		return err
	}

	// Set cast text only if it's non-empty

	_, err := neo4j.ExecuteQuery(
		ctx,
		driver,
		"MERGE (c:Cast {hash: $hash}) SET c.text = $text",
		map[string]any{"hash": hash, "text": text},
		neo4j.EagerResultTransformer,
	)
	return err
}

// CreateFollowsEdge creates a directional FOLLOW edge from `fid1` to `fid2`
func CreateFollowsEdge(ctx context.Context, driver neo4j.DriverWithContext, fid1 int, fid2 int) error {
	// There's a UNIQUE constraint on FID, hence it's sufficient to match a node
	// without username
	_, err := neo4j.ExecuteQuery(
		ctx,
		driver,
		"MATCH (u1:User {fid: $fid1}), (u2:User {fid: $fid2}) MERGE (u1)-[r:FOLLOWS]->(u2)",
		map[string]any{"fid1": fid1, "fid2": fid2},
		neo4j.EagerResultTransformer,
	)
	return err
}

func CreatePublishedEdge(ctx context.Context, driver neo4j.DriverWithContext, fid int, hash string) error {
	// There's a UNIQUE constraint on FID, hence it's sufficient to match a node
	// without username
	_, err := neo4j.ExecuteQuery(
		ctx,
		driver,
		"MATCH (c:Cast {hash: $hash}), (u:User {fid: $fid}) MERGE (u)-[r:PUBLISHED]->(c)",
		map[string]any{"hash": hash, "fid": fid},
		neo4j.EagerResultTransformer,
	)
	return err
}

func CreateLikedEdge(ctx context.Context, driver neo4j.DriverWithContext, fid int, hash string) error {
	// There's a UNIQUE constraint on FID, hence it's sufficient to match a node
	// without username
	_, err := neo4j.ExecuteQuery(
		ctx,
		driver,
		"MATCH (c:Cast {hash: $hash}), (u:User {fid: $fid}) MERGE (u)-[r:LIKED]->(c)",
		map[string]any{"hash": hash, "fid": fid},
		neo4j.EagerResultTransformer,
	)
	return err
}

func CreateRecastedEdge(ctx context.Context, driver neo4j.DriverWithContext, fid int, hash string) error {
	// There's a UNIQUE constraint on FID, hence it's sufficient to match a node
	// without username
	_, err := neo4j.ExecuteQuery(
		ctx,
		driver,
		"MATCH (c:Cast {hash: $hash}), (u:User {fid: $fid}) MERGE (u)-[r:RECASTED]->(c)",
		map[string]any{"hash": hash, "fid": fid},
		neo4j.EagerResultTransformer,
	)
	return err
}

func CreateChildOfEdge(ctx context.Context, driver neo4j.DriverWithContext, parentHash string, childHash string) error {
	// There's a UNIQUE constraint on FID, hence it's sufficient to match a node
	// without username
	_, err := neo4j.ExecuteQuery(
		ctx,
		driver,
		"MATCH (pc:Cast {hash: $parentHash}), (cc:Cast {hash: $childHash}) MERGE (cc)-[r:CHILD_OF]->(pc)",
		map[string]any{"parentHash": parentHash, "childHash": childHash},
		neo4j.EagerResultTransformer,
	)
	return err
}
