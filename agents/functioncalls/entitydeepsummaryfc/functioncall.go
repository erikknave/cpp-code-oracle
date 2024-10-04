package entitydeepsummaryfc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"

	"github.com/erikknave/go-code-oracle/cypher/cypherqueries"
	"github.com/erikknave/go-code-oracle/search"
	"github.com/erikknave/go-code-oracle/types"
	"github.com/tmc/langchaingo/llms"
)

const name = "EntityDeepSummary"

type FunctionCall struct {
	User *types.User
}

func CreateNewFunctionCall(c context.Context) *FunctionCall {
	user := c.Value(types.CtxKey("user")).(types.User)
	return &FunctionCall{
		User: &user,
	}
}

func (f *FunctionCall) Name() string {
	return name
}

func (f *FunctionCall) ToolDefinition() llms.Tool {
	return llms.Tool{
		Type: "function",
		Function: &llms.FunctionDefinition{
			Name:        name,
			Description: "Returns the 'deep' summary of an entity, i.e. variable or function, within a certain file within a repository",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"entitySearchId": map[string]any{
						"type":        "string",
						"description": "The search id of the entity, i.e. variable or function, to get the deep summary of",
					},
					// "query": map[string]any{
					// 	"type":        "string",
					// 	"description": "The query to find the packages related to",
					// },
					// "unit": map[string]any{
					// 	"type": "string",
					// 	"enum": []string{"fahrenheit", "celsius"},
					// },
				},
				"required": []string{"entitySearchId"},
			},
		},
	}
}

func (f *FunctionCall) Execute(args json.RawMessage) (string, error) {
	fmt.Printf("\n - Execute function %s called\n", name)
	var params struct {
		EntitySearchId string `json:"entitySearchId"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return "", err
	}
	return f.Function(params.EntitySearchId), nil
}

const responseTemplate = `
The name of the function, struct, variable, const or type is is {{.Name}} with search id entity-{{.Dbid}}.
The signature is: {{.Signature}}
The summary is: {{.Summary}}

Information regarding the file containing the entity:
Import path: {{.ImportPath}}
File search id (not shown to user): file-{{.FileDbid}}
File summary: {{.FileSummary}}

Information regarding the Package containing the file:
Package path: {{.PackageImportPath}}
Package search id (not shown to user): package-{{.PackageDbid}}
Package summary: {{.PackageShortSummary}}

Information regarding the Repository containing the package:
Repository name: {{.RepoName}}
Repository search-id (not shown to user): repository-{{.RepoDbid}}
Repository summary: {{.RepoShortSummary}}
`

func (f *FunctionCall) Function(entitySearchId string) string {
	requestedType := search.GetTypeFromSearchId(entitySearchId)
	if requestedType != "entity" {
		return "The entity search id provided does not correspond to an entity, but to a " + requestedType
	}
	entityDbid := search.GetDbidFromSearchId(entitySearchId)
	repoResult, err := cypherqueries.PerformEntityCypherQuery(entityDbid)
	if err != nil {
		return "An neo4j error occurred while performing the query: " + err.Error()
	}
	tmpl, err := template.New("systemMessage").Parse(responseTemplate)
	if err != nil {
		return fmt.Sprintf("Error in template.New: %v", err)
	}

	var result bytes.Buffer
	err = tmpl.Execute(&result, repoResult)
	if err != nil {
		return fmt.Sprintf("Error in tmpl.Execute: %v", err)
	}

	response := result.String()
	return response
}
