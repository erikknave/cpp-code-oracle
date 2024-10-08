package searchallfilesfc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/erikknave/go-code-oracle/helpers"
	"github.com/erikknave/go-code-oracle/search"
	"github.com/erikknave/go-code-oracle/types"
	"github.com/tmc/langchaingo/llms"
)

const name = "SearchAllFiles"

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
			Description: "Returns the short summaries of all files, regardless of repository, related to a query, including what package and repository they belong to",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"query": map[string]any{
						"type":        "string",
						"description": "The query to search for",
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

func (f *FunctionCall) Execute(args json.RawMessage) (string, error) {
	fmt.Printf("\n - Execute function %s called\n", name)
	var params struct {
		Query string `json:"query"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return "", err
	}
	return f.Function(params.Query), nil

}

// type queryResponseType struct {
// 	Name string `json:"name"`
// 	Dbid string `json:"dbid"`
// }

const responseTemplate = `
The following files were found (presented by relevance):
{{range .}}
- Type: {{.Type}}
- Name: {{.Name}}
- Path: {{.Path}}
- Search id: file-{{.Dbid}}
- 
- Summary: {{if .ShortSummary}}{{.ShortSummary}}{{else}}{{.Summary}}{{end}}
{{end}}
`

func (f *FunctionCall) Function(queryString string) string {
	limit := 8
	searchDocs, err := search.SearchAllFiles(queryString, limit)
	if err != nil {
		return fmt.Sprintf("Error in search.SearchFiles: %v", err)
	}
	summaryString, err := helpers.CreateStringFromTemplate(responseTemplate, searchDocs)
	if err != nil {
		return fmt.Sprintf("search files function call: Error in helpers.CreateStringFromTemplate: %v", err)
	}
	return summaryString
}
