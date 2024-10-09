package listfilesforsearchidfc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/erikknave/go-code-oracle/cypher/cypherqueries"
	"github.com/erikknave/go-code-oracle/helpers"
	"github.com/erikknave/go-code-oracle/types"
	"github.com/tmc/langchaingo/llms"
)

const name = "ListFilesForSearchId"

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
			Description: "Lists all files within a repository or directory based on a search id (the search id must start with repository- or directory-)",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"searchId": map[string]any{
						"type":        "string",
						"description": "The search id of the repository or directory to list files in",
					},
					// "unit": map[string]any{
					// 	"type": "string",
					// 	"enum": []string{"fahrenheit", "celsius"},
					// },
				},
				"required": []string{"searchId"},
			},
		},
	}
}

func (f *FunctionCall) Execute(args json.RawMessage) (string, error) {
	fmt.Printf("\n - Execute function %s called\n", name)
	var params struct {
		SearchId string `json:"searchId"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return "", err
	}
	return f.Function(params.SearchId), nil

}

// type queryResponseType struct {
// 	Name string `json:"name"`
// 	Dbid string `json:"dbid"`
// }

const responseTemplate = `
The following files are part of the search id {{.Type}}-{{.Dbid}} with path {{.ImportPath}}:
{{range .Files}}
- Name: {{.Name}}
- Path: {{.ImportPath}}
- Search id (not shown to user): file-{{.Dbid}}

{{end}}
`

func (f *FunctionCall) Function(searchId string) string {
	listFilesResult, err := cypherqueries.ListFilesBasedOnSearchId(searchId)
	if err != nil {
		return fmt.Sprintf("Error in cypherqueries.ListFilesBasedOnSearchId: %v", err)
	}
	summaryString, err := helpers.CreateStringFromTemplate(responseTemplate, listFilesResult)
	if err != nil {
		return fmt.Sprintf("search files function call: Error in helpers.CreateStringFromTemplate: %v", err)
	}
	return summaryString
}
