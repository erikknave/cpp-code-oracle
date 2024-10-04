package fileshortsummaryfc

import (
	"encoding/json"
	"fmt"

	"github.com/erikknave/go-code-oracle/helpers"
	"github.com/erikknave/go-code-oracle/search"
	"github.com/erikknave/go-code-oracle/types"
	"github.com/tmc/langchaingo/llms"
)

const name = "FilesShortSummary"

type FunctionCall struct {
	Dbid int
	User types.User
}

func CreateNewFunctionCall(user *types.User, dbid int) *FunctionCall {

	return &FunctionCall{
		Dbid: dbid,
		User: *user,
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
			Description: "Returns the short summaries and dbid:s (used when asking the file agent) of a number of go files related to a query",
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
- Name: {{.Name}}
- Import Path: {{.Path}}
- Summary: {{.Summary}}
- Dbid (Used when asking the file agent): {{.Dbid}}
{{end}}
`

func (f *FunctionCall) Function(queryString string) string {
	limit := 5
	dbid := f.Dbid
	searchDocs, err := search.SearchFiles(queryString, fmt.Sprintf("%d", dbid), limit)
	if err != nil {
		return fmt.Sprintf("Error in search.SearchFiles: %v", err)
	}
	summaryString, err := helpers.CreateStringFromTemplate(responseTemplate, searchDocs)
	if err != nil {
		return fmt.Sprintf("Files Short summary function call: Error in helpers.CreateStringFromTemplate: %v", err)
	}
	return fmt.Sprintf("The search responded with these %d files:\n%s", limit, summaryString)
}
