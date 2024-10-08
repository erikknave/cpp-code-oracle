package server

import (
	"log"

	"github.com/erikknave/go-code-oracle/agents/agenthelpers"
	"github.com/erikknave/go-code-oracle/database"
	"github.com/erikknave/go-code-oracle/dbhelpers"
	"github.com/erikknave/go-code-oracle/search"
	"github.com/erikknave/go-code-oracle/server/chromaclient"
	"github.com/joho/godotenv"
)

func Main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	database.Init()
	dbhelpers.AddInitialUser()
	chromaclient.Init()
	agenthelpers.InitAgentDescriptions()
	search.Init()
	ServerInit()
}
