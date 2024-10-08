package endpoints

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/erikknave/go-code-oracle/agents/agentflows"
	"github.com/erikknave/go-code-oracle/dbhelpers"
	"github.com/erikknave/go-code-oracle/server/serverhelpers"
	"github.com/erikknave/go-code-oracle/types"
	"github.com/erikknave/go-code-oracle/web/templates"
	"github.com/erikknave/go-code-oracle/web/webhelpers"
	"github.com/gofiber/fiber/v2"
)

func SendMessageEndPoint(c *fiber.Ctx) error {
	ctx := context.Background()
	promptStr := c.FormValue("prompt")
	userInterface := c.Locals("user")
	user := userInterface.(types.User)
	words := strings.Fields(promptStr)
	if words[0] == "/stats" {
		return StatsViewWrapperEndPoint(c)
	}
	if words[0] == "/search" {
		return PerformSearchEndPoint(c)
	}
	if words[0] == "/help" {
		return HelpViewWrapperEndPoint(c)
	}
	userAgentType, searchableDocument, err := serverhelpers.InterceptAgentType(c, user)
	if err != nil {
		return err
	}

	ctx = context.WithValue(ctx, types.CtxKey("agentType"), userAgentType)
	ctx = context.WithValue(ctx, types.CtxKey("searchableDocument"), searchableDocument)
	ctx = context.WithValue(ctx, types.CtxKey("user"), user)
	ctx = context.WithValue(ctx, types.CtxKey("prompt"), promptStr)
	if strings.Contains(promptStr, "/clear") {
		go SendDeleteAllMessagesToUser(ctx, &user)
	} else if strings.Contains(promptStr, "/chat") {
		c.Request().SetRequestURI(c.Path())
		dbhelpers.SetUserAgentType(&user, "codeBaseAgent", nil)
		dbhelpers.ClearChatMessagesForUser(&user)
		return ChatViewWrapperEndPoint(c)
	} else {
		go StartCodeAgentFlow(ctx)
	}
	if promptStr == "/chat" {
		return PerformSearchEndPoint(c)
	}
	emptyChatPrompt := templates.ChatPrompt(userAgentType, searchableDocument)
	var buf bytes.Buffer
	emptyChatPrompt.Render(c.Context(), &buf)
	compStr := buf.String()
	c.Response().Header.SetContentType("text/html")
	c.SendString(compStr)
	return nil
}

func StartCodeAgentFlow(c context.Context) error {
	ctx := c
	user := ctx.Value(types.CtxKey("user")).(types.User)
	messageHistory, err := dbhelpers.LoadChatMessagesForUser(&user)
	if err != nil {
		messageHistory = []types.ChatMessage{}
	}
	c = context.WithValue(c, types.CtxKey("chatMessages"), messageHistory)
	promptStr := ctx.Value(types.CtxKey("prompt")).(string)
	tmpPromptMsg := types.ChatMessage{
		Role:    "user",
		Content: promptStr,
		Date:    time.Now(),
		User:    user,
	}
	go webhelpers.SendHtmxChatMessageToUser(ctx, &user, tmpPromptMsg)
	userAgentType := ctx.Value(types.CtxKey("agentType")).(types.UserAgentType)
	switch userAgentType.AgentType {
	case "codeBaseAgent":
		return agentflows.StartCodeBaseAgentFlow(messageHistory, user, c, promptStr)
	case "repoAgent":
		return agentflows.StartRepoAgentFlow(messageHistory, user, c, promptStr)
	case "packageAgent":
		return agentflows.StartPackageAgentFlow(messageHistory, user, c, promptStr)
	case "fileAgent":
		return agentflows.StartFileAgentFlow(messageHistory, user, c, promptStr)
	}

	fmt.Printf("Unknown agent type: %s\n", userAgentType.AgentType)
	incorrectAgentMsg := types.ChatMessage{
		Role:    "server",
		Content: fmt.Sprintf("Error: Unknown agent type: %s", userAgentType.AgentType),
		Date:    time.Now(),
		User:    user,
	}
	go webhelpers.SendHtmxChatMessageToUser(ctx, &user, incorrectAgentMsg)
	return nil
}

func SendDeleteAllMessagesToUser(ctx context.Context, user *types.User) {
	dbhelpers.ClearChatMessagesForUser(user)
	component := templates.MsgDeleteAllWrapper()
	var buf bytes.Buffer
	err := component.Render(ctx, &buf)
	if err != nil {
		panic(err)
	}
	compStr := buf.String()
	serverhelpers.SendStringToUser(user.Username, compStr)
}
