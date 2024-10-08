package endpoints

import (
	"log"

	"github.com/erikknave/go-code-oracle/dbhelpers"
	"github.com/erikknave/go-code-oracle/server/serverhelpers"
	"github.com/erikknave/go-code-oracle/types"
	"github.com/erikknave/go-code-oracle/web/templates"
	"github.com/erikknave/go-code-oracle/web/webhelpers"
	"github.com/gofiber/fiber/v2"
)

func ChatPageEndPoint(c *fiber.Ctx) error {
	ctx := c.Context()
	var user types.User
	var err error
	userInterface := c.Locals("user")
	user = userInterface.(types.User)
	viewName := c.Query("view")
	if viewName == "search" {
		searchResults, err := dbhelpers.LoadUserSearchResults(&user)
		if err != nil {
			log.Fatalf("Error loading search results: %v", err)
		}
		return webhelpers.RenderHttpComponent(templates.SearchPage(&searchResults, ""), c, ctx)

	}
	chatMessages, _ := dbhelpers.LoadChatMessagesForUser(&user)
	agentType, searchableDocument, err := serverhelpers.InterceptAgentType(c, user)
	if err != nil {
		log.Printf("Error loading agent type: %v", err)
		return c.SendStatus(400)
	}
	chatView := templates.ChatPage(chatMessages, agentType, searchableDocument)
	return webhelpers.RenderHttpComponent(chatView, c, ctx)
}

func ChatViewWrapperEndPoint(c *fiber.Ctx) error {
	ctx := c.Context()
	var user types.User
	var err error
	userInterface := c.Locals("user")
	user = userInterface.(types.User)
	chatMessages, _ := dbhelpers.LoadChatMessagesForUser(&user)
	agentType, searchableDocument, err := serverhelpers.InterceptAgentType(c, user)
	if err != nil {
		log.Printf("Error loading agent type: %v", err)
		return c.SendStatus(400)
	}
	chatView := templates.ChatViewWrapper(chatMessages, agentType, searchableDocument)
	return webhelpers.RenderHttpComponent(chatView, c, ctx)
}
