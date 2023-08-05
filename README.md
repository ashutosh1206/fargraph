# Farcaster Graph

Export your entire Farcaster social graph / activity to a Neo4j Graph Database

## Setting up and running the app

To setup the tool:

+ Make sure you have a Neo4j instance up and running
  + If you are using Docker, run: `docker run --publish=7474:7474 --publish=7687:7687 --volume=$HOME/neo4j/data:/data --volume=$HOME/neo4j/logs:/logs neo4j`
+ Generate an application bearer token for Farcaster, see docs [here](https://warpcast.notion.site/Warpcast-v2-API-Documentation-c19a9494383a4ce0bd28db6d44d99ea8#c8290028e8f64238bdd2db8938b29b9b)
+ Create .env and update values inside it: `cp example.env .env`
+ Compile the binary by running: `go build .`, assuming you are in the root directory of the project

To run the tool:
`./farcaster <username>`, where `username` is the user whose social graph you want to export
