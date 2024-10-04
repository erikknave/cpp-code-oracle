package cypherqueries

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"github.com/erikknave/go-code-oracle/cypher"
	"github.com/erikknave/go-code-oracle/types"
	"github.com/gofiber/fiber/v3/log"
)

const fileQueryTemplate = `
// Query 1: Matching and collecting repositories using or being used by the repository with dbid "{{.DBID}}"
OPTIONAL MATCH (r1:Repository )-[:HAS_MODULE]->(m1:Module)<-[:PART_OF_MODULE]-(p1:Package )-[:CONTAINS]->(f1:File {dbid: "{{.DBID}}"})-[:DEFINES]->(e1:Entity)-[:USES]->(e2:Entity)<-[:DEFINES]-(f2:File)<-[:CONTAINS]-(p2:Package)-[:PART_OF_MODULE]->(m2:Module)<-[:HAS_MODULE]-(r2:Repository)
WHERE f2.dbid <> "{{.DBID}}" and f2.name <> "NON_EXISTING_FILE.go"
WITH f2, p2, r2, COUNT(*) AS f2_count
ORDER BY f2_count DESC
WITH COLLECT(DISTINCT {name: f2.name, dbid: f2.dbid, importpath: f2.repoPath, reponame: r2.name, repodbid: r2.dbid, count: f2_count}) AS is_used_by_files

// Query 2: Matching repositories using the file with dbid "{{.DBID}}"
OPTIONAL MATCH (r3:Repository)-[:HAS_MODULE]->(m3:Module)<-[:PART_OF_MODULE]-(p3:Package)-[:CONTAINS]->(f3:File)-[:DEFINES]->(e3:Entity)-[:USES]->(e4:Entity)<-[:DEFINES]-(f4:File {dbid: "{{.DBID}}"})<-[:CONTAINS]-(p4:Package)-[:PART_OF_MODULE]->(m4:Module)<-[:HAS_MODULE]-(r4:Repository )
WHERE f3.dbid <> "{{.DBID}}" and f3.name <> "NON_EXISTING_FILE.go"
WITH is_used_by_files, f3, p3, r3, COUNT(*) AS f3_count
ORDER BY f3_count DESC
WITH is_used_by_files, COLLECT(DISTINCT {name: f3.name, dbid: f3.dbid, importpath: f3.repoPath, reponame: r3.name, repodbid: r3.dbid, count: f3_count}) AS is_using_files

// Query 3: Matching packages and entities of the repository with dbid "{{.DBID}}"
OPTIONAL MATCH (r:Repository )-[:HAS_MODULE]->(m:Module)-[:PART_OF_MODULE]-(p:Package)-[:CONTAINS]-(f:File {dbid: "{{.DBID}}"})-[:DEFINES]-(e:Entity)
WITH is_used_by_files, is_using_files, f, p, r, COLLECT({name: e.name, summary: e.summary, signature: e.signature, dbid: e.dbid}) AS entities

// Query 4: Matching file commits affecting the repository with dbid "{{.DBID}}"
OPTIONAL MATCH (r)-[:HAS_MODULE]->(m:Module)-[:PART_OF_MODULE]-(p)-[:CONTAINS]-(f)-[:AFFECTS]-(fc:FileCommit)
WITH is_used_by_files, is_using_files, entities, r, f, p, COLLECT(DISTINCT fc.authorName) AS authors, MAX(fc.commitDate) AS latestUpdate

// Combining all results
RETURN {
    is_used_by_files: is_used_by_files,
    is_using_files: is_using_files,
    name: f.name,
    summary: f.summary,
    authors: authors,
    latestUpdate: latestUpdate,
    entities: entities,
    dbid: f.dbid,
    repodbid: r.dbid,
    reponame: r.name,
	packagedbid: p.dbid,
	packageimportpath: p.repoPath,
    importpath: f.repoPath
} AS result


`

func PerformFileCypherQuery(dbid string) (types.FileQueryReponseResult, error) {
	dbidStr := fmt.Sprintf("%s", dbid) // The arbitrary value to replace "47"
	queryParams := struct {
		DBID string
	}{
		DBID: dbidStr,
	}

	tmpl, err := template.New("query").Parse(fileQueryTemplate)
	if err != nil {
		log.Fatalf("Error parsing query template: %v", err)
	}

	var result bytes.Buffer
	if err := tmpl.Execute(&result, queryParams); err != nil {
		log.Fatalf("Error executing query template: %v", err)
	}
	cypherResult := cypher.InjectCypher(result.String())
	cypherResultJson, _ := json.Marshal(cypherResult)
	var typedResult []types.FileQueryReponseResult
	err = json.Unmarshal(cypherResultJson, &typedResult)
	if err != nil {
		return types.FileQueryReponseResult{}, err
	}
	finalResult := typedResult[0]
	return finalResult, nil
}

const listFilesInRepoQueryTemplate = `
MATCH (r:Repository{dbid:"%s"})-[]-(m:Module)-[]-(p:Package)-[:CONTAINS]-(f:File)
WITH r,  collect({importPath: f.repoPath,  dbid: f.dbid, name: f.name}) AS files
RETURN {type: 'repository', dbid: r.dbid, importPath: r.name, files: files} AS result
`

// const listFilesInRepoQueryTemplate = `
// MATCH (r:Repository{dbid:"18"})-[]-(m:Module)-[]-(p:Package)-[:CONTAINS]-(f:File)
// WITH r, collect(DISTINCT f) AS files
// RETURN {type: 'repository', dbid: r.dbid, importPath: r.name, files:
//   [f IN files | {importPath: f.repoPath, dbid: f.dbid, name: f.name}]} AS result
// `

const listFilesInPackageQueryTemplate = `
MATCH (r:Repository)-[]-(m:Module)-[]-(p:Package {dbid:"%s"})-[:CONTAINS]-(f:File)
WITH  p, collect({importPath: f.repoPath,  dbid: f.dbid, name: f.name}) AS files
RETURN {type:'package', dbid: p.dbid,  importPath: p.repoPath, files: files} AS result
`

func ListFilesBasedOnSearchId(searchId string) (types.ListFilesResponseResult, error) {
	words := strings.Split(searchId, "-")
	typeOfSearch := words[0]
	switch typeOfSearch {
	case "repository":
		repoTemplate := fmt.Sprintf(listFilesInRepoQueryTemplate, words[1])
		cypherResult := cypher.InjectCypher(repoTemplate)
		cypherResultJson, _ := json.Marshal(cypherResult)
		var typedResult []types.ListFilesResponseResult
		err := json.Unmarshal(cypherResultJson, &typedResult)
		if err != nil {
			return types.ListFilesResponseResult{}, err
		}
		finalResult := typedResult[0]
		return finalResult, nil
	case "package":
		packageTemplate := fmt.Sprintf(listFilesInPackageQueryTemplate, words[1])
		cypherResult := cypher.InjectCypher(packageTemplate)
		cypherResultJson, _ := json.Marshal(cypherResult)
		var typedResult []types.ListFilesResponseResult
		err := json.Unmarshal(cypherResultJson, &typedResult)
		if err != nil {
			return types.ListFilesResponseResult{}, err
		}
		finalResult := typedResult[0]
		return finalResult, nil
	}
	return types.ListFilesResponseResult{}, fmt.Errorf("The search id %s is not valid", searchId)
}
