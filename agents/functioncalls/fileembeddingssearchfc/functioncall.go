package fileembeddingssearchfc

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

const name = "FileSemanticSearch"

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
			Description: "Returns the  summary of a c++ file based on an embeddings (semantic) search based on a search string",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{

					// "fileSearchId": map[string]any{
					// 	"type":        "string",
					// 	"description": "The search id of the file to get the deep summary of",
					// },
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
	return f.Function(params.Query, tCtx), nil
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

func (f *FunctionCall) Function(query string, tCtx *types.ToolContext) string {
	searchDocs := pgvector.PerformFileQuery(query, 10)
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
	for _, doc := range searchDocs {
		tCtx.MentionedFiles = append(tCtx.MentionedFiles, doc.Path)
	}
	response := result.String()
	return response
}
