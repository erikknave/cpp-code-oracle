package dbhelpers

import (
	"encoding/json"
	"log"

	"github.com/erikknave/go-code-oracle/database"
	"github.com/erikknave/go-code-oracle/types"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func LoadChatMessagesForUser(user *types.User) ([]types.ChatMessage, error) {
	var messages []types.ChatMessage
	result := database.DB.Preload("User").Where("user_id = ?", user.ID).Order("date ASC").Find(&messages)
	if result.Error != nil || messages == nil {
		messages = []types.ChatMessage{}
	}
	return messages, nil
}

func LoadChatMessagesForWSUser(c *websocket.Conn) ([]types.ChatMessage, error) {
	user, ok := c.Locals("user").(types.User)
	if !ok {
		log.Fatalf("User not found in context")
	}
	var messages []types.ChatMessage
	result := database.DB.Preload("User").Where("user_id = ?", user.ID).Order("date ASC").Find(&messages)
	if result.Error != nil || messages == nil {
		messages = []types.ChatMessage{}
	}
	c.Locals("chatMessages", messages)
	return messages, nil
}

func LoadUserFromUserName(username string) (types.User, error) {
	var user types.User
	result := database.DB.First(&user, "username = ?", username)
	if result.Error != nil {
		return types.User{}, result.Error
	}
	return user, nil
}

func LoadUserFromHttpGetQuery(c *fiber.Ctx) (types.User, error) {
	var user types.User
	result := database.DB.First(&user, "username = ?", c.Query("user"))
	if result.Error != nil {
		return types.User{}, result.Error
	}
	c.Locals("user", user)
	return user, nil
}

func LoadUserFromWSGetQuery(c *websocket.Conn) (types.User, error) {
	var user types.User
	result := database.DB.First(&user, "username = ?", c.Query("user"))
	if result.Error != nil {
		return types.User{}, result.Error
	}
	c.Locals("user", user)
	return user, nil
}

func LoadUserSearchResults(user *types.User) ([]types.SearchableDocument, error) {
	var results types.UserSearchResults
	result := database.DB.Where("user_id = ?", user.ID).Find(&results)
	if result.Error != nil {
		return []types.SearchableDocument{}, result.Error
	}
	if results.ID == 0 {
		return []types.SearchableDocument{}, nil
	}
	var searchResults []types.SearchableDocument
	err := json.Unmarshal([]byte(results.Results), &searchResults)
	if err != nil {
		log.Fatalf("Error unmarshalling search results: %v", err)
	}
	return searchResults, nil
}

func LoadUserAgentType(user *types.User) (types.UserAgentType, types.SearchableDocument, error) {
	var userAgentType types.UserAgentType
	var searchableDocument types.SearchableDocument
	result := database.DB.First(&userAgentType, "user_id = ?", user.ID)
	if result.Error != nil {
		return types.UserAgentType{}, types.SearchableDocument{}, result.Error
	}
	searchableDocumentStr := userAgentType.SearchableDocument
	err := json.Unmarshal([]byte(searchableDocumentStr), &searchableDocument)
	if err != nil {
		log.Fatalf("LoadUserAgentType, error unmarshaling searchable document,%s", err)
	}
	return userAgentType, searchableDocument, nil

}
