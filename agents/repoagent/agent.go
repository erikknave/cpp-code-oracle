package repoagent

import (
	"context"
	"fmt"

	"github.com/erikknave/go-code-oracle/agents"
	"github.com/erikknave/go-code-oracle/agents/functioncalls/deepsummaryfc"
	"github.com/erikknave/go-code-oracle/agents/functioncalls/directoriesshortsummaryfc"
	"github.com/erikknave/go-code-oracle/agents/functioncalls/filecommitsfcdbid"
	"github.com/erikknave/go-code-oracle/agents/functioncalls/getfilecontentsfc"
	"github.com/erikknave/go-code-oracle/agents/functioncalls/listfilesforsearchidfc"
	"github.com/erikknave/go-code-oracle/cypher/cypherqueries"
	"github.com/erikknave/go-code-oracle/types"
)

const NAME = "repoAgent"

type Agent struct {
	agents.BaseAgent
}

func (a *Agent) Init(messageHistory []types.ChatMessage, u *types.User, dbid int, c context.Context) error {
	availableTools := []types.FunctionCall{
		// &modulesshortsummaryfc.FunctionCall{},
		// &packagesshortsummaryfc.FunctionCall{},
		// repositoriesshortsummaryfc.CreateNewFunctionCall(c),
		filecommitsfcdbid.CreateNewFunctionCall(c),
		deepsummaryfc.CreateNewFunctionCall(c),
		listfilesforsearchidfc.CreateNewFunctionCall(c),
		directoriesshortsummaryfc.CreateNewFunctionCall(u, dbid),
		getfilecontentsfc.CreateNewFunctionCall(c),
	}
	repoResult, err := cypherqueries.PerformRepoCypherQuery(fmt.Sprintf("%d", dbid))
	if err != nil {
		return err
	}
	a.InitBaseAgent(NAME, messageHistory, availableTools, u, repoResult)
	return nil
}

func (a *Agent) Invoke(queryString string, messages []types.ChatMessage, u *types.User) (string, []types.ChatMessage, error) {

	templateData := map[string]interface{}{
		"Content": queryString,
	}
	response, messageHistory, err := a.InvokeBaseAgent(templateData, nil, u)
	return response, messageHistory, err
}
