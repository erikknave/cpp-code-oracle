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

func PerformFileSearchEndPoint(c *fiber.Ctx) error {
	ctx := context.Background()
	user := c.Locals("user").(types.User)
	promptStr := c.FormValue("prompt")
	dbidStr := c.Query("dbid")
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
		searchResults, err := search.SearchFiles("", dbidStr, 100)
		if err != nil {
			fmt.Printf("Error performing search: %v", err)
			// dbhelpers.SetUserSearchResults(&user, []types.SearchableDocument{})
			return webhelpers.RenderHttpComponent(templates.SearchFilesContainerWrapper(searchResults, dbidStr), c, ctx)
		}
		// dbhelpers.SetUserSearchResults(&user, searchResults)
		return webhelpers.RenderHttpComponent(templates.SearchFilesContainerWrapper(searchResults, dbidStr), c, ctx)
	}
	searchResults, err := search.SearchFiles(promptStr, dbidStr, 20)
	if err != nil {
		fmt.Printf("Error performing search: %v", err)
		// dbhelpers.SetUserSearchResults(&user, []types.SearchableDocument{})
		return webhelpers.RenderHttpComponent(templates.SearchFilesContainerWrapper(searchResults, dbidStr), c, ctx)
	}
	// dbhelpers.SetUserSearchResults(&user, searchResults)
	return webhelpers.RenderHttpComponent(templates.SearchFilesContainerWrapper(searchResults, dbidStr), c, ctx)
}
