package serverhelpers

import (
	"fmt"
	"strconv"

	"github.com/erikknave/go-code-oracle/dbhelpers"
	"github.com/erikknave/go-code-oracle/search"
	"github.com/erikknave/go-code-oracle/types"
	"github.com/gofiber/fiber/v2"
)

func InterceptAgentType(c *fiber.Ctx, user types.User) (types.UserAgentType, types.SearchableDocument, error) {
	var userAgentType types.UserAgentType
	var searchDoc types.SearchableDocument
	var err error
	dbidStr := c.Query("dbid")
	agentTypeStr := c.Query("agentType")
	if agentTypeStr != "" {
		if dbidStr == "" {
			fmt.Printf("Error: dbid required for agent type: %s\n", agentTypeStr)
			return types.UserAgentType{}, types.SearchableDocument{}, c.SendStatus(400)
		}
		dbid, err := strconv.Atoi(dbidStr)
		if err != nil {
			fmt.Printf("Error converting dbid to int: %v\n", err)
			return types.UserAgentType{}, types.SearchableDocument{}, c.SendStatus(400)
		}
		searchDoc, err := getSearchableDocument(agentTypeStr, dbid)
		if err != nil {
			fmt.Printf("Error getting searchable document: %v\n", err)
			return types.UserAgentType{}, types.SearchableDocument{}, c.SendStatus(400)
		}
		userAgentType = dbhelpers.SetUserAgentType(&user, agentTypeStr, &searchDoc)
		return userAgentType, searchDoc, nil
	} else {
		userAgentType, searchDoc, err = dbhelpers.LoadUserAgentType(&user)
		if err != nil {
			fmt.Printf("Error loading user agent type: %v\n", err)
			return types.UserAgentType{}, types.SearchableDocument{}, c.SendStatus(400)
		}
	}
	return userAgentType, searchDoc, nil
}

func getSearchableDocument(
	agentTypeString string,
	dbid int,
) (types.SearchableDocument, error) {
	switch agentTypeString {
	case "codeBaseAgent":
		return types.SearchableDocument{}, nil
	case "repoAgent":
		searchString := fmt.Sprintf("repository-%d", dbid)
		searchDocs, err := search.SearchReporisitories(searchString, 1)
		if err != nil {
			return types.SearchableDocument{}, err
		}
		if len(searchDocs) == 0 {
			return types.SearchableDocument{}, fmt.Errorf("no search results for repository: %d", dbid)
		}
		return searchDocs[0], nil
	case "packageAgent":
		searchString := fmt.Sprintf("package-%d", dbid)
		searchDocs, err := search.SearchPackages(searchString, "", 1)
		if err != nil {
			return types.SearchableDocument{}, err
		}
		if len(searchDocs) == 0 {
			return types.SearchableDocument{}, fmt.Errorf("no search results for package: %d", dbid)
		}
		return searchDocs[0], nil
	case "fileAgent":
		searchString := fmt.Sprintf("file-%d", dbid)
		searchDocs, err := search.SearchFiles(searchString, "", 1)
		if err != nil {
			return types.SearchableDocument{}, err
		}
		if len(searchDocs) == 0 {
			return types.SearchableDocument{}, fmt.Errorf("no search results for file: %d", dbid)
		}
		return searchDocs[0], nil
	}
	return types.SearchableDocument{}, fmt.Errorf("unknown agent type: %s", agentTypeString)
}
