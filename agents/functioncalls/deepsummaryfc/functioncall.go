package deepsummaryfc

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/erikknave/go-code-oracle/agents/functioncalls/directorydeepsummaryfc"
	"github.com/erikknave/go-code-oracle/agents/functioncalls/entitydeepsummaryfc"
	"github.com/erikknave/go-code-oracle/agents/functioncalls/filedeepsummaryfc"
	"github.com/erikknave/go-code-oracle/agents/functioncalls/repositorydeepsummaryfc"
	"github.com/erikknave/go-code-oracle/types"
	"github.com/tmc/langchaingo/llms"
)

const name = "DeepSummaryForSearchId"

type FunctionCall struct {
	User *types.User
	Ctx  context.Context
}

func CreateNewFunctionCall(c context.Context) *FunctionCall {
	user := c.Value(types.CtxKey("user")).(types.User)
	return &FunctionCall{
		User: &user,
		Ctx:  c,
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
			Description: "Returns the 'deep' summary based on a search id (a repository, directory, file or entity (An entity can be a variable or a function, is a block within a go file",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"searchId": map[string]any{
						"type":        "string",
						"description": "The search id to get the summary for",
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

// const responseTemplate = `
// The name of the function, struct, variable, const or type is is {{.Name}} with dbid {{.Dbid}}.
// The signature is: {{.Signature}}
// The summary is: {{.Summary}}

// Information regarding the file containing the entity:
// Import path: {{.ImportPath}}
// File dbid: {{.FileDbid}}
// File summary: {{.FileSummary}}

// Information regarding the Package containing the file:
// Package path: {{.PackageImportPath}}
// Package dbid: {{.PackageDbid}}
// Package summary: {{.PackageShortSummary}}

// Information regarding the Repository containing the package:
// Repository name: {{.RepoName}}
// Repository dbid: {{.RepoDbid}}
// Repository summary: {{.RepoShortSummary}}
// `

func (f *FunctionCall) Function(searchId string) string {
	words := strings.Split(searchId, "-")
	typeOfSearch := words[0]
	switch typeOfSearch {
	case "repository":
		fc := repositorydeepsummaryfc.CreateNewFunctionCall(f.Ctx)
		returnString := fc.Function(searchId)
		return returnString
	case "directory":
		fc := directorydeepsummaryfc.CreateNewFunctionCall(f.Ctx)
		returnString := fc.Function(searchId)
		return returnString
	case "file":
		fc := filedeepsummaryfc.CreateNewFunctionCall(f.Ctx)
		returnString := fc.Function(searchId)
		return returnString
	case "entity":
		fc := entitydeepsummaryfc.CreateNewFunctionCall(f.Ctx)
		returnString := fc.Function(searchId)
		return returnString
	default:
		return "The search id provided is not a valid search id"
	}

	// requestedType := search.GetTypeFromSearchId(entitySearchId)
	// if requestedType != "entity" {
	// 	return "The entity dbid provided does not correspond to an entity, but to a " + requestedType
	// }
	// entityDbid := search.GetDbidFromSearchId(entitySearchId)
	// repoResult, err := cypherqueries.PerformEntityCypherQuery(entityDbid)
	// if err != nil {
	// 	return "An neo4j error occurred while performing the query: " + err.Error()
	// }
	// tmpl, err := template.New("systemMessage").Parse(responseTemplate)
	// if err != nil {
	// 	return fmt.Sprintf("Error in template.New: %v", err)
	// }

	// var result bytes.Buffer
	// err = tmpl.Execute(&result, repoResult)
	// if err != nil {
	// 	return fmt.Sprintf("Error in tmpl.Execute: %v", err)
	// }

	// response := result.String()
	// return response
}
