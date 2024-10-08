package endpoints

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/erikknave/go-code-oracle/cypher"
	"github.com/erikknave/go-code-oracle/types"
	"github.com/erikknave/go-code-oracle/web/templates"
	"github.com/erikknave/go-code-oracle/web/webhelpers"
	"github.com/gofiber/fiber/v2"
)

//go:embed stats-query.cql
var queryTemplate string

func StatsPageEndPoint(c *fiber.Ctx) error {
	ctx := context.Background()
	stats := getStats(queryTemplate)

	return webhelpers.RenderHttpComponent(templates.StatsPage(stats), c, ctx)
}

func StatsViewWrapperEndPoint(c *fiber.Ctx) error {
	ctx := context.Background()
	stats := getStats(queryTemplate)
	return webhelpers.RenderHttpComponent(templates.StatsViewWrapper(stats), c, ctx)
}

func getStats(queryString string) *types.Stats {
	cypherResult := cypher.InjectCypher(queryString)
	cypherResultJson, err := json.Marshal(cypherResult)
	if err != nil {
		fmt.Printf("error in json.Marshal: %v", err)
	}
	var typedResult []types.Stats
	err = json.Unmarshal(cypherResultJson, &typedResult)
	if err != nil {
		fmt.Printf("error in json.Unmarshal: %v", err)
	}
	returnedStats := typedResult[0]
	return &returnedStats
}

func CommandEndPoint(c *fiber.Ctx) error {
	ctx := context.Background()
	command := c.FormValue("prompt")
	if strings.Contains(command, "/chat") {
		return ChatViewWrapperEndPoint(c)
	}
	if strings.Contains(command, "/search") {
		return PerformSearchEndPoint(c)
	}
	if strings.Contains(command, "/stats") {
		return StatsViewWrapperEndPoint(c)
	}
	if strings.Contains(command, "/help") {
		return HelpViewWrapperEndPoint(c)
	}
	return webhelpers.RenderHttpComponent(templates.StatsPrompt(), c, ctx)
}
