package agentflows

import (
	"context"
	"log"

	"github.com/erikknave/go-code-oracle/agents/codebaseagent"
	"github.com/erikknave/go-code-oracle/agents/containeragent"
	"github.com/erikknave/go-code-oracle/agents/directoryagent"
	"github.com/erikknave/go-code-oracle/agents/fileagent"
	"github.com/erikknave/go-code-oracle/agents/repoagent"
	"github.com/erikknave/go-code-oracle/dbhelpers"
	"github.com/erikknave/go-code-oracle/types"
	"github.com/erikknave/go-code-oracle/web/webhelpers"
)

func StartCodeBaseAgentFlow(
	messageHistory []types.ChatMessage,
	user types.User,
	c context.Context,
	promptStr string,
) error {
	agent := &codebaseagent.Agent{}
	agent.Init(messageHistory, &user, c)
	var messages []types.ChatMessage
	_, messages, err := agent.Invoke(promptStr, messageHistory, &user)
	if err != nil {
		log.Println("Error invoking code base agent:", err)
	}
	messages = dbhelpers.SetChatMessages(messages)
	c = context.WithValue(c, types.CtxKey("chatMessages"), messages)
	webhelpers.SendHtmxChatMessageToUser(c, &user, messages[len(messages)-1])
	return nil
}

func StartRepoAgentFlow(
	messageHistory []types.ChatMessage,
	user types.User,
	c context.Context,
	promptStr string,
) error {
	searchDoc := c.Value(types.CtxKey("searchableDocument")).(types.SearchableDocument)
	agent := &repoagent.Agent{}
	agent.Init(messageHistory, &user, searchDoc.Dbid, c)
	var messages []types.ChatMessage
	_, messages, err := agent.Invoke(promptStr, messageHistory, &user)
	if err != nil {
		log.Println("Error invoking code base agent:", err)
	}
	messages = dbhelpers.SetChatMessages(messages)
	c = context.WithValue(c, types.CtxKey("chatMessages"), messages)
	webhelpers.SendHtmxChatMessageToUser(c, &user, messages[len(messages)-1])
	return nil
}

func StartDirectoryAgentFlow(
	messageHistory []types.ChatMessage,
	user types.User,
	c context.Context,
	promptStr string,
) error {
	agent := &directoryagent.Agent{}
	searchDoc := c.Value(types.CtxKey("searchableDocument")).(types.SearchableDocument)
	agent.Init(messageHistory, &user, searchDoc.Dbid, c)
	var messages []types.ChatMessage
	_, messages, err := agent.Invoke(promptStr, messageHistory, &user)
	if err != nil {
		log.Println("Error invoking code base agent:", err)
	}
	messages = dbhelpers.SetChatMessages(messages)
	c = context.WithValue(c, types.CtxKey("chatMessages"), messages)
	webhelpers.SendHtmxChatMessageToUser(c, &user, messages[len(messages)-1])
	return nil
}

func StartFileAgentFlow(
	messageHistory []types.ChatMessage,
	user types.User,
	c context.Context,
	promptStr string,
) error {
	agent := &fileagent.Agent{}
	searchDoc := c.Value(types.CtxKey("searchableDocument")).(types.SearchableDocument)
	agent.Init(messageHistory, &user, searchDoc.Dbid, c)
	var messages []types.ChatMessage
	_, messages, err := agent.Invoke(promptStr, messageHistory, &user)
	if err != nil {
		log.Println("Error invoking code base agent:", err)
	}
	messages = dbhelpers.SetChatMessages(messages)
	c = context.WithValue(c, types.CtxKey("chatMessages"), messages)
	webhelpers.SendHtmxChatMessageToUser(c, &user, messages[len(messages)-1])
	return nil
}

func StartContainerAgentFlow(
	messageHistory []types.ChatMessage,
	user types.User,
	c context.Context,
	promptStr string,
) error {
	agent := &containeragent.Agent{}
	searchDoc := c.Value(types.CtxKey("searchableDocument")).(types.SearchableDocument)
	agent.Init(messageHistory, &user, searchDoc.Dbid, c)
	var messages []types.ChatMessage
	_, messages, err := agent.Invoke(promptStr, messageHistory, &user)
	if err != nil {
		log.Println("Error invoking code base agent:", err)
	}
	messages = dbhelpers.SetChatMessages(messages)
	c = context.WithValue(c, types.CtxKey("chatMessages"), messages)
	webhelpers.SendHtmxChatMessageToUser(c, &user, messages[len(messages)-1])
	return nil
}
