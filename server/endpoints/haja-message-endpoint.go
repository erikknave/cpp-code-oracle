package endpoints

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/erikknave/go-code-oracle/agents/codebaseagent"
	"github.com/erikknave/go-code-oracle/dbhelpers"
	"github.com/erikknave/go-code-oracle/filecontent"
	"github.com/erikknave/go-code-oracle/types"
	"github.com/gofiber/fiber/v2"
)

func HajaMessageEndpoint(c *fiber.Ctx) error {
	if os.Getenv("HAJA_AGENT_TOOL_KEY") != c.Get("Authorization") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	promptStr := c.FormValue("message")
	username := c.FormValue("thread_id")
	apiVersion := c.FormValue("api_version")
	fmt.Println("API Version:", apiVersion)
	fmt.Printf("Username/thread_id: %s\n", username)
	user, err := dbhelpers.LoadUserFromUserName(username)
	if err != nil || user.ID == 0 {
		user, err = dbhelpers.CreateUser(username)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error creating user"})
		}
	}
	if promptStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Message is required"})
	}
	if promptStr == "/clear" {
		dbhelpers.ClearChatMessagesForUser(&user)
		return c.JSON(fiber.Map{"message": "Chat cleared"})
	}
	messageHistory, _ := dbhelpers.LoadChatMessagesForUser(&user)
	ctx := context.WithValue(context.WithValue(context.Background(), types.CtxKey("user"), user), types.CtxKey("prompt"), promptStr)
	response, tCtx, err := sendCodeBaseAgentMessage(messageHistory, user, ctx, promptStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error sending message"})
	}
	fileContents := []types.FileContent{}
	for _, filePath := range tCtx.MentionedFiles {
		content, err := filecontent.GetFileContent(filePath)
		if err != nil {
			log.Println("Error getting file content:", err)
			continue
		}
		fileContents = append(fileContents, types.FileContent{
			FilePath: filePath,
			Content:  content,
		})
	}
	if apiVersion == "" {
		return c.Status(fiber.StatusOK).SendString(response)
	}

	return c.JSON(fiber.Map{
		"response": response,
		"context":  tCtx,
		"files":    fileContents,
	})
}

func sendCodeBaseAgentMessage(
	messageHistory []types.ChatMessage,
	user types.User,
	c context.Context,
	promptStr string,
) (string, types.ToolContext, error) {
	agent := &codebaseagent.Agent{}
	agent.Init(messageHistory, &user, c)
	var messages []types.ChatMessage
	_, messages, err := agent.Invoke(promptStr, messageHistory, &user)
	if err != nil {
		log.Println("Error invoking code base agent:", err)
	}
	messages = dbhelpers.SetChatMessages(messages)
	finalMessage := messages[len(messages)-1]
	toolContext := types.ToolContext{}
	err = json.Unmarshal([]byte(finalMessage.Context), &toolContext)
	if err != nil {
		log.Println("Error unmarshalling tool context:", err)
	}
	// c = context.WithValue(c, types.CtxKey("chatMessages"), messages)
	return finalMessage.Content, toolContext, nil
}
