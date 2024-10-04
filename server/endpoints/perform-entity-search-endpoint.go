package endpoints

import (
	"context"
	"fmt"

	"github.com/erikknave/go-code-oracle/dbhelpers"
	"github.com/erikknave/go-code-oracle/search"
	"github.com/erikknave/go-code-oracle/types"
	"github.com/erikknave/go-code-oracle/web/templates"
	"github.com/erikknave/go-code-oracle/web/webhelpers"
	"github.com/gofiber/fiber/v2"
)

func PerformEntitySearchEndPoint(c *fiber.Ctx) error {
	ctx := context.Background()
	promptStr := c.FormValue("prompt")
	dbidStr := c.Query("dbid")
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
		searchResults, err := search.SearchEntities("", dbidStr, 50)
		if err != nil {
			fmt.Printf("Error performing search: %v", err)
			return webhelpers.RenderHttpComponent(templates.SearchEntitiesContainerWrapper(searchResults, dbidStr), c, ctx)
		}
		return webhelpers.RenderHttpComponent(templates.SearchEntitiesContainerWrapper(searchResults, dbidStr), c, ctx)
	}
	searchResults, err := search.SearchEntities(promptStr, dbidStr, 20)
	if err != nil {
		fmt.Printf("Error performing search: %v", err)
		return webhelpers.RenderHttpComponent(templates.SearchEntitiesContainerWrapper(searchResults, dbidStr), c, ctx)
	}
	return webhelpers.RenderHttpComponent(templates.SearchEntitiesContainerWrapper(searchResults, dbidStr), c, ctx)
}
