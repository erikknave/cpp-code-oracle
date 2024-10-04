package endpoints

import (
	"fmt"
	"strconv"

	"github.com/erikknave/go-code-oracle/cypher/cypherqueries"
	"github.com/erikknave/go-code-oracle/web/templates"
	"github.com/erikknave/go-code-oracle/web/webhelpers"
	"github.com/gofiber/fiber/v2"
)

func FilePageEndPoint(c *fiber.Ctx) error {
	ctx := c.Context()
	var err error
	dbidStr := c.Query("dbid")
	dbid, err := strconv.Atoi(dbidStr)
	if err != nil {
		return c.SendStatus(400)
	}
	result, err := cypherqueries.PerformFileCypherQuery(fmt.Sprintf("%d", dbid))
	if err != nil {
		fmt.Printf("Error performing file query: %v", err)
		return c.SendStatus(500)
	}
	repositoryView := templates.FilePage(result)
	return webhelpers.RenderHttpComponent(repositoryView, c, ctx)
}
