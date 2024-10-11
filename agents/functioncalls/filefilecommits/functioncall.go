package filefilecommits

import (
	"encoding/json"
	"fmt"

	"github.com/erikknave/go-code-oracle/cypher"
	"github.com/erikknave/go-code-oracle/helpers"
	"github.com/erikknave/go-code-oracle/types"
	"github.com/tmc/langchaingo/llms"
)

const name = "FileCommits"

type FunctionCall struct {
	Dbid int
	User types.User
}

func CreateNewFunctionCall(u *types.User, dbid int) *FunctionCall {

	return &FunctionCall{
		Dbid: dbid,
		User: *u,
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
			Description: "Returns a summary of the latest commits in the file (if there are any). The function will return who updated the file, when it was updated, and a summary of the changes. (Note that it might contain information on changes done also outside of this file)",
			Parameters: map[string]any{
				"type":       "object",
				"properties": map[string]any{
					// "query": map[string]any{
					// 	"type":        "string",
					// 	"description": "The query to find the packages related to",
					// },
					// "unit": map[string]any{
					// 	"type": "string",
					// 	"enum": []string{"fahrenheit", "celsius"},
					// },
				},
				// "required": []string{"query"},
			},
		},
	}
}

func (f *FunctionCall) Execute(args json.RawMessage) (string, error) {
	fmt.Printf("\n - Execute function %s called\n", name)
	// var params struct {
	// 	Query string `json:"query"`
	// }
	// if err := json.Unmarshal(args, &params); err != nil {
	// 	return "", err
	// }
	// return f.Function(params.Query), nil
	return f.Function(), nil

}

type queryResponseType struct {
	Summary    string `json:"summary"`
	Author     string `json:"author"`
	CommitDate string `json:"commitDate"`
}

const queryString = `
MATCH (r:REPOSITORY )-[]-(d:DIRECTORY)-[]-(f:FILE  {dbid:"%s"})-[]-(fc:FILECOMMIT)
with {author:fc.authorName, commitDate:fc.commitDate, summary: fc.summary} as commit 
return collect(commit) as result
`

const responseTemplate = `
The following commits were found:
{{range .}}
- Name: {{.Author}}
- Commit summary: {{.Summary}}
- Commit Date: {{.CommitDate}}

{{end}}
`

func (f *FunctionCall) Function() string {
	queryString := fmt.Sprintf(queryString, f.Dbid)
	cypherResult := cypher.InjectCypher(queryString)
	cypherResultJson, err := json.Marshal(cypherResult)
	if err != nil {
		return fmt.Sprintf("error in json.Marshal: %v", err)
	}

	var typedResult [][]queryResponseType
	err = json.Unmarshal(cypherResultJson, &typedResult)
	if err != nil {
		return fmt.Sprintf("error in json.Unmarshal: %v", err)
	}

	if len(typedResult) == 0 {
		return fmt.Sprintf("unexpected result format: %v", typedResult)
	}
	finalResult := typedResult[0]
	summaryString, err := helpers.CreateStringFromTemplate(responseTemplate, finalResult)
	if err != nil {
		return fmt.Sprintf("Error in helpers.CreateStringFromTemplate: %v", err)
	}
	return fmt.Sprintf("The search responded with these commits (Note that this is raw data that might need to be formatted):\n%s", summaryString)
	// limit := 5
	// dbid := f.SearchableDocument.Dbid
	// searchDocs, err := search.SearchPackages(queryString, fmt.Sprintf("%d", dbid), limit)
	// if err != nil {
	// 	return fmt.Sprintf("Error in search.SearchPackages: %v", err)
	// }
	// summaryString, err := helpers.CreateStringFromTemplate(responseTemplate, searchDocs)
	// if err != nil {
	// 	return fmt.Sprintf("Error in helpers.CreateStringFromTemplate: %v", err)
	// }
	// return fmt.Sprintf("The search responded with these %d packages:\n%s", limit, summaryString)
}
