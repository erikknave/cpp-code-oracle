package agents

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"text/template"
	"time"

	"github.com/erikknave/go-code-oracle/agents/languagemodel"
	"github.com/erikknave/go-code-oracle/maps"
	"github.com/erikknave/go-code-oracle/types"
	"github.com/tmc/langchaingo/llms"
)

type BaseAgent struct {
	SystemMessage    string
	PromptMessage    types.ChatMessage
	Name             string
	AgentDescription types.AgentDescription
	Messages         []types.ChatMessage
	CallingUser      *types.User
	FunctionCalls    []types.FunctionCall
}

func (a *BaseAgent) InitBaseAgent(
	name string,
	messages []types.ChatMessage,
	functionCalls []types.FunctionCall,
	u *types.User,
	systemMsgTemplateData interface{}) {
	if len(messages) > 0 {
		a.Messages = messages
	} else {
		a.Messages = []types.ChatMessage{}
	}
	var ok bool
	a.AgentDescription, ok = maps.AgentDescriptions.Load(name)
	if !ok {
		log.Fatalf("Agent %s not found in agent descriptions", name)
	}
	if u != nil {
		a.CallingUser = u
	}
	a.FunctionCalls = functionCalls
	systemMsgTemplate := a.AgentDescription.SystemMessage
	tmpl, err := template.New("systemMessage").Parse(systemMsgTemplate)
	if err != nil {
		return
	}

	var result bytes.Buffer
	err = tmpl.Execute(&result, systemMsgTemplateData)
	if err != nil {
		return
	}
	println(result.String())

	a.SystemMessage = result.String()
}

func (a *BaseAgent) InvokeBaseAgent(templateData interface{}, messageHistory []types.ChatMessage, u *types.User) (string, []types.ChatMessage, error) {
	if u != nil {
		a.CallingUser = u
	}
	user := a.CallingUser
	if len(messageHistory) > 0 {
		a.Messages = messageHistory
	}
	// a.SystemMessage = a.AgentDescription.SystemMessage
	var toolContext types.ToolContext
	promptTemplate := a.AgentDescription.PromptTemplate
	tmpl, err := template.New("prompt").Parse(promptTemplate)
	if err != nil {
		return "", nil, err
	}

	var result bytes.Buffer
	err = tmpl.Execute(&result, templateData)
	if err != nil {
		return "", nil, err
	}

	prompt := result.String()
	a.Messages = append(a.Messages, types.ChatMessage{
		User:    *user,
		Role:    "user",
		Content: prompt,
		Date:    time.Now(),
	})

	var messages []llms.MessageContent
	messages = append(messages, llms.TextParts(llms.ChatMessageTypeSystem, a.SystemMessage))
	for _, message := range a.Messages {
		switch message.Role {
		case "user":
			messages = append(messages, llms.TextParts(llms.ChatMessageTypeHuman, message.Content))
		case "assistant":
			messages = append(messages, llms.TextParts(llms.ChatMessageTypeAI, message.Content))
		}
	}
	messages = append(messages, llms.TextParts(llms.ChatMessageTypeHuman, prompt))
	var content string

	content, _, a.Messages, err = languagemodel.GenerateResponse(messages, a.FunctionCalls, a.AgentDescription.Model, a.Messages, user, &toolContext)
	if err != nil || len(content) == 0 {
		return "", nil, err
	}
	toolContextBytes, err := json.Marshal(toolContext)
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)

	}
	toolContextString := string(toolContextBytes)

	a.Messages = append(a.Messages, types.ChatMessage{
		User:    *user,
		Role:    "assistant",
		Content: content,
		Date:    time.Now(),
		Context: toolContextString,
	})

	return content, a.Messages, nil
}
