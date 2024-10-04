package dbhelpers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/erikknave/go-code-oracle/database"
	"github.com/erikknave/go-code-oracle/types"
	"github.com/gofiber/contrib/websocket"
)

func SetChatMessagesForUser(c *websocket.Conn, messages []types.ChatMessage) []types.ChatMessage {
	database.DB.Save(messages)
	c.Locals("chatMessages", messages)
	return messages
}

func SetChatMessages(messages []types.ChatMessage) []types.ChatMessage {
	database.DB.Save(messages)
	return messages
}

func AddChatMessage(message *types.ChatMessage, c *websocket.Conn) []types.ChatMessage {
	database.DB.Save(message)
	userChatMessages := c.Locals("chatMessages").([]types.ChatMessage)
	userChatMessages = append(userChatMessages, *message)
	c.Locals("chatMessages", userChatMessages)
	return userChatMessages
}

func SetUserSearchResults(user *types.User, searchResults []types.SearchableDocument) []types.SearchableDocument {
	resultsJsonStr, err := json.Marshal(searchResults)
	if err != nil {
		log.Fatalf("Error marshalling search results: %v", err)
	}
	userSearchResults := types.UserSearchResults{
		UserID:  user.ID,
		Results: string(resultsJsonStr),
	}
	database.DB.Unscoped().Where("user_id = ?", user.ID).Delete(&types.UserSearchResults{})
	database.DB.Save(&userSearchResults)
	return searchResults
}

func SetUserAgentType(user *types.User, agentType string, searchDoc *types.SearchableDocument) types.UserAgentType {
	if searchDoc == nil {
		searchDoc = &types.SearchableDocument{}
	}
	searchDocBytes, err := json.Marshal(searchDoc)
	if err != nil {
		log.Fatalf("Error marshalling search doc: %v", err)
	}
	searchDocStr := string(searchDocBytes)
	userAgentType := types.UserAgentType{
		UserID:             user.ID,
		AgentType:          agentType,
		SearchableDocument: searchDocStr,
	}
	database.DB.Save(&userAgentType)
	return userAgentType
}

func AddInitialUser() {
	database.Init()
	// var users []types.User
	// database.DB.Unscoped().Where("1=1").Delete(&types.User{})
	// database.DB.Unscoped().Where("1=1").Delete(&types.UserAgentType{})
	// users = []types.User{
	// 	{Username: "erik"},
	// 	{Username: "martin"},
	// 	{Username: "knave"},
	// }
	// database.DB.Create(&users)
	// SetUserAgentType(&users[0], "codeBaseAgent", nil)
	// SetUserAgentType(&users[1], "codeBaseAgent", nil)
	// SetUserAgentType(&users[2], "codeBaseAgent", nil)
}

func CreateUser(userName string) (types.User, error) {
	user, err := LoadUserFromUserName(userName)
	if err == nil {
		return types.User{}, fmt.Errorf("User already exists")
	}
	user = types.User{Username: userName}
	database.DB.Save(&user)
	SetUserAgentType(&user, "codeBaseAgent", nil)
	return user, nil
}
