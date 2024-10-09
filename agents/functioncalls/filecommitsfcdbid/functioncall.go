package filecommitsfcdbid

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/erikknave/go-code-oracle/cypher"
	"github.com/erikknave/go-code-oracle/helpers"
	"github.com/erikknave/go-code-oracle/search"
	"github.com/erikknave/go-code-oracle/types"
	"github.com/tmc/langchaingo/llms"
)

const name = "SearchIdCommits"

type FunctionCall struct {
	User *types.User
}

func CreateNewFunctionCall(c context.Context) *FunctionCall {
	u := c.Value(types.CtxKey("user")).(types.User)
	return &FunctionCall{
		User: &u,
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
			Description: "Returns a summary of the latest commits based on a search-id regardless if it is a repository, package, file or entity. The function will return who updated the code, when it was updated, and a summary of the changes.",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"searchId": map[string]any{
						"type":        "string",
						"description": "The search id of the repository to find ",
					},
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

type queryFileCommitType struct {
	Author         string `json:"author"`
	CommitDate     string `json:"commitDate"`
	Summary        string `json:"summary"`
	FileImportPath string `json:"fileImportPath"`
}

type queryResponseType struct {
	SearchId    string `json:"searchId"`
	FileCommits []queryFileCommitType
}

const queryStringRepo = `
MATCH (r:REPOSITORY {dbid:"%s"})-[]-(d:DIRECTORY)-[]-(f:FILE)-[]-(fc:FILECOMMIT)
with {author:fc.authorName, commitDate:fc.commitDate, summary: fc.summary, dbid: r.dbid, fileImportPath: f.repoPath} as commit 
return collect(commit) as result
`

const queryStringDirectory = `
MATCH (r:REPOSITORY )-[]-(d:DIRECTORY {dbid:"%s"})-[]-(f:FILE)-[]-(fc:FILECOMMIT)
with {author:fc.authorName, commitDate:fc.commitDate, summary: fc.summary, dbid: p.dbid, fileImportPath: f.repoPath} as commit 
return collect(commit) as result
`

const queryStringFile = `
MATCH (r:REPOSITORY )-[]-(d:DIRECTORY)-[]-(f:FILE  {dbid:"%s"})-[]-(fc:FILECOMMIT)
with {author:fc.authorName, commitDate:fc.commitDate, summary: fc.summary, dbid: f.dbid, fileImportPath: f.repoPath} as commit 
return collect(commit) as result
`

const responseTemplate = `
Here are the commits for the search id {{.SearchId}}:
{{range .FileCommits}}
- Name: {{.Author}}
- Commit summary: {{.Summary}}
- Commit Date: {{.CommitDate}}
- Updated file import path: {{.FileImportPath}}

{{end}}
`

func (f *FunctionCall) Function(searchId string) string {
	words := strings.Split(searchId, "-")
	typeOfSearch := words[0]
	switch typeOfSearch {
	case "repository":
		return repositoryFileCommit(searchId)
	case "directory":
		return directoryFileCommit(searchId)
	case "file":
		return fileFileCommit(searchId)

	}
	return fmt.Sprintf("The search id %s is not valid", searchId)
}

func repositoryFileCommit(searchId string) string {
	words := strings.Split(searchId, "-")
	dbid := words[1]
	queryString := fmt.Sprintf(queryStringRepo, dbid)
	cypherResult := cypher.InjectCypher(queryString)
	cypherResultJson, err := json.Marshal(cypherResult)
	if err != nil {
		return fmt.Sprintf("error in json.Marshal: %v", err)
	}

	var typedResult [][]queryFileCommitType
	err = json.Unmarshal(cypherResultJson, &typedResult)
	if err != nil {
		return fmt.Sprintf("error in json.Unmarshal: %v", err)
	}

	if len(typedResult) == 0 {
		return fmt.Sprintf("unexpected result format: %v", typedResult)
	}
	finalResult := typedResult[0]
	responseResult := queryResponseType{
		SearchId:    "repository-" + dbid,
		FileCommits: finalResult,
	}

	summaryString, err := helpers.CreateStringFromTemplate(responseTemplate, responseResult)
	if err != nil {
		return fmt.Sprintf("Error in helpers.CreateStringFromTemplate: %v", err)
	}
	return summaryString
}

func directoryFileCommit(searchId string) string {
	words := strings.Split(searchId, "-")
	dbid := words[1]
	queryString := fmt.Sprintf(queryStringDirectory, dbid)
	cypherResult := cypher.InjectCypher(queryString)
	cypherResultJson, err := json.Marshal(cypherResult)
	if err != nil {
		return fmt.Sprintf("error in json.Marshal: %v", err)
	}

	var typedResult [][]queryFileCommitType
	err = json.Unmarshal(cypherResultJson, &typedResult)
	if err != nil {
		return fmt.Sprintf("error in json.Unmarshal: %v", err)
	}

	if len(typedResult) == 0 {
		return fmt.Sprintf("unexpected result format: %v", typedResult)
	}
	finalResult := typedResult[0]
	responseResult := queryResponseType{
		SearchId:    "package-" + dbid,
		FileCommits: finalResult,
	}
	if len(finalResult) == 0 {
		repoSearchId, err := search.GetRepoSearchIdFromSearchId("package-" + dbid)
		if err != nil {
			return fmt.Sprintf("Error in search.GetRepoSearchIdFromSearchId: %v", err)
		}
		repoSummaryString := repositoryFileCommit(repoSearchId)
		return fmt.Sprintf("No commits found for package %s, however there might have been done changes to other parts of the same repository, in such case they are presented below:\n%s", dbid, repoSummaryString)
	}
	summaryString, err := helpers.CreateStringFromTemplate(responseTemplate, responseResult)
	if err != nil {
		return fmt.Sprintf("Error in helpers.CreateStringFromTemplate: %v", err)
	}
	return fmt.Sprintf("%s", summaryString)
}

func fileFileCommit(searchId string) string {
	words := strings.Split(searchId, "-")
	dbid := words[1]
	queryString := fmt.Sprintf(queryStringFile, dbid)
	cypherResult := cypher.InjectCypher(queryString)
	cypherResultJson, err := json.Marshal(cypherResult)
	if err != nil {
		return fmt.Sprintf("error in json.Marshal: %v", err)
	}

	var typedResult [][]queryFileCommitType
	err = json.Unmarshal(cypherResultJson, &typedResult)
	if err != nil {
		return fmt.Sprintf("error in json.Unmarshal: %v", err)
	}

	if len(typedResult) == 0 {
		return fmt.Sprintf("unexpected result format: %v", typedResult)
	}
	finalResult := typedResult[0]
	responseResult := queryResponseType{
		SearchId:    "file-" + dbid,
		FileCommits: finalResult,
	}
	if len(finalResult) == 0 {
		repoSearchId, err := search.GetRepoSearchIdFromSearchId(searchId)
		if err != nil {
			return fmt.Sprintf("Error in search.GetRepoSearchIdFromSearchId: %v", err)
		}
		repoSummaryString := repositoryFileCommit(repoSearchId)
		return fmt.Sprintf("No commits found for file %s, however there might have been done changes to other parts of the same repository, in such case they are presented below:\n%s", dbid, repoSummaryString)
	}
	summaryString, err := helpers.CreateStringFromTemplate(responseTemplate, responseResult)
	if err != nil {
		return fmt.Sprintf("Error in helpers.CreateStringFromTemplate: %v", err)
	}
	return fmt.Sprintf("%s", summaryString)
}

func entityFileCommit(searchId string) string {
	docs, err := search.SearchAllDocuments(searchId, 1)
	if err != nil {
		return fmt.Sprintf("entityFileCommit: Error in search.SearchAllDocuments: %v", err)
	}
	if len(docs) == 0 {
		return fmt.Sprintf("The search id %s is not valid", searchId)
	}
	if docs[0].Type != "entity" {
		return fmt.Sprintf("entityFileCommit: Error - The search id %s is not an entity", searchId)
	}
	fileSearchId := fmt.Sprintf("file-%d", docs[0].FileID)
	return fileFileCommit(fileSearchId)

}
