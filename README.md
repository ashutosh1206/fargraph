# Farcaster Graph

Export your entire [Farcaster](https://www.farcaster.xyz/) social graph / activity to a Neo4j Graph Database

## Setting up and running the tool

To setup the tool:

+ Make sure you have a Neo4j instance up and running
  + If you are using Docker, run: `docker run --publish=7474:7474 --publish=7687:7687 --volume=$HOME/neo4j/data:/data --volume=$HOME/neo4j/logs:/logs neo4j`
+ Generate an application bearer token for Farcaster, see docs [here](https://warpcast.notion.site/Warpcast-v2-API-Documentation-c19a9494383a4ce0bd28db6d44d99ea8#c8290028e8f64238bdd2db8938b29b9b)
+ Clone the repository: `git clone https://github.com/ashutosh1206/fargraph` and `cd fargraph`
+ Create .env and update values inside it: `cp example.env .env`
+ You can use `go` to install `fargraph` and `fargraph-binaries` in your `GOPATH`: `GO111MODULE=on go install github.com/ashutosh1206/fargraph/cmd/...`
+ Run the `fargraph-script` binary to setup DB constraints: `fargraph-script`

To run the tool:

Run `./fargraph <username>`, where `username` is the user whose social graph you want to export
