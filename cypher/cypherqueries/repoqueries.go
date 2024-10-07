package cypherqueries

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"text/template"

	"github.com/erikknave/go-code-oracle/cypher"
	"github.com/erikknave/go-code-oracle/types"
	"github.com/gofiber/fiber/v3/log"
)

//go:embed repo-query.cql
var repoQueryTemplate string

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
