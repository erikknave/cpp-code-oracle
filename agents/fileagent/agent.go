package fileagent

import (
	"context"
	"fmt"

	"github.com/erikknave/go-code-oracle/agents"
	"github.com/erikknave/go-code-oracle/agents/functioncalls/deepsummaryfc"
	"github.com/erikknave/go-code-oracle/agents/functioncalls/filecommitsfcdbid"
	"github.com/erikknave/go-code-oracle/agents/functioncalls/getfilecontentsfc"
	"github.com/erikknave/go-code-oracle/cypher/cypherqueries"
	"github.com/erikknave/go-code-oracle/types"
)

const NAME = "fileAgent"

// type queryResponseType struct {
// 	FileName       string   `json:"fileName"`
// 	PackageName    string   `json:"packageName"`
// 	RepositoryName string   `json:"repositoryName"`
// 	Dbid           string   `json:"dbid"`
// 	Summary        string   `json:"summary"`
// 	Entities       []string `json:"entities"`
// }

// const systemMsgQueryTemplate = `
// MATCH (r1:Repository)-[:HAS_MODULE]->(m1:Module)<-[:PART_OF_MODULE]-(p1:Package)-[:CONTAINS]->(f1:File {dbid: "%d"})-[:DEFINES]->(e1:Entity)
// WHERE f1.name <> "NON_EXISTING_FILE.go"
// WITH r1, p1, f1, collect(e1.signature) as entities
// RETURN {
//     repositoryName: r1.name,
//     packageName: p1.repoPath,
//     fileName: f1.repoPath,
//     dbid: f1.dbid,
//     summary: f1.summary,
//     entities: entities
// } AS result
// `

type Agent struct {
	agents.BaseAgent
}

func (a *Agent) Init(messageHistory []types.ChatMessage, u *types.User, dbid int, c context.Context) {
	availableTools := []types.FunctionCall{
		filecommitsfcdbid.CreateNewFunctionCall(c),
		deepsummaryfc.CreateNewFunctionCall(c),
		getfilecontentsfc.CreateNewFunctionCall(c),
	}
	dbidStr := fmt.Sprintf("%d", dbid)
	endResult, err := cypherqueries.PerformFileCypherQuery(dbidStr)
	if err != nil {
		fmt.Printf("%s: Error executing cypher query when creating system message\n", NAME)
	}
	a.InitBaseAgent(NAME, messageHistory, availableTools, u, endResult)
}

func (a *Agent) Invoke(queryString string, messages []types.ChatMessage, u *types.User) (string, []types.ChatMessage, error) {

	templateData := map[string]interface{}{
		"Content": queryString,
	}
	// a.InitBaseAgent(NAME, messages, a.FunctionCalls, u, nil)
	response, messageHistory, err := a.InvokeBaseAgent(templateData, nil, u)
	return response, messageHistory, err
}
