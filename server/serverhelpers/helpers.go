package serverhelpers

import (
	"fmt"
	"time"

	"github.com/erikknave/go-code-oracle/dbhelpers"
	"github.com/erikknave/go-code-oracle/maps"
	"github.com/erikknave/go-code-oracle/types"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func SendInitialMessage(c *websocket.Conn) {
	var messages []types.ChatMessage
	var err error
	messages, err = dbhelpers.LoadChatMessagesForWSUser(c)
	if err != nil {
		messages = []types.ChatMessage{}
	}

	message := types.ChatMessage{
		Role:    "server",
		Content: "Welcome to the chat!",
		Date:    time.Now(),
	}
	// var messages []types.ChatMessage
	messages = append([]types.ChatMessage{message}, messages...)
	chatPacket := types.ChatMessagesPacket{
		Messages:        messages,
		Type:            "ChatMessagesPacket",
		UserInputStatus: "can_respond",
	}
	c.WriteJSON(chatPacket)
}

func SendServerMessage(c *websocket.Conn, message string) {
	serverChatMessage := types.ChatMessage{
		Role:    "server",
		Content: message,
		Date:    time.Now(),
	}
	packet := types.ChatMessagePacket{
		Message:         serverChatMessage,
		UserInputStatus: "can_respond",
		Type:            "ChatMessagePacket",
	}
	c.WriteJSON(packet)
}

func SendServerMessageToUser(user *types.User, message string) {
	maps.UserWSConnections.Range(func(key string, value []*websocket.Conn) bool {
		if key == user.Username {
			for _, conn := range value {
				SendServerMessage(conn, message)
			}
		}
		return true
	})
}

func SendStringToUser(userName string, message string) {
	maps.UserWSConnections.Range(func(key string, value []*websocket.Conn) bool {
		if key == userName {
			for _, conn := range value {
				conn.WriteMessage(websocket.TextMessage, []byte(message))
			}
		}
		return true
	})
}

func Cleanup(c *websocket.Conn) {
	maps.RemoveUserConnection(c)
	c.Close()
}

func GetUserFromCookie(c *fiber.Ctx) (types.User, error) {
	var user types.User
	var err error
	username := ""
	username = c.Cookies("username")
	if username == "" {
		username := c.Cookies("CF_Authorization")
		if username == "" {
			return types.User{}, fmt.Errorf("no user found in cookie")
		}
		user, _ = dbhelpers.LoadUserFromUserName(username)
		if user.ID == 0 {
			user, err = dbhelpers.CreateUser(username)
			if err != nil {
				return types.User{}, err
			}
		}
		cookie := new(fiber.Cookie)
		cookie.Name = "username"
		cookie.Value = username                              // replace with your actual username value
		cookie.Expires = time.Now().Add(24 * 31 * time.Hour) // Cookie expires in 24 hours
		c.Cookie(cookie)
	} else {
		user, err = dbhelpers.LoadUserFromUserName(username)
		if err != nil {
			return types.User{}, err
		}
	}
	c.Locals("user", user)
	return user, nil
}
