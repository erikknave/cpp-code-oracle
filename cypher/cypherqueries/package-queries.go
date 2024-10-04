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

const packageQueryTemplate = `
// Query 1: Matching and collecting repositories using or being used by the repository with dbid "{{.DBID}}"
MATCH (r1:Repository )-[:HAS_MODULE]->(m1:Module)<-[:PART_OF_MODULE]-(p1:Package {dbid: "{{.DBID}}"})-[:CONTAINS]->(f1:File)-[:DEFINES]->(e1:Entity)-[:USES]->(e2:Entity)<-[:DEFINES]-(f2:File)<-[:CONTAINS]-(p2:Package)-[:PART_OF_MODULE]->(m2:Module)<-[:HAS_MODULE]-(r2:Repository)
WHERE p2.dbid <> "{{.DBID}}"
WITH p2.name AS p2_name, p2.dbid AS p2_dbid, p2.repoPath as p2_repopath, r2.name as r2_name, r2.dbid as r2_dbid, COUNT(*) AS p2_count
WITH p2_name, p2_dbid, p2_count, p2_repopath, r2_name, r2_dbid
ORDER BY p2_count DESC
WITH COLLECT(DISTINCT {name: p2_name, dbid: p2_dbid, importpath: p2_repopath, reponame: r2_name, repodbid: r2_dbid, count: p2_count}) AS is_used_by_packages
OPTIONAL MATCH (r3:Repository)-[:HAS_MODULE]->(m3:Module)<-[:PART_OF_MODULE]-(p3:Package)-[:CONTAINS]->(f3:File)-[:DEFINES]->(e3:Entity)-[:USES]->(e4:Entity)<-[:DEFINES]-(f4:File)<-[:CONTAINS]-(p4:Package {dbid: "{{.DBID}}"})-[:PART_OF_MODULE]->(m4:Module)<-[:HAS_MODULE]-(r4:Repository )
WHERE p3.dbid <> "{{.DBID}}"
WITH is_used_by_packages, p3.name AS p3_name, p3.dbid AS p3_dbid, p3.repoPath as p3_repopath, r3.dbid as r3_dbid, r3.name as r3_name, COUNT(*) AS p3_count
WITH is_used_by_packages, p3_name, p3_dbid, p3_repopath, p3_count, r3_name, r3_dbid
ORDER BY p3_count DESC
WITH is_used_by_packages, COLLECT(DISTINCT {name: p3_name, dbid: p3_dbid, importpath: p3_repopath, repodbid: r3_dbid, reponame: r3_name, count: p3_count}) AS is_using_packages

// Query 2: Matching packages of the repository with dbid "{{.DBID}}"
OPTIONAL MATCH (r:Repository )-[:HAS_MODULE]->(m:Module)-[:PART_OF_MODULE]-(p:Package{dbid: "{{.DBID}}"})-[:CONTAINS]-(f:File)
where f.name <> "NON_EXISTING_FILE.go"
WITH is_used_by_packages, is_using_packages, p, r,COLLECT({name: f.name, summary: f.summary,  importpath: f.repoPath, dbid: f.dbid}) AS files
// Query 3: Matching file commits affecting the repository with dbid "{{.DBID}}"
OPTIONAL MATCH (r)-[:HAS_MODULE]->(m:Module)-[:PART_OF_MODULE]-(p)-[:CONTAINS]-(f5:File)-[:AFFECTS]-(fc:FileCommit)
WITH is_used_by_packages, is_using_packages, files,r, p.name AS name, p.dbid AS dbid, p.summary as summary, p.shortsummary AS shortsummary, p.repoPath as importpath, r.name as reponame, r.dbid as repodbid, COLLECT(DISTINCT fc.authorName) AS authors, MAX(fc.commitDate) AS latestUpdate
// Combining all results
RETURN COLLECT({
    is_used_by_packages: is_used_by_packages,
    is_using_packages: is_using_packages,
    name: name,
    shortsummary: shortsummary,
	summary: summary,
    authors: authors,
    latestUpdate: latestUpdate,
    files: files,
    dbid: dbid,
    repodbid: repodbid,
    reponame: reponame,
    importpath: importpath
}) AS result
`

func PerformPackageCypherQuery(dbid string) (types.PackageQueryReponseResult, error) {
	dbidStr := fmt.Sprintf("%s", dbid) // The arbitrary value to replace "47"
	queryParams := struct {
		DBID string
	}{
		DBID: dbidStr,
	}

	tmpl, err := template.New("query").Parse(packageQueryTemplate)
	if err != nil {
		log.Fatalf("Error parsing query template: %v", err)
	}

	var result bytes.Buffer
	if err := tmpl.Execute(&result, queryParams); err != nil {
		log.Fatalf("Error executing query template: %v", err)
	}
	cypherResult := cypher.InjectCypher(result.String())
	cypherResultJson, _ := json.Marshal(cypherResult)
	var typedResult [][]types.PackageQueryReponseResult
	err = json.Unmarshal(cypherResultJson, &typedResult)
	if err != nil {
		return types.PackageQueryReponseResult{}, err
	}
	finalResult := typedResult[0][0]
	return finalResult, nil
}
