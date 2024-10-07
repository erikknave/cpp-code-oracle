package packagedeepsummaryfc

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

const name = "PackageDeepSummary"

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
			Description: "Returns the 'deep' summary of a package (roughly folder) within a certain repository (including what files it contains)",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{

					"packageSearchId": map[string]any{
						"type":        "string",
						"description": "The search id of the package to get the deep summary of",
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
				"required": []string{"packageSearchId"},
			},
		},
	}
}

func (f *FunctionCall) Execute(args json.RawMessage) (string, error) {
	fmt.Printf("\n - Execute function %s called\n", name)
	var params struct {
		PackageSearchId string `json:"packageSearchId"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return "", err
	}
	return f.Function(params.PackageSearchId), nil
}

const responseTemplate = `
The name of the package is {{.Name}} with search id package-{{.Dbid}}
The import path of the package is {{.ImportPath}}

The package belongs to the repository {{.RepoName}} with search id repository-{{.RepoDbid}}

Here is a summary of the package: {{.Summary}}

Here is a list of all files within the package:
{{- range .Files}}
Import path: {{ .ImportPath }}
Search id (not shown to user): file-{{ .Dbid }}
Summary: {{ .Summary }}

{{- end }}
`

func (f *FunctionCall) Function(packageSearchId string) string {
	requestedType := search.GetTypeFromSearchId(fmt.Sprintf("%v", packageSearchId))
	if requestedType != "package" {
		return "The dbid provided does not correspond to a package, but to a " + requestedType
	}
	packageDbid := search.GetDbidFromSearchId(fmt.Sprintf("%v", packageSearchId))
	repoResult, err := cypherqueries.PerformDirectoryCypherQuery(packageDbid)
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
