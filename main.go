package main

import (
	"fmt"

	"github.com/erikknave/go-code-oracle/agents/agenthelpers"
	"github.com/erikknave/go-code-oracle/database"
	"github.com/erikknave/go-code-oracle/dbhelpers"
	"github.com/erikknave/go-code-oracle/search"
	"github.com/erikknave/go-code-oracle/server"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	// if err != nil {
	// 	log.Fatalf("Error loading .env file")
	// }
	database.Init()
	dbhelpers.AddInitialUser()
	// chromaclient.Init()
	agenthelpers.InitAgentDescriptions()
	search.Init()
	fmt.Println("Server starting...")
	server.ServerInit()
	// repos, _ := cypherqueries.PerformRepoListCypherQuery()
	// helpers.PrettyPrintJSONInterface(repos)
}
