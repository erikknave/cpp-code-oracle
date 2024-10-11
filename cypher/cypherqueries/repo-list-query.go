package cypherqueries

import (
	"bytes"
	"encoding/json"
	"text/template"

	"github.com/erikknave/go-code-oracle/cypher"
	"github.com/erikknave/go-code-oracle/types"
	"github.com/gofiber/fiber/v3/log"
)

const repoListQueryTemplate = `
MATCH (r:REPOSITORY) where NOT r.name CONTAINS 'DUMMY_'
WITH {name: r.name, dbid: r.dbid, summary: r.summary} AS repo
RETURN collect(repo) AS result
`

func PerformRepoListCypherQuery() ([]types.RepoListQueryResult, error) {
	tmpl, err := template.New("query").Parse(repoListQueryTemplate)
	if err != nil {
		log.Fatalf("Error parsing query template: %v", err)
	}

	var query bytes.Buffer
	err = tmpl.Execute(&query, nil)
	if err != nil {
		log.Fatalf("Error executing query template: %v", err)
	}

	cypherResult := cypher.InjectCypher(query.String())
	cypherResultJson, _ := json.Marshal(cypherResult)

	var typedResponse [][]types.RepoListQueryResult
	err = json.Unmarshal(cypherResultJson, &typedResponse)
	if err != nil {
		log.Fatalf("Error unmarshalling result: %v", err)
	}

	return typedResponse[0], nil
}
