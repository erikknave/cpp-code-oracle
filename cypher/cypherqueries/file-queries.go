package cypherqueries

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"github.com/erikknave/go-code-oracle/cypher"
	"github.com/erikknave/go-code-oracle/types"
	"github.com/gofiber/fiber/v3/log"
)

//go:embed file-query.cql
var fileQueryTemplate string

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
MATCH (r:REPOSITORY{dbid:%s})-[]-(d:DIRECTORY)-[]-(f:FILE)
WITH r,  collect({importPath: f.importPath,  dbid: f.dbid, name: f.fileName}) AS files
RETURN {type: 'repository', dbid: r.dbid, importPath: r.name, files: files} AS result
`

// const listFilesInRepoQueryTemplate = `
// MATCH (r:Repository{dbid:"18"})-[]-(m:Module)-[]-(p:Package)-[:CONTAINS]-(f:File)
// WITH r, collect(DISTINCT f) AS files
// RETURN {type: 'repository', dbid: r.dbid, importPath: r.name, files:
//   [f IN files | {importPath: f.repoPath, dbid: f.dbid, name: f.name}]} AS result
// `

const listFilesInDirectoryQueryTemplate = `
MATCH (r:REPOSITORY)-[]-(d:DIRECTORY {dbid:%s})-[]-(f:FILE)
WITH  d, collect({importPath: f.importPath,  dbid: f.dbid, name: f.fileName}) AS files
RETURN {type:'directory', dbid: d.dbid,  importPath: d.importPath, files: files} AS result
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
	case "directory":
		packageTemplate := fmt.Sprintf(listFilesInDirectoryQueryTemplate, words[1])
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
