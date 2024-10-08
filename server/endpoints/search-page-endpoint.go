package endpoints

import (
	"context"
	"fmt"
	"strings"

	"github.com/erikknave/go-code-oracle/dbhelpers"
	"github.com/erikknave/go-code-oracle/server/serverhelpers"
	"github.com/erikknave/go-code-oracle/types"
	"github.com/erikknave/go-code-oracle/web/templates"
	"github.com/erikknave/go-code-oracle/web/webhelpers"
	"github.com/gofiber/fiber/v2"
)

func SearchPageEndPoint(c *fiber.Ctx) error {
	ctx := context.Background()
	userInterface := c.Locals("user")
	user := userInterface.(types.User)
	promptStr := c.FormValue("prompt")
	if promptStr == "" {
		searchResults, err := dbhelpers.LoadUserSearchResults(&user)
		if err != nil {
			errStr := fmt.Sprintf("Error loading search results: %v", err)
			webhelpers.RenderHttpComponent(templates.SearchPage(&searchResults, errStr), c, ctx)
		}
		return webhelpers.RenderHttpComponent(templates.SearchPage(&searchResults, ""), c, ctx)
	}
	words := strings.Fields(promptStr)
	if words[0] == "/stats" {
		return StatsViewWrapperEndPoint(c)
	}
	if words[0] == "/chat" {
		messages, err := dbhelpers.LoadChatMessagesForUser(&user)
		if err != nil {
			fmt.Printf("Error loading chat messages: %v\n", err)
			return c.SendStatus(400)
		}
		agentType, searchableDocument, err := serverhelpers.InterceptAgentType(c, user)
		if err != nil {
			fmt.Printf("Error intercepting agent type: %v\n", err)
			return c.SendStatus(400)
		}
		return webhelpers.RenderHttpComponent(templates.ChatViewWrapper(messages, agentType, searchableDocument), c, ctx)
	}
	if words[0] == "/search" {
		searchResults, err := dbhelpers.LoadUserSearchResults(&user)
		if err != nil {
			errStr := fmt.Sprintf("Error loading search results: %v", err)
			webhelpers.RenderHttpComponent(templates.SearchViewWrapper(&searchResults, errStr), c, ctx)
		}
		return webhelpers.RenderHttpComponent(templates.SearchViewWrapper(&searchResults, ""), c, ctx)
	}
	searchResults, err := serverhelpers.PerformSearch(promptStr)
	if err != nil {
		errStr := fmt.Sprintf("Error performing search: %v", err)
		dbhelpers.SetUserSearchResults(&user, []types.SearchableDocument{})
		return webhelpers.RenderHttpComponent(templates.SearchViewWrapper(&searchResults, errStr), c, ctx)
	}
	dbhelpers.SetUserSearchResults(&user, searchResults)
	webhelpers.RenderHttpComponent(templates.SearchViewWrapper(&searchResults, ""), c, ctx)
	return nil
}
