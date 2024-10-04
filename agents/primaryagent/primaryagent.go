package primaryagent

import (
	"context"

	"github.com/erikknave/go-code-oracle/agents"
	"github.com/erikknave/go-code-oracle/agents/functioncalls/repositoriesshortsummaryfc"
	"github.com/erikknave/go-code-oracle/cypher/cypherqueries"
	"github.com/erikknave/go-code-oracle/types"
)

const NAME = "primaryAgent"

type PrimaryAgent struct {
	agents.BaseAgent
}

func (a *PrimaryAgent) Init(messageHistory []types.ChatMessage, u *types.User, c context.Context) {
	availableTools := []types.FunctionCall{
		// &modulesshortsummaryfc.FunctionCall{},
		// &packagesshortsummaryfc.FunctionCall{},
		repositoriesshortsummaryfc.CreateNewFunctionCall(c),
	}
	a.InitBaseAgent(NAME, messageHistory, availableTools, u, nil)
}

func (a *PrimaryAgent) Invoke(queryString string, messages []types.ChatMessage, u *types.User) (string, []types.ChatMessage, error) {

	templateData, err := cypherqueries.PerformRepoListCypherQuery()
	if err != nil {
		return "", nil, err
	}
	a.InitBaseAgent(NAME, messages, a.FunctionCalls, u, nil)
	response, messageHistory, err := a.InvokeBaseAgent(templateData, nil, u)
	return response, messageHistory, err
}
