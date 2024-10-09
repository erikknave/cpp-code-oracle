package directoriesshortsummaryfc

import (
	"encoding/json"
	"fmt"

	"github.com/erikknave/go-code-oracle/helpers"
	"github.com/erikknave/go-code-oracle/search"
	"github.com/erikknave/go-code-oracle/types"
	"github.com/tmc/langchaingo/llms"
)

const name = "DirectoriesShortSummary"

type FunctionCall struct {
	Dbid int
	User *types.User
}

func CreateNewFunctionCall(user *types.User, dbid int) *FunctionCall {
	return &FunctionCall{
		Dbid: dbid,
		User: user,
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
			Description: "Returns the short summaries of a number of C++ directories related to a query",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"query": map[string]any{
						"type":        "string",
						"description": "The query to find the directories related to",
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
The following directories were found (presented by relevance):
{{range .}}
- Name: {{.Name}}
- Import Path: {{.Path}}
- Summary: {{.ShortSummary}}
{{end}}
`

func (f *FunctionCall) Function(queryString string) string {
	limit := 5
	dbid := f.Dbid
	searchDocs, err := search.SearchDirectories(queryString, fmt.Sprintf("%d", dbid), limit)
	if err != nil {
		return fmt.Sprintf("Error in search.SearchPackages: %v", err)
	}
	summaryString, err := helpers.CreateStringFromTemplate(responseTemplate, searchDocs)
	if err != nil {
		return fmt.Sprintf("Error in helpers.CreateStringFromTemplate: %v", err)
	}
	return fmt.Sprintf("The search responded with these %d packages:\n%s", limit, summaryString)
}
