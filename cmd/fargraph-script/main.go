package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

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

	fmt.Println("Creating UNIQUE constraints for Cast hash")
	_, err = neo4j.ExecuteQuery(
		ctx,
		driver,
		"create constraint cast_hash_unique if not exists for (n:Cast) require n.hash is UNIQUE",
		make(map[string]any),
		neo4j.EagerResultTransformer,
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Creating UNIQUE constraints for User fid")
	_, err = neo4j.ExecuteQuery(
		ctx,
		driver,
		"create constraint user_fid_unique if not exists for (n:User) require n.fid is UNIQUE",
		make(map[string]any),
		neo4j.EagerResultTransformer,
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Created all constraints")
}
