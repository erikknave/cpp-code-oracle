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
