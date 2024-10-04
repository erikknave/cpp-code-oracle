package askpackageagentfc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/erikknave/go-code-oracle/agents/packageagent"
	"github.com/erikknave/go-code-oracle/helpers"
	"github.com/erikknave/go-code-oracle/search"
	"github.com/erikknave/go-code-oracle/types"
	"github.com/erikknave/go-code-oracle/web/webhelpers"
	"github.com/tmc/langchaingo/llms"
)

const name = "AskPackageAgent"

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
			Description: "Asks a question to the go package agent (specializing and having all information regarding a specific package within the repository). The more information and context the question has, the better the answer.",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"question": map[string]any{
						"type":        "string",
						"description": "The question to ask the package agent.",
					},
					"packageSearchId": map[string]any{
						"type":        "string",
						"description": "The search id of the package the question is related to.",
					},
					// "unit": map[string]any{
					// 	"type": "string",
					// 	"enum": []string{"fahrenheit", "celsius"},
					// },
				},
				"required": []string{"question", "packageSearchId"},
			},
		},
	}
}

func (f FunctionCall) Execute(args json.RawMessage) (string, error) {
	var params struct {
		Question        string `json:"question"`
		PackageSearchId string `json:"packageSearchId"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return "", err
	}
	fmt.Printf("Calling function %s with the following args: %v\n", name, params)

	webhelpers.SendServerMessageToUser(f.Context, &f.User, fmt.Sprintf("A question was sent to the Package agent: %v..\n", helpers.SafeSubstring(params.Question, 150)))
	response := f.Function(params.Question, params.PackageSearchId)
	webhelpers.SendServerMessageToUser(f.Context, &f.User, fmt.Sprintf("The Package agent responded: %v..\n", helpers.SafeSubstring(response, 150)))
	return response, nil

}

// type queryResponseType struct {
// 	Name string `json:"name"`
// 	Dbid string `json:"dbid"`
// }

func (f FunctionCall) Function(question string, searchId string) string {
	agent := &packageagent.Agent{}
	dbidStr := search.GetDbidFromSearchId(searchId)
	dbidInt, err := strconv.Atoi(dbidStr)
	if err != nil {
		return "An error occurred when trying to understand dbid"
	}

	agent.Init(nil, &f.User, dbidInt, f.Context)
	var messages []types.ChatMessage
	_, messages, err = agent.Invoke(question, nil, &f.User)
	if err != nil {
		log.Println("Error invoking code base agent:", err)
	}
	lastMsgContent := messages[len(messages)-1].Content
	return lastMsgContent
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
