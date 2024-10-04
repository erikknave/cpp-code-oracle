package cypherqueries

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"

	"github.com/erikknave/go-code-oracle/cypher"
	"github.com/erikknave/go-code-oracle/types"
	"github.com/gofiber/fiber/v3/log"
)

const entityQueryTemplate = `
// Query 1: Matching and collecting repositories using or being used by the repository with dbid "{{.DBID}}"
OPTIONAL MATCH (r1:Repository)-[:HAS_MODULE]->(m1:Module)<-[:PART_OF_MODULE]-(p1:Package)-[:CONTAINS]->(f1:File)-[:DEFINES]->(e1:Entity {dbid: "{{.DBID}}"})-[:USES]->(e2:Entity)<-[:DEFINES]-(f2:File)<-[:CONTAINS]-(p2:Package)-[:PART_OF_MODULE]->(m2:Module)<-[:HAS_MODULE]-(r2:Repository)
WHERE e2.dbid <> "{{.DBID}}"
WITH e2, f2, p2, r2, COUNT(*) AS e2_count
ORDER BY e2_count DESC
WITH COLLECT(DISTINCT {name: e2.name, dbid: e2.dbid, signature: e2.signature, reponame: r2.name, repodbid: r2.dbid, count: e2_count}) AS is_used_by_entities

// Query 2: Matching repositories using the entity with dbid "{{.DBID}}"
OPTIONAL MATCH (r3:Repository)-[:HAS_MODULE]->(m3:Module)<-[:PART_OF_MODULE]-(p3:Package)-[:CONTAINS]->(f3:File)-[:DEFINES]->(e3:Entity)-[:USES]->(e4:Entity {dbid: "{{.DBID}}"})<-[:DEFINES]-(f4:File)<-[:CONTAINS]-(p4:Package)-[:PART_OF_MODULE]->(m4:Module)<-[:HAS_MODULE]-(r4:Repository)
WHERE e3.dbid <> "{{.DBID}}"
WITH is_used_by_entities, e3, f3, p3, r3, COUNT(*) AS e3_count
ORDER BY e3_count DESC
WITH is_used_by_entities, COLLECT(DISTINCT {name: e3.name, dbid: e3.dbid, signature: e3.signature, reponame: r3.name, repodbid: r3.dbid, count: e3_count}) AS is_using_entities

// Query 3: Matching file and package of the entity with dbid "{{.DBID}}"
OPTIONAL MATCH (r:Repository)-[:HAS_MODULE]->(m:Module)-[:PART_OF_MODULE]-(p:Package)-[:CONTAINS]-(f:File)-[:DEFINES]-(e:Entity {dbid: "{{.DBID}}"})
WITH is_used_by_entities, is_using_entities, e, f, p, r

// Query 4: Matching entity commits affecting the repository with dbid "{{.DBID}}"
OPTIONAL MATCH (r)-[:HAS_MODULE]->(m:Module)-[:PART_OF_MODULE]-(p)-[:CONTAINS]-(f)-[:DEFINES]-(e)-[:AFFECTS]-(fc:FileCommit)
WITH is_used_by_entities, is_using_entities, e, f, p, r, COLLECT(DISTINCT fc.authorName) AS authors, MAX(fc.commitDate) AS latestUpdate

// Combining all results
RETURN {
    is_used_by_entities: is_used_by_entities,
    is_using_entities: is_using_entities,
    name: e.name,
    summary: e.summary,
    signature: e.signature,
    authors: authors,
    latestUpdate: latestUpdate,
    dbid: e.dbid,
    repodbid: r.dbid,
    reponame: r.name,
	reposhortsummary: r.shortsummary,
	packagedbid: p.dbid,
	packageimportpath: p.repoPath,
	packageshortsummary: p.shortsummary,
	filedbid: f.dbid,
	filename: f.name,
	filesummary: f.summary,
    importpath: f.repoPath
} AS result
`

func PerformEntityCypherQuery(dbid string) (types.EntityQueryResponseResult, error) {
	dbidStr := fmt.Sprintf("%s", dbid) // The arbitrary value to replace "47"
	queryParams := struct {
		DBID string
	}{
		DBID: dbidStr,
	}

	tmpl, err := template.New("query").Parse(entityQueryTemplate)
	if err != nil {
		log.Fatalf("Error parsing query template: %v", err)
	}

	var result bytes.Buffer
	if err := tmpl.Execute(&result, queryParams); err != nil {
		log.Fatalf("Error executing query template: %v", err)
	}
	resultString := result.String()
	cypherResult := cypher.InjectCypher(resultString)
	cypherResultJson, _ := json.Marshal(cypherResult)
	var typedResult []types.EntityQueryResponseResult
	err = json.Unmarshal(cypherResultJson, &typedResult)
	if err != nil {
		return types.EntityQueryResponseResult{}, err
	}
	finalResult := typedResult[0]
	return finalResult, nil
}
