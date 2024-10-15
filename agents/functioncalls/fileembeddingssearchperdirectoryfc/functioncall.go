package fileembeddingssearchperdirectoryfc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"

	"github.com/erikknave/go-code-oracle/server/pgvector"
	"github.com/erikknave/go-code-oracle/types"
	"github.com/tmc/langchaingo/llms"
)

const name = "FileSemanticSearchPerDirectory"

type FunctionCall struct {
	User *types.User
	Dbid int
}

func CreateNewFunctionCall(c context.Context, dbid int) *FunctionCall {
	user := c.Value(types.CtxKey("user")).(types.User)
	return &FunctionCall{
		User: &user,
		Dbid: dbid,
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
			Description: "Returns the  summary of a c++ file based on an embeddings (semantic) search based on a search string within a directory",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"query": map[string]any{
						"type":        "string",
						"description": "The query to find the files related to",
					},
					// "unit": map[string]any{
					// 	"type": "string",
					// 	"enum": []string{"fahrenheit", "celsius"},
					// },
				},
				"required": []string{"query"},
			},
		},
	}
}

func (f *FunctionCall) Execute(args json.RawMessage, tCtx *types.ToolContext) (string, error) {
	fmt.Printf("\n - Execute function %s called\n", name)
	var params struct {
		Query string `json:"query"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return "", err
	}
	return f.Function(params.Query, f.Dbid), nil
}

const responseTemplate = `
The following files were found (presented by relevance):
{{range .}}
- Search ID: file-{{.Dbid}}
- Name: {{.Name}}
- Import Path: {{.Path}}
- Summary: {{.Summary}}
{{end}}
`

func (f *FunctionCall) Function(query string, dbid int) string {
	searchDocs := pgvector.PerformFileQueryPerDirectory(query, 20, dbid)
	if len(searchDocs) == 0 {
		return "No files found with the provided query"
	}

	tmpl, err := template.New("systemMessage").Parse(responseTemplate)
	if err != nil {
		return fmt.Sprintf("Error in template.New: %v", err)
	}

	var result bytes.Buffer
	err = tmpl.Execute(&result, searchDocs)
	if err != nil {
		return fmt.Sprintf("Error in tmpl.Execute: %v", err)
	}

	response := result.String()
	return response
}
