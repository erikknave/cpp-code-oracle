package endpoints

import (
	"context"
	_ "embed"
	"log"
	"os"

	"github.com/erikknave/go-code-oracle/agents/codebaseagent"
	"github.com/erikknave/go-code-oracle/dbhelpers"
	"github.com/erikknave/go-code-oracle/types"
	"github.com/gofiber/fiber/v2"
)

func HajaMessageEndpoint(c *fiber.Ctx) error {
	if os.Getenv("HAJA_AGENT_TOOL_KEY") != c.Get("Authorization") {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}
	promptStr := c.FormValue("message")
	username := "haja_user"
	user, err := dbhelpers.LoadUserFromUserName(username)
	if err != nil {
		log.Println("Error loading user:", err)
	}
	if user.ID == 0 {
		var err error
		user, err = dbhelpers.CreateUser(username)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error creating user")
		}
	}
	c.Locals("user", user)
	if promptStr == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Message is required")
	}
	if promptStr == "/clear" {
		dbhelpers.ClearChatMessagesForUser(&user)
		return c.SendString("Chat cleared")
	}
	messageHistory, err := dbhelpers.LoadChatMessagesForUser(&user)
	if err != nil {
		log.Println("Error loading chat messages:", err)
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, types.CtxKey("user"), user)
	ctx = context.WithValue(ctx, types.CtxKey("prompt"), promptStr)
	response, err := sendCodeBaseAgentMessage(messageHistory, user, ctx, promptStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error sending message")
	}
	return c.SendString(response)
}

func sendCodeBaseAgentMessage(
	messageHistory []types.ChatMessage,
	user types.User,
	c context.Context,
	promptStr string,
) (string, error) {
	agent := &codebaseagent.Agent{}
	agent.Init(messageHistory, &user, c)
	var messages []types.ChatMessage
	_, messages, err := agent.Invoke(promptStr, messageHistory, &user)
	if err != nil {
		log.Println("Error invoking code base agent:", err)
	}
	messages = dbhelpers.SetChatMessages(messages)
	// c = context.WithValue(c, types.CtxKey("chatMessages"), messages)
	return messages[len(messages)-1].Content, nil
}
