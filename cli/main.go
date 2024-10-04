package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/erikknave/go-code-oracle/cli/connection"
	"github.com/erikknave/go-code-oracle/cli/helpers"
	"github.com/erikknave/go-code-oracle/types"
)

const (
	Reset = "\033[0m"
	Bold  = "\033[1m"
	Blue  = "\033[34m"
	Red   = "\033[31m"
)

// Init is the main entry point for the CLI client.
func main() {
	packets := make(chan json.RawMessage)
	go connection.Init(packets)
	scanner := bufio.NewScanner(os.Stdin)
	var userInputStatus = "can_respond"
	chatMessages, err := helpers.FetchChatMessages("test")
	if err != nil {
		log.Fatal("Failed to fetch chat messages:", err)
	}
	helpers.PrintChatMessages(chatMessages)
	helpers.RequestInput(scanner, packets)
	for packet := range packets {
		var msgMap map[string]interface{}
		err := json.Unmarshal(packet, &msgMap)
		if err != nil {
			log.Print("Failed to unmarshal raw message:", err)
			continue
		}
		switch msgMap["type"] {
		case "ChatMessagePacket":
			chatPacket := types.ChatMessagePacket{}
			err := json.Unmarshal(packet, &chatPacket)
			if err != nil {
				log.Print("Failed to unmarshal chat message packet:", err)
				continue
			}
			helpers.PrintChatMessage(chatPacket.Message)
			userInputStatus = chatPacket.UserInputStatus
		}

		if userInputStatus == "can_respond" {
			helpers.RequestInput(scanner, packets)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}
}
