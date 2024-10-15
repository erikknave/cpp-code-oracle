package repositoriesshortsummaryfc

import (
	"context"
	"encoding/json"

	"github.com/erikknave/go-code-oracle/cypher/cypherhelpers"
	chromaclient "github.com/erikknave/go-code-oracle/server/chromaclient"
	"github.com/erikknave/go-code-oracle/types"
	"github.com/tmc/langchaingo/llms"
)

const name = "RepositoriessShortSummary"

type FunctionCall struct {
	AgentType          types.UserAgentType
	SearchableDocument types.SearchableDocument
	User               types.User
	Context            context.Context
}

func CreateNewFunctionCall(c context.Context) *FunctionCall {
	agentType := c.Value(types.CtxKey("agentType")).(types.UserAgentType)
	searchableDocument := c.Value(types.CtxKey("searchableDocument")).(types.SearchableDocument)
	user := c.Value(types.CtxKey("user")).(types.User)

	return &FunctionCall{
		AgentType:          agentType,
		SearchableDocument: searchableDocument,
		User:               user,
		Context:            c,
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
			Description: "Returns the short summaries of a number of repositories related to a query",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"query": map[string]any{
						"type":        "string",
						"description": "The query to find the repositories related to",
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
	var params struct {
		Query string `json:"query"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return "", err
	}
	return Function(params.Query), nil

}

// type queryResponseType struct {
// 	Name string `json:"name"`
// 	Dbid string `json:"dbid"`
// }

func Function(queryString string) string {
	searchDocs := chromaclient.PerformRepositoryQuery(queryString, 5)
	var dbids []int
	for _, searchDoc := range searchDocs {
		dbids = append(dbids, searchDoc.Dbid)
	}
	summaryString := cypherhelpers.PrintRepositoryShortSummaries(dbids)
	return "The following repositories were found:\n" + summaryString

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
