package packagesshortsummaryfcdbid

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/erikknave/go-code-oracle/helpers"
	"github.com/erikknave/go-code-oracle/search"
	"github.com/erikknave/go-code-oracle/types"
	"github.com/tmc/langchaingo/llms"
)

const name = "PackagesShortSummaryWithinRepository"

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
			Description: "Returns the short summaries of a number of go packages (roughly translated to folders within a repository) within a certain repository related to a query",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"repositorySearchId": map[string]any{
						"type":        "string",
						"description": "The search id of the repository to search in",
					},
					"query": map[string]any{
						"type":        "string",
						"description": "The query to find the packages related to",
					},
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
		Query              string `json:"query"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return "", err
	}
	return f.Function(params.RepositorySearchId, params.Query), nil
}

// type queryResponseType struct {
// 	Name string `json:"name"`
// 	Dbid string `json:"dbid"`
// }

const responseTemplate = `
The following packages were found (presented by relevance):
{{range .}}
- search id: package-{{.Dbid}}
- Name: {{.Name}}
- Import Path: {{.Path}}
- Summary: {{.ShortSummary}}

{{end}}
`

func (f *FunctionCall) Function(inputSearchId string, queryString string) string {
	limit := 8
	requestedType := search.GetTypeFromSearchId(inputSearchId)
	if requestedType != "repository" {
		return "The repository dbid is not a repository, but a " + requestedType
	}
	dbid := search.GetDbidFromSearchId(inputSearchId)
	searchDocs, err := search.SearchPackages(queryString, dbid, limit)
	if err != nil {
		return fmt.Sprintf("Error in search.SearchPackages: %v", err)
	}
	summaryString, err := helpers.CreateStringFromTemplate(responseTemplate, searchDocs)
	if err != nil {
		return fmt.Sprintf("Error in helpers.CreateStringFromTemplate: %v", err)
	}
	return fmt.Sprintf("The search responded with these %d packages:\n%s", limit, summaryString)
}
