package getfilecontentsfc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/erikknave/go-code-oracle/filecontent"
	"github.com/erikknave/go-code-oracle/search"
	"github.com/erikknave/go-code-oracle/types"
	"github.com/tmc/langchaingo/llms"
)

const name = "GetFileContents"

type FunctionCall struct {
	User *types.User
}

func CreateNewFunctionCall(c context.Context) *FunctionCall {
	user := c.Value(types.CtxKey("user")).(types.User)
	return &FunctionCall{
		User: &user,
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
			Description: "Returns the contents of a file based on a search id",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"fileSearchId": map[string]any{
						"type":        "string",
						"description": "The search id of the actual file from which to get the contents",
					},
					// "query": map[string]any{
					// 	"type":        "string",
					// 	"description": "The query to find the files within the repository",
					// },
					// "unit": map[string]any{
					// 	"type": "string",
					// 	"enum": []string{"fahrenheit", "celsius"},
					// },
				},
				"required": []string{"fileSearchId"},
			},
		},
	}
}

func (f *FunctionCall) Execute(args json.RawMessage, tCtx *types.ToolContext) (string, error) {
	fmt.Printf("\n - Execute function %s called\n", name)
	var params struct {
		FileSearchId string `json:"fileSearchId"`
		// Query              string `json:"query"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return "", err
	}
	return f.Function(params.FileSearchId), nil
}

func (f *FunctionCall) Function(searchId string) string {
	requestedType := search.GetTypeFromSearchId(searchId)
	if requestedType != "file" {
		return "The file dbid provided does not correspond to a file, but to a " + requestedType
	}
	dbid := search.GetDbidFromSearchId(searchId)
	searchDocs, err := search.SearchAllDocuments(dbid, 1)
	if err != nil {
		return fmt.Sprintf("Error in search.SearchPackages: %v", err)
	}
	if len(searchDocs) == 0 {
		return "No file found with the provided search id"
	}
	filePath := searchDocs[0].Path
	fileContents, err := filecontent.GetFileContent(filePath)
	if err != nil {
		return fmt.Sprintf("Error in filecontent.GetFileContent: %v", err)
	}
	return fileContents
}
