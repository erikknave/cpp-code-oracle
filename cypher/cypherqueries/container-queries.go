package cypherqueries

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"text/template"

	"github.com/erikknave/go-code-oracle/cypher"
	"github.com/erikknave/go-code-oracle/types"
	"github.com/gofiber/fiber/v3/log"
)

//go:embed container-query.cql
var containerQueryTemplate string

func PerformContainerCypherQuery(dbid string) (types.ContainerQueryResponseResult, error) {
	dbidStr := dbid
	queryParams := struct {
		DBID string
	}{
		DBID: dbidStr,
	}

	tmpl, err := template.New("query").Parse(containerQueryTemplate)
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
	var typedResult []types.ContainerQueryResponseResult
	err = json.Unmarshal(cypherResultJson, &typedResult)
	if err != nil {
		return types.ContainerQueryResponseResult{}, err
	}
	finalResult := typedResult[0]
	return finalResult, nil
}
