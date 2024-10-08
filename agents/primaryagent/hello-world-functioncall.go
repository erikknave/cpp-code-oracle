package primaryagent

import (
	"encoding/json"
	"fmt"

	"github.com/tmc/langchaingo/llms"
)

const name = "HelloWorldFunction"

type HelloWorldFunctionCall struct{}

func (f *HelloWorldFunctionCall) Name() string {
	return name
}

func (f *HelloWorldFunctionCall) ToolDefinition() llms.Tool {
	return llms.Tool{
		Type: "function",
		Function: &llms.FunctionDefinition{
			Name:        name,
			Description: "Says Hello World to a certain person and return a response",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"name": map[string]any{
						"type":        "string",
						"description": "The name of the person saying hello to",
					},
					// "unit": map[string]any{
					// 	"type": "string",
					// 	"enum": []string{"fahrenheit", "celsius"},
					// },
				},
				"required": []string{"name"},
			},
		},
	}
}

func (f *HelloWorldFunctionCall) Execute(args json.RawMessage) (string, error) {
	var params struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return "", err
	}
	return HelloWorldFunction(params.Name), nil
}

func HelloWorldFunction(name string) string {
	fmt.Printf("HelloWorldFunction called with argument: name: " + name)
	return "Hello, world right back at you, you happy camper responded " + name
}
