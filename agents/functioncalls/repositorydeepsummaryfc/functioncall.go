package repositorydeepsummaryfc

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

const name = "RepositoryDeepSummary"

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
			Description: "Returns the 'deep' summary of a repository based on the repository's dbid",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"repositorySearchId": map[string]any{
						"type":        "string",
						"description": "The search id of the repository to search in",
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
				"required": []string{"repositorySearchId", "query"},
			},
		},
	}
}

func (f *FunctionCall) Execute(args json.RawMessage) (string, error) {
	fmt.Printf("\n - Execute function %s called\n", name)
	var params struct {
		RepositorySearchId string `json:"repositorySearchId"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return "", err
	}
	return f.Function(params.RepositorySearchId), nil
}

// type queryResponseType struct {
// 	Name string `json:"name"`
// 	Dbid string `json:"dbid"`
// }

const responseTemplate = `
The search id of the repository is repository-{{.Dbid}} 

Here is a summary of the repository: {{.Summary}}

Here is a list of all packages within the repository:
{{- range .Packages}}
Package path: {{.ImportPath}} with name {{.Name}} and search id package-{{.Dbid}}
{{- end }}

{{ if .Authors }}
Here is a list of code update authors (might be displayed as names, logins or e-mails) that have contributed to the repository
{{ range .Authors }}
 - {{ . }}
{{ end }}
 {{ else }}
There are currently no registered updates and authors for this repository
{{ end }}
`

func (f *FunctionCall) Function(inputSearchId string) string {
	requestedType := search.GetTypeFromSearchId(inputSearchId)
	if requestedType != "repository" {
		return "The dbid provided does not correspond to a repository, but to a " + requestedType
	}
	dbid := search.GetDbidFromSearchId(inputSearchId)
	repoResult, err := cypherqueries.PerformRepoCypherQuery(dbid)
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
