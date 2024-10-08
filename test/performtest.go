package test

import (
	"fmt"
	"log"
	"time"

	"github.com/erikknave/go-code-oracle/helpers"
	"github.com/erikknave/go-code-oracle/search"
	"github.com/erikknave/go-code-oracle/types"
	"github.com/gofiber/contrib/websocket"
)

func PerformTest(c *websocket.Conn) (types.ChatMessagePacket, error) {
	serverChatMessage := types.ChatMessage{
		Role:    "server",
		Content: "Test performed",
		Date:    time.Now(),
	}
	docs, err := search.SearchDocuments("logging", 10, "type = package")
	if err != nil {
		log.Println("Error searching documents:", err)
	}
	fmt.Printf("Found %d documents\n", len(docs))
	for _, doc := range docs {
		// fmt.Printf("Document: %v\n", doc)
		yamlStr, err := helpers.PrettyPrintYAMLInterface(doc)
		if err != nil {
			log.Println("Error pretty printing YAML:", err)
		}
		fmt.Printf("Document: \n%v\n", yamlStr)

	}
	packet := types.ChatMessagePacket{
		Message:         serverChatMessage,
		UserInputStatus: "can_respond",
		Type:            "ChatMessagePacket",
	}
	err = c.WriteJSON(packet)
	if err != nil {
		log.Println("write:", err)
		return packet, err
	}
	return packet, nil
}
