package endpoints

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/erikknave/go-code-oracle/cypher"
	"github.com/erikknave/go-code-oracle/types"
	"github.com/erikknave/go-code-oracle/web/templates"
	"github.com/erikknave/go-code-oracle/web/webhelpers"
	"github.com/gofiber/fiber/v2"
)

const queryTemplate = `
MATCH (r:Repository)
WITH COUNT(r) AS repositoryCount
MATCH (m:Module)
WITH repositoryCount, COUNT(m) AS moduleCount
MATCH (p:Package)
WITH repositoryCount, moduleCount, COUNT(p) AS packageCount
MATCH (f:File)
WITH repositoryCount, moduleCount, packageCount, COUNT(f) AS fileCount
MATCH (e:Entity)
WITH repositoryCount, moduleCount, packageCount, fileCount, COUNT(e) AS entityCount
MATCH (fc:FileCommit)
WITH repositoryCount, moduleCount, packageCount, fileCount, entityCount, COUNT(fc) AS fileCommitCount, COUNT(DISTINCT fc.authorName) AS authorCount
MATCH ()-[r]->()
WITH repositoryCount, moduleCount, packageCount, fileCount, entityCount, fileCommitCount, authorCount, COUNT(r) AS relationshipCount

// Subquery for most depended on repository
CALL {
    MATCH (r1:Repository)-[:HAS_MODULE]->(m1:Module)<-[:PART_OF_MODULE]-(p1:Package)-[:CONTAINS]->(f1:File)-[:DEFINES]->(e1:Entity)-[:USES]->(e2:Entity)<-[:DEFINES]-(f2:File)<-[:CONTAINS]-(p2:Package)-[:PART_OF_MODULE]->(m2:Module)<-[:HAS_MODULE]-(r2:Repository)
    WHERE r1 <> r2
    WITH r2, COUNT(DISTINCT r1) AS dependedOnByCount
    ORDER BY dependedOnByCount DESC
    LIMIT 1
    RETURN r2.name AS mostDependedOnName, r2.dbid AS mostDependedOnDbid, dependedOnByCount AS mostDependedOnCount
}

// Subquery for repository with most dependencies
CALL {
    MATCH (r1:Repository)-[:HAS_MODULE]->(m1:Module)<-[:PART_OF_MODULE]-(p1:Package)-[:CONTAINS]->(f1:File)-[:DEFINES]->(e1:Entity)-[:USES]->(e2:Entity)<-[:DEFINES]-(f2:File)<-[:CONTAINS]-(p2:Package)-[:PART_OF_MODULE]->(m2:Module)<-[:HAS_MODULE]-(r2:Repository)
    WHERE r1 <> r2
    WITH r1, COUNT(DISTINCT r2) AS dependingOnCount
    ORDER BY dependingOnCount DESC
    LIMIT 1
    RETURN r1.name AS mostDependenciesName, r1.dbid AS mostDependenciesDbid, dependingOnCount AS mostDependenciesCount
}

WITH repositoryCount, moduleCount, packageCount, fileCount, entityCount, fileCommitCount, authorCount, relationshipCount, 
     mostDependedOnName, mostDependedOnDbid, mostDependedOnCount, 
     mostDependenciesName, mostDependenciesDbid, mostDependenciesCount

RETURN {
  repositories: repositoryCount,
  modules: moduleCount,
  packages: packageCount,
  files: fileCount,
  entities: entityCount,
  fileCommits: fileCommitCount,
  relationships: relationshipCount,
  authors: authorCount,
  mostDependedOn: {name: mostDependedOnName, dbid: mostDependedOnDbid, count: mostDependedOnCount},
  mostDependencies: {name: mostDependenciesName, dbid: mostDependenciesDbid, count: mostDependenciesCount}
} AS results


`

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
