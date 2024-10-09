package codebaseagent

import (
	"context"

	"github.com/erikknave/go-code-oracle/agents"
	"github.com/erikknave/go-code-oracle/agents/functioncalls/deepsummaryfc"
	"github.com/erikknave/go-code-oracle/agents/functioncalls/directoriesshortsummaryfcdbid"
	"github.com/erikknave/go-code-oracle/agents/functioncalls/filecommitsfcdbid"
	"github.com/erikknave/go-code-oracle/agents/functioncalls/fileswithinrepositoryfc"
	"github.com/erikknave/go-code-oracle/agents/functioncalls/getfilecontentsfc"
	"github.com/erikknave/go-code-oracle/agents/functioncalls/performreposearchfc"
	"github.com/erikknave/go-code-oracle/agents/functioncalls/searchallfilesfc"
	"github.com/erikknave/go-code-oracle/agents/functioncalls/searcheverythingfc"
	"github.com/erikknave/go-code-oracle/cypher/cypherqueries"
	"github.com/erikknave/go-code-oracle/types"
)

const NAME = "codeBaseAgent"

type Agent struct {
	agents.BaseAgent
}

func (a *Agent) Init(messageHistory []types.ChatMessage, u *types.User, c context.Context) {
	availableTools := []types.FunctionCall{
		// askrepoagentfc.CreateNewFunctionCall(c),
		searcheverythingfc.CreateNewFunctionCall(c),
		searchallfilesfc.CreateNewFunctionCall(c),
		deepsummaryfc.CreateNewFunctionCall(c),
		performreposearchfc.CreateNewFunctionCall(c),
		directoriesshortsummaryfcdbid.CreateNewFunctionCall(c),
		fileswithinrepositoryfc.CreateNewFunctionCall(c),
		filecommitsfcdbid.CreateNewFunctionCall(c),
		getfilecontentsfc.CreateNewFunctionCall(c),
		// &modulesshortsummaryfc.FunctionCall{},
		// &packagesshortsummaryfc.FunctionCall{},
		// &repositoriesshortsummaryfc.FunctionCall{},
	}
	repoList, _ := cypherqueries.PerformRepoListCypherQuery()
	a.InitBaseAgent(NAME, messageHistory, availableTools, u, repoList)
}

func (a *Agent) Invoke(queryString string, messages []types.ChatMessage, u *types.User) (string, []types.ChatMessage, error) {

	templateData := map[string]interface{}{
		"Content": queryString,
	}
	// a.InitBaseAgent(NAME, messages, a.FunctionCalls, u, )
	response, messageHistory, err := a.InvokeBaseAgent(templateData, nil, u)
	return response, messageHistory, err
}
