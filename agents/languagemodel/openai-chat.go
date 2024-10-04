package languagemodel

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/erikknave/go-code-oracle/types"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

func GenerateResponse(content []llms.MessageContent, functionCalls []types.FunctionCall, model string, chatMessageHistory []types.ChatMessage, user *types.User) (string, []llms.MessageContent, []types.ChatMessage, error) {
	var err error
	hostMode := os.Getenv("AI_HOST_MODE")
	var llm *openai.LLM

	if hostMode == "azure" {
		llm, err = openai.New(
			openai.WithAPIType(openai.APITypeAzure),
			openai.WithModel(os.Getenv("LANGUAGE_MODEL")),
			openai.WithEmbeddingModel(os.Getenv("LANGUAGE_MODEL")),
			openai.WithBaseURL(os.Getenv("AZURE_BASE_URL")),
			openai.WithToken(os.Getenv("OPENAI_API_KEY")),
		)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		llm, err = openai.New(
			openai.WithAPIType(openai.APITypeOpenAI),
			openai.WithModel(os.Getenv("LANGUAGE_MODEL")),
			openai.WithEmbeddingModel(os.Getenv("LANGUAGE_MODEL")),
			openai.WithBaseURL(os.Getenv("AZURE_BASE_URL")),
			openai.WithToken(os.Getenv("OPENAI_API_KEY")),
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	// llm, err := openai.New(
	// 	openai.WithAPIType(openai.APITypeAzure),
	// 	openai.WithModel(os.Getenv("AZURE_DEPLOYMENT")),
	// 	openai.WithEmbeddingModel(os.Getenv("AZURE_DEPLOYMENT")),
	// 	openai.WithBaseURL(os.Getenv("AZURE_BASE_URL")))

	// llm, err := openai.New(openai.WithModel(model))
	if err != nil {
		return "", nil, nil, err
	}
	ctx := context.Background()

	// completion, err := llm.GenerateContent(ctx, content, llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
	// 	fmt.Print(string(chunk))
	// 	return nil
	// }))
	var availableTools []llms.Tool
	for _, tool := range functionCalls {
		availableTools = append(availableTools, tool.ToolDefinition())
	}
	completion, err := llm.GenerateContent(ctx, content, llms.WithTools(availableTools), llms.WithTemperature(0))
	if err != nil {
		return "", nil, nil, err
	}
	content = updateMessageHistory(content, completion)
	for {
		var executedContent []llms.MessageContent
		var toolExecuted bool
		executedContent, chatMessageHistory, toolExecuted = executeToolCalls(content, functionCalls, completion, chatMessageHistory, user)
		if !toolExecuted {
			break
		}
		completion, err = llm.GenerateContent(ctx, executedContent, llms.WithTools(availableTools), llms.WithTemperature(0))
		if err != nil {
			return "", nil, nil, err
		}
		executedContent = updateMessageHistory(executedContent, completion)
		content = executedContent

	}

	response := completion.Choices[0].Content
	return response, content, chatMessageHistory, nil
}

func updateMessageHistory(messageHistory []llms.MessageContent, resp *llms.ContentResponse) []llms.MessageContent {
	respchoice := resp.Choices[0]

	assistantResponse := llms.TextParts(llms.ChatMessageTypeAI, respchoice.Content)
	for _, tc := range respchoice.ToolCalls {
		assistantResponse.Parts = append(assistantResponse.Parts, tc)
	}
	return append(messageHistory, assistantResponse)
}

func executeToolCalls(messageHistory []llms.MessageContent, tools []types.FunctionCall, resp *llms.ContentResponse, chatMessageHistory []types.ChatMessage, user *types.User) ([]llms.MessageContent, []types.ChatMessage, bool) {
	fmt.Println("Executing", len(resp.Choices[0].ToolCalls), "tool calls")
	fmt.Println("Executing the following tool calls: " + fmt.Sprint(resp.Choices[0].ToolCalls))

	// List of function implementations

	var toolExecuted bool
	toolMap := make(map[string]types.FunctionCall)
	for _, tool := range tools {
		toolMap[tool.Name()] = tool
	}

	for _, toolCall := range resp.Choices[0].ToolCalls {
		tool, exists := toolMap[toolCall.FunctionCall.Name]
		if !exists {
			log.Fatalf("Unsupported tool: %s", toolCall.FunctionCall.Name)
		}
		fmt.Printf(" - Executing tool: %s\n", toolCall.FunctionCall.Name)
		response, err := tool.Execute(json.RawMessage(toolCall.FunctionCall.Arguments))
		toolExecuted = true
		if err != nil {
			log.Fatal(err)
		}

		toolCallResponse := llms.MessageContent{
			Role: llms.ChatMessageTypeTool,
			Parts: []llms.ContentPart{
				llms.ToolCallResponse{
					ToolCallID: toolCall.ID,
					Name:       toolCall.FunctionCall.Name,
					Content:    response,
				},
			},
		}
		chatResponse := "The function " + toolCall.FunctionCall.Name + " was called and responded with: \n" + response + ""
		chatMessageResponse := types.ChatMessage{
			User:         *user,
			Role:         "assistant",
			Content:      chatResponse,
			Date:         time.Now(),
			HideFromUser: true,
		}

		messageHistory = append(messageHistory, toolCallResponse)
		chatMessageHistory = append(chatMessageHistory, chatMessageResponse)
	}

	return messageHistory, chatMessageHistory, toolExecuted
}
