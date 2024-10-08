package endpoints

import (
	"fmt"
	"strconv"

	"github.com/erikknave/go-code-oracle/cypher/cypherqueries"
	"github.com/erikknave/go-code-oracle/web/templates"
	"github.com/erikknave/go-code-oracle/web/webhelpers"
	"github.com/gofiber/fiber/v2"
)

func RepositoryPageEndPoint(c *fiber.Ctx) error {
	ctx := c.Context()
	var err error

	dbidStr := c.Query("dbid")
	dbid, err := strconv.Atoi(dbidStr)
	if err != nil {
		return c.SendStatus(400)
	}
	// user := c.Locals("user").(types.User)
	// viewName := c.Query("view")
	// if viewName == "search" {
	// 	searchResults, err := dbhelpers.LoadUserSearchResults(&user)
	// 	if err != nil {
	// 		log.Fatalf("Error loading search results: %v", err)
	// 	}
	// 	return webhelpers.RenderHttpComponent(templates.SearchPage(&searchResults, ""), c, ctx)

	// }
	// chatMessages, _ := dbhelpers.LoadChatMessagesForUser(&user)
	repoResult, err := cypherqueries.PerformRepoCypherQuery(fmt.Sprintf("%d", dbid))
	if err != nil {
		fmt.Printf("Error performing repo query: %v", err)
		return c.SendStatus(500)
	}
	repositoryView := templates.RepositoryPage(repoResult)
	return webhelpers.RenderHttpComponent(repositoryView, c, ctx)
}
