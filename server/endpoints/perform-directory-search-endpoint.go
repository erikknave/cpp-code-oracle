package endpoints

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/erikknave/go-code-oracle/dbhelpers"
	"github.com/erikknave/go-code-oracle/search"
	"github.com/erikknave/go-code-oracle/types"
	"github.com/erikknave/go-code-oracle/web/templates"
	"github.com/erikknave/go-code-oracle/web/webhelpers"
	"github.com/gofiber/fiber/v2"
)

func PerformDirectorySearchEndPoint(c *fiber.Ctx) error {
	ctx := context.Background()
	user := c.Locals("user").(types.User)
	promptStr := c.FormValue("prompt")
	dbidStr := c.Query("dbid")
	words := strings.Fields(promptStr)
	dbid, err := strconv.Atoi(dbidStr)
	if err != nil {
		log.Fatalf("Error converting dbidStr to int: %v", err)
	}
	if words[0] == "/chat" {
		// messages, err := dbhelpers.LoadChatMessagesForUser(&user)
		// if err != nil {
		// 	fmt.Printf("Error loading chat messages: %v\n", err)
		// 	return c.SendStatus(400)
		// }
		// return webhelpers.RenderHttpComponent(templates.ChatViewWrapper(messages), c, ctx)
	}
	if words[0] == "/search" {
		searchResults, err := dbhelpers.LoadUserSearchResults(&user)
		if err != nil {
			errStr := fmt.Sprintf("Error loading search results: %v", err)
			webhelpers.RenderHttpComponent(templates.SearchViewWrapper(&searchResults, errStr), c, ctx)
		}
		return webhelpers.RenderHttpComponent(templates.SearchViewWrapper(&searchResults, ""), c, ctx)
	}
	if promptStr == "/stats" {
		return StatsViewWrapperEndPoint(c)
	}
	if promptStr == "/search" {
		return PerformSearchEndPoint(c)
	}
	if promptStr == "/help" {
		return HelpViewWrapperEndPoint(c)
	}

	if promptStr == "/chat" {
		dbhelpers.ClearChatMessagesForUser(&user)
		return ChatViewWrapperEndPoint(c)
	}
	if promptStr == "/all" {
		searchResults, err := search.SearchDirectories("", dbidStr, 100)
		if err != nil {
			fmt.Printf("Error performing search: %v", err)
			dbhelpers.SetUserSearchResults(&user, []types.SearchableDocument{})
			return webhelpers.RenderHttpComponent(templates.SearchPackagesContainerWrapper(searchResults, dbid), c, ctx)
		}
		dbhelpers.SetUserSearchResults(&user, searchResults)
		return webhelpers.RenderHttpComponent(templates.SearchPackagesContainerWrapper(searchResults, dbid), c, ctx)
	}
	searchResults, err := search.SearchDirectories(promptStr, dbidStr, 20)
	if err != nil {
		fmt.Printf("Error performing search: %v", err)
		dbhelpers.SetUserSearchResults(&user, []types.SearchableDocument{})
		return webhelpers.RenderHttpComponent(templates.SearchPackagesContainerWrapper(searchResults, dbid), c, ctx)
	}
	dbhelpers.SetUserSearchResults(&user, searchResults)
	return webhelpers.RenderHttpComponent(templates.SearchPackagesContainerWrapper(searchResults, dbid), c, ctx)
	// return c.SendStatus(200)
}
