package askfileagentfc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/erikknave/go-code-oracle/agents/fileagent"
	"github.com/erikknave/go-code-oracle/helpers"
	"github.com/erikknave/go-code-oracle/search"
	"github.com/erikknave/go-code-oracle/types"
	"github.com/erikknave/go-code-oracle/web/webhelpers"
	"github.com/tmc/langchaingo/llms"
)

const name = "AskFileAgent"

type FunctionCall struct {
	User    *types.User
	Dbid    int
	Context context.Context
}

func CreateNewFunctionCall(user *types.User, dbid int, c context.Context) *FunctionCall {

	return &FunctionCall{
		User:    user,
		Dbid:    dbid,
		Context: c,
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
			Description: "Asks a question to the go file agent specializing and having all information regarding ONE specific file within a package, including its entities (i.e. variables, functions, consts, types and methods)). The more information and context the question has, the better the answer.",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"question": map[string]any{
						"type":        "string",
						"description": "The question to ask the file agent.",
					},
					"fileSearchId": map[string]any{
						"type":        "string",
						"description": "The search id of the File the question is related to.",
					},
					// "unit": map[string]any{
					// 	"type": "string",
					// 	"enum": []string{"fahrenheit", "celsius"},
					// },
				},
				"required": []string{"question", "fileSearchId"},
			},
		},
	}
}

func (f FunctionCall) Execute(args json.RawMessage) (string, error) {
	var params struct {
		Question     string `json:"question"`
		FileSearchId string `json:"fileSearchId"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return "", err
	}
	fmt.Printf("Calling function %s with the following args: %v\n", name, params)
	webhelpers.SendServerMessageToUser(f.Context, f.User, fmt.Sprintf("A question was sent to the File agent: %v..\n", helpers.SafeSubstring(params.Question, 150)))
	response := f.Function(params.Question, params.FileSearchId)
	webhelpers.SendServerMessageToUser(f.Context, f.User, fmt.Sprintf("The File agent responded: %v..\n", helpers.SafeSubstring(response, 150)))
	return response, nil

}

// type queryResponseType struct {
// 	Name string `json:"name"`
// 	Dbid string `json:"dbid"`
// }

func (f FunctionCall) Function(question string, searchId string) string {
	agent := &fileagent.Agent{}
	dbidStr := search.GetDbidFromSearchId(searchId)
	dbidInt, err := strconv.Atoi(dbidStr)
	if err != nil {
		return "An error occurred when trying to understand dbid"
	}

	agent.Init(nil, f.User, dbidInt, f.Context)
	var messages []types.ChatMessage
	_, messages, err = agent.Invoke(question, nil, f.User)
	if err != nil {
		log.Println("Error invoking code base agent:", err)
	}
	lastMsgContent := messages[len(messages)-1].Content
	return lastMsgContent
}
