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

const repoQueryTemplate = `
// Query 1: Matching and collecting repositories using or being used by the repository with dbid "{{.DBID}}"
MATCH (r1:Repository {dbid: "{{.DBID}}"})-[:HAS_MODULE]->(m1:Module)<-[:PART_OF_MODULE]-(p1:Package)-[:CONTAINS]->(f1:File)-[:DEFINES]->(e1:Entity)-[:USES]->(e2:Entity)<-[:DEFINES]-(f2:File)<-[:CONTAINS]-(p2:Package)-[:PART_OF_MODULE]->(m2:Module)<-[:HAS_MODULE]-(r2:Repository)
WHERE r2.dbid <> "{{.DBID}}" AND r2.name <> "DUMMY_REPO_NAME_NOT_EXIST"
WITH r2.name AS r2_name, r2.dbid AS r2_dbid, COUNT(*) AS r2_count
WITH r2_name, r2_dbid, r2_count
ORDER BY r2_count DESC
WITH COLLECT(DISTINCT {name: r2_name, dbid: r2_dbid, count: r2_count}) AS is_used_by_repos
OPTIONAL MATCH (r3:Repository)-[:HAS_MODULE]->(m3:Module)<-[:PART_OF_MODULE]-(p3:Package)-[:CONTAINS]->(f3:File)-[:DEFINES]->(e3:Entity)-[:USES]->(e4:Entity)<-[:DEFINES]-(f4:File)<-[:CONTAINS]-(p4:Package)-[:PART_OF_MODULE]->(m4:Module)<-[:HAS_MODULE]-(r4:Repository {dbid: "{{.DBID}}"})
WHERE r3.dbid <> "{{.DBID}}" AND r3.name <> "DUMMY_REPO_NAME_NOT_EXIST"
WITH is_used_by_repos, r3.name AS r3_name, r3.dbid AS r3_dbid, COUNT(*) AS r3_count
WITH is_used_by_repos, r3_name, r3_dbid, r3_count
ORDER BY r3_count DESC
WITH is_used_by_repos, COLLECT(DISTINCT {name: r3_name, dbid: r3_dbid, count: r3_count}) AS is_using_repos

// Query 2: Matching packages of the repository with dbid "{{.DBID}}"
OPTIONAL MATCH (r:Repository {dbid: "{{.DBID}}"})-[:HAS_MODULE]->(m:Module)-[:PART_OF_MODULE]-(p:Package)
WITH is_used_by_repos, is_using_repos, r, COLLECT({name: p.name, shortsummary: p.shortsummary, importpath: p.repoPath, dbid: p.dbid}) AS packages
// Query 3: Matching file commits affecting the repository with dbid "{{.DBID}}"
OPTIONAL MATCH (r)-[:HAS_MODULE]->(m:Module)-[:PART_OF_MODULE]-(p:Package)-[:CONTAINS]-(f:File)-[:AFFECTS]-(fc:FileCommit)
WITH is_used_by_repos, is_using_repos, packages, r.name AS name, r.dbid AS dbid, r.shortsummary AS shortsummary, r.summary as summary, COLLECT(DISTINCT fc.authorName) AS authors, MAX(fc.commitDate) AS latestUpdate
// Combining all results
RETURN COLLECT({
    is_used_by_repos: is_used_by_repos,
    is_using_repos: is_using_repos,
    name: name,
    shortsummary: shortsummary,
	summary: summary,
    authors: authors,
    latestUpdate: latestUpdate,
    packages: packages,
    dbid: dbid
}) AS result
`

func PerformRepoCypherQuery(dbid string) (types.RepoQueryReponseResult, error) {
	dbidStr := fmt.Sprintf("%s", dbid) // The arbitrary value to replace "47"
	queryParams := struct {
		DBID string
	}{
		DBID: dbidStr,
	}

	tmpl, err := template.New("query").Parse(repoQueryTemplate)
	if err != nil {
		log.Fatalf("Error parsing query template: %v", err)
	}

	var result bytes.Buffer
	if err := tmpl.Execute(&result, queryParams); err != nil {
		log.Fatalf("Error executing query template: %v", err)
	}
	cypherResult := cypher.InjectCypher(result.String())
	cypherResultJson, _ := json.Marshal(cypherResult)
	var typedResult [][]types.RepoQueryReponseResult
	err = json.Unmarshal(cypherResultJson, &typedResult)
	if err != nil {
		return types.RepoQueryReponseResult{}, err
	}
	finalResult := typedResult[0][0]
	return finalResult, nil
}
