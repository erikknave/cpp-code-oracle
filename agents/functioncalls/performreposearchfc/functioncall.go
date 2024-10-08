package performreposearchfc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/erikknave/go-code-oracle/helpers"
	"github.com/erikknave/go-code-oracle/search"
	"github.com/erikknave/go-code-oracle/types"
	"github.com/tmc/langchaingo/llms"
)

const name = "PerformElasticSearch"

type FunctionCall struct {
	AgentType          types.UserAgentType
	SearchableDocument types.SearchableDocument
	User               types.User
	Context            context.Context
}

func CreateNewFunctionCall(c context.Context) FunctionCall {
	agentType := c.Value(types.CtxKey("agentType")).(types.UserAgentType)
	searchableDocument := c.Value(types.CtxKey("searchableDocument")).(types.SearchableDocument)
	user := c.Value(types.CtxKey("user")).(types.User)

	return FunctionCall{
		AgentType:          agentType,
		SearchableDocument: searchableDocument,
		User:               user,
		Context:            c,
	}
}

func (f FunctionCall) Name() string {
	return name
}

func (f FunctionCall) ToolDefinition() llms.Tool {
	return llms.Tool{
		Type: "function",
		Function: &llms.FunctionDefinition{
			Name:        name,
			Description: "Perform an elastic search of all repositories. Responds with the repositories sorted by relevance.",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"searchstring": map[string]any{
						"type":        "string",
						"description": "The searchstring",
					},
					// "dbid": map[string]any{
					// 	"type":        "string",
					// 	"description": "The dbid of the repository the question is related to.",
					// },
					// "unit": map[string]any{
					// 	"type": "string",
					// 	"enum": []string{"fahrenheit", "celsius"},
					// },
				},
				"required": []string{"searchstring"},
			},
		},
	}
}

func (f FunctionCall) Execute(args json.RawMessage) (string, error) {
	var params struct {
		SearchString string `json:"searchstring"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return "", err
	}
	fmt.Printf("Calling function %s with the following args: %v\n", name, params)
	return Function(params.SearchString), nil

}

const responseTemplate = `
The following repositories were found (presented by relevance):
{{range .}}
- Name: {{.Name}}
- Summary: {{.ShortSummary}}
- Search ID: repository-{{.Dbid}}
{{end}}
`

// type queryResponseType struct {
// 	Name string `json:"name"`
// 	Dbid string `json:"dbid"`
// }

func Function(searchString string) string {
	searchDocs, err := search.SearchReporisitories(searchString, 10)
	if err != nil {
		return "Error in search"
	}
	summaryString, err := helpers.CreateStringFromTemplate(responseTemplate, searchDocs)
	if err != nil {
		return fmt.Sprintf("Error in helpers.CreateStringFromTemplate: %v", err)
	}
	return summaryString
	// dbids := chromaclient.PerformRepositoryQuery(queryString, 5)
	// summaryString := cypherhelpers.PrintRepositoryShortSummaries(dbids)
	// return "The following repositories were found:\n" + summaryString

	// var typedQueryResponse []queryResponseType
	// queryString := "MATCH (r:Repository) RETURN {name: r.name, dbid: r.dbid} AS repository LIMIT 2"
	// queryString := "MATCH (r:Repository) RETURN r.name, r.dbid LIMIT 2"

	// queryResponse := cypher.InjectCypher(queryString)
	// queryResponseStr, err := helpers.PrettyPrintYAMLInterface(queryResponse)
	// if err != nil {
	// 	return "Error in responseString"
	// }
	// responseString := "The following test results were returned:\n" + queryResponseStr
	// fmt.Printf("\nresponseString: %s\n", responseString)
	// return responseString
	// queryResponseStr
	// typedQueryResponse := queryResponse.([]queryResponseType)

	// typedQueryResponse := queryResponse.([]cypher.QueryResult)
	// prettyJson, err := json.MarshalIndent(queryResponse, "", "    ")
	// if err != nil {
	// 	fmt.Printf("\nError in prettyJson\n")
	// 	return "Error in prettyJson"
	// }
	// err = json.Unmarshal(prettyJson, &typedQueryResponse)
	// if err != nil {
	// 	fmt.Printf("\nError in prettyJson\n")
	// 	return "Error in prettyJson"
	// }
	// for _, response := range typedQueryResponse {
	// 	fmt.Printf("Response name: %s\n", response.Name)
	// 	fmt.Printf("Response dbid: %s\n", response.Dbid)
	// }
	// fmt.Printf("\nqueryResponse: %v\n", typedQueryResponse)

	// if err != nil {
	// 	fmt.Printf("\nError in prettyJson\n")
	// 	return "Error in prettyJson"
	// }
	// fmt.Printf("\nqueryResponse: \n%s\n", prettyJson)
	// fmt.Printf("Type of queryResponse: %s\n", reflect.TypeOf(queryResponse))
	// queryResponseStr, ok := queryResponse.(string)
	// if !ok {
	// 	fmt.Printf("\nError in queryResponse\n")
	// 	return "Error in queryResponse"
	// }
	// helpers.PrettyPrintJSON(queryResponseStr)

	// fmt.Printf("\nThe Test function is called\n")

	// dbids := chromaclient.PerformModuleQuery("What modules are related to observability", 5)
	// fmt.Printf("\nDBIDs: %v\n", dbids)
	// return "Test function performed"
}
