package dbhelpers

import (
	"log"
	"time"

	"github.com/erikknave/go-code-oracle/database"
	"github.com/erikknave/go-code-oracle/types"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func ClearChatMessagesForWSUser(c *websocket.Conn) (types.ChatMessagePacket, error) {
	c.Locals("chatMessages", []types.ChatMessage{})
	database.DB.Unscoped().Where("user_id = ?", c.Locals("user").(*types.User).ID).Delete(&types.ChatMessage{})
	serverChatMessage := types.ChatMessage{
		Role:    "server",
		Content: "Chat history cleared",
		Date:    time.Now(),
	}
	packet := types.ChatMessagePacket{
		Message:         serverChatMessage,
		UserInputStatus: "can_respond",
		Type:            "ChatMessagePacket",
	}
	err := c.WriteJSON(packet)
	if err != nil {
		log.Println("write:", err)
		return packet, err

	}
	return packet, nil
}

func ClearChatMessagesForHttpUser(c *fiber.Ctx) {
	c.Locals("chatMessages", []types.ChatMessage{})
	database.DB.Unscoped().Where("user_id = ?", c.Locals("user").(*types.User).ID).Delete(&types.ChatMessage{})
}

func ClearChatMessagesForUser(u *types.User) {
	database.DB.Unscoped().Where("user_id = ?", u.ID).Delete(&types.ChatMessage{})
}
