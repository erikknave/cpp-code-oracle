package endpoints

import (
	"fmt"
	"strconv"

	"github.com/erikknave/go-code-oracle/cypher/cypherqueries"
	"github.com/erikknave/go-code-oracle/filecontent"
	"github.com/erikknave/go-code-oracle/web/templates"
	"github.com/erikknave/go-code-oracle/web/webhelpers"
	"github.com/gofiber/fiber/v2"
)

func CodeblockPageEndPoint(c *fiber.Ctx) error {
	ctx := c.Context()
	var err error
	dbidStr := c.Query("dbid")
	dbid, err := strconv.Atoi(dbidStr)
	if err != nil {
		return c.SendStatus(400)
	}
	result, err := cypherqueries.PerformCodeblockCypherQuery(fmt.Sprintf("%d", dbid))
	if err != nil {
		fmt.Printf("Error performing file query: %v", err)
		return c.SendStatus(500)
	}
	content, err := filecontent.GetFilePartContent(result.ImportPath, int64(result.StartOffSet), int64(result.EndOffSet))
	if err != nil {
		result.Content = "Error getting file content" + err.Error()
	} else {
		result.Content = content
	}
	codeblockView := templates.CodeblockPage(result)
	return webhelpers.RenderHttpComponent(codeblockView, c, ctx)
}
