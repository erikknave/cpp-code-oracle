package filedeepsummaryfc

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

const name = "FileDeepSummary"

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
			Description: "Returns the 'deep' summary of a c++ file within a certain repository (including what functions and variables it contains)",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{

					"fileSearchId": map[string]any{
						"type":        "string",
						"description": "The search id of the file to get the deep summary of",
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
				"required": []string{"fileSearchId"},
			},
		},
	}
}

func (f *FunctionCall) Execute(args json.RawMessage) (string, error) {
	fmt.Printf("\n - Execute function %s called\n", name)
	var params struct {
		FileSearchId string `json:"fileSearchId"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return "", err
	}
	return f.Function(params.FileSearchId), nil
}

const responseTemplate = `
The name of the file is {{.Name}} with search id file-{{.Dbid}}.
The import path of the file is {{.ImportPath}}.

The file is located in the repository {{.RepoName}} with search id repository-{{.RepoDbid}} and in the directory {{.DirectoryImportPath}} with search id directory-{{.DirectoryDbid}}.
The summary of the file is: {{.Summary}}

Here is a list of all variables, constants, functions, types and methods  within the file:
{{- range .Entities}}
Name: {{  .Name }}
Signature: {{ .Signature }}
Search id (not shown to user): entity-{{ .Dbid }}
Summary: {{ .Summary }}
{{- end }}

`

func (f *FunctionCall) Function(fileSearchId string) string {
	requestedType := search.GetTypeFromSearchId(fileSearchId)
	if requestedType != "file" {
		return "The dbid provided does not correspond to a file, but to a " + requestedType
	}
	fileDbid := search.GetDbidFromSearchId(fileSearchId)

	repoResult, err := cypherqueries.PerformFileCypherQuery(fileDbid)
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
