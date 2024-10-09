package endpoints

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/erikknave/go-code-oracle/dbhelpers"
	"github.com/erikknave/go-code-oracle/search"
	"github.com/erikknave/go-code-oracle/types"
	"github.com/erikknave/go-code-oracle/web/templates"
	"github.com/erikknave/go-code-oracle/web/webhelpers"
	"github.com/gofiber/fiber/v2"
)

func PerformCodeblockSearchEndPoint(c *fiber.Ctx) error {
	ctx := context.Background()
	promptStr := c.FormValue("prompt")
	dbidStr := c.Query("dbid")
	dbid := convertDbid(dbidStr)
	userInterface := c.Locals("user")
	user := userInterface.(types.User)
	if promptStr == "/help" {
		return HelpViewWrapperEndPoint(c)
	}
	if promptStr == "/stats" {
		return StatsViewWrapperEndPoint(c)
	}
	if promptStr == "/search" {
		return PerformSearchEndPoint(c)
	}
	if promptStr == "/chat" {
		dbhelpers.ClearChatMessagesForUser(&user)
		return ChatViewWrapperEndPoint(c)
	}
	if promptStr == "/all" {
		searchResults, err := search.SearchCodeblocks("", dbidStr, 50)
		if err != nil {
			fmt.Printf("Error performing search: %v", err)

			return webhelpers.RenderHttpComponent(templates.SearchEntitiesContainerWrapper(searchResults, dbid), c, ctx)
		}
		return webhelpers.RenderHttpComponent(templates.SearchEntitiesContainerWrapper(searchResults, dbid), c, ctx)
	}
	searchResults, err := search.SearchCodeblocks(promptStr, dbidStr, 20)
	if err != nil {
		fmt.Printf("Error performing search: %v", err)
		return webhelpers.RenderHttpComponent(templates.SearchEntitiesContainerWrapper(searchResults, dbid), c, ctx)
	}
	return webhelpers.RenderHttpComponent(templates.SearchEntitiesContainerWrapper(searchResults, dbid), c, ctx)
}

func convertDbid(dbidStr string) int {
	dbid, err := strconv.Atoi(dbidStr)
	if err != nil {
		log.Fatalf("Error converting dbidStr to int: %v", err)
	}
	return dbid
}
