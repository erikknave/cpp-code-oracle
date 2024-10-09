package directoryagent

import (
	"context"
	"fmt"

	"github.com/erikknave/go-code-oracle/agents"
	"github.com/erikknave/go-code-oracle/agents/functioncalls/deepsummaryfc"
	"github.com/erikknave/go-code-oracle/agents/functioncalls/filecommitsfcdbid"
	"github.com/erikknave/go-code-oracle/agents/functioncalls/fileshortsummaryfc"
	"github.com/erikknave/go-code-oracle/agents/functioncalls/getfilecontentsfc"
	"github.com/erikknave/go-code-oracle/agents/functioncalls/listfilesforsearchidfc"
	"github.com/erikknave/go-code-oracle/cypher/cypherqueries"
	"github.com/erikknave/go-code-oracle/types"
)

const NAME = "directoryAgent"

type queryResponseType struct {
	Name           string   `json:"name"`
	Path           string   `json:"path"`
	RepositoryName string   `json:"repositoryName"`
	Dbid           string   `json:"dbid"`
	Summary        string   `json:"summary"`
	Files          []string `json:"files"`
}

// const systemMsgQueryTemplate = `
// MATCH (r1:Repository )-[:HAS_MODULE]->(m1:Module)<-[:PART_OF_MODULE]-(p1:Package{dbid:"%d"} )-[:CONTAINS]->(f1:File )
// where f1.name <> "NON_EXISTING_FILE.go" with r1,p1, collect(f1.repoPath) as files
// return {name:p1.name, path:p1.repoPath, repositoryName:r1.name, summary: p1.summary, dbid: p1.dbid, files:files} as result
// `

type Agent struct {
	agents.BaseAgent
}

func (a *Agent) Init(messageHistory []types.ChatMessage, u *types.User, dbid int, c context.Context) {
	availableTools := []types.FunctionCall{
		fileshortsummaryfc.CreateNewFunctionCall(u, dbid),
		filecommitsfcdbid.CreateNewFunctionCall(c),
		deepsummaryfc.CreateNewFunctionCall(c),
		listfilesforsearchidfc.CreateNewFunctionCall(c),
		getfilecontentsfc.CreateNewFunctionCall(c),
	}
	endResult, err := cypherqueries.PerformDirectoryCypherQuery(fmt.Sprintf("%d", dbid))
	if err != nil {
		fmt.Printf("%s: Error executing cypher query when creating system message\n", NAME)
	}
	a.InitBaseAgent(NAME, messageHistory, availableTools, u, endResult)
}

func (a *Agent) Invoke(queryString string, messages []types.ChatMessage, u *types.User) (string, []types.ChatMessage, error) {

	templateData := map[string]interface{}{
		"Content": queryString,
	}
	response, messageHistory, err := a.InvokeBaseAgent(templateData, nil, u)
	return response, messageHistory, err
}
