package helpers

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/charmbracelet/glamour"
	"github.com/erikknave/go-code-oracle/cli/cliglobals"
	"github.com/erikknave/go-code-oracle/cli/connection"
	"github.com/erikknave/go-code-oracle/types"
)

func SendMessage(msg string, ch chan json.RawMessage) {
	ctx := context.Background()
	err := cliglobals.Conn.Write(ctx, 1, []byte(msg))
	if err != nil {
		connection.Init(ch)
		err = cliglobals.Conn.Write(ctx, 1, []byte(msg))
		if err != nil {
			panic(err)
		}
	}
}

const (
	Reset = "\033[0m"
	Bold  = "\033[1m"
	Blue  = "\033[34m"
	Red   = "\033[31m"
)

func PrintChatMessages(chatMessages []types.ChatMessage) {
	for _, chatMessage := range chatMessages {
		PrintChatMessage(chatMessage)
	}
}

func PrintChatMessage(chatMessage types.ChatMessage) {
	if chatMessage.HideFromUser {
		return
	}
	var roleString string
	switch chatMessage.Role {
	case "user":
		roleString = Red + Bold + "User: " + Reset
		// fmt.Println(Red + Bold + "User: " + Reset + chatMessage.Content + "\n")
	case "assistant":
		roleString = Blue + Bold + "Assistant: " + Reset
		// messageString = Blue + Bold + "Assistant: " + Reset + chatMessage.Content + "\n"
		// fmt.Println(Blue + Bold + "Assistant: " + Reset + chatMessage.Content + "\n")
	case "server":
		roleString = Bold + "Server: " + Reset
		// messageString = Bold + "Server: " + Reset + chatMessage.Content + "\n"
		// fmt.Println("", chatMessage.Content+"\n")
	}
	renderedString, err := glamour.Render(chatMessage.Content, "dark")
	if err != nil {
		fmt.Println("Failed to render message:", err)
	}
	fmt.Println(roleString + "\n" + renderedString)
}

func FetchChatMessages(userName string) ([]types.ChatMessage, error) {
	resp, err := http.Get("http://127.0.0.1:8080/ChatMessages?user=" + userName)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var messages []types.ChatMessage
	err = json.NewDecoder(resp.Body).Decode(&messages)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func RequestInput(scanner *bufio.Scanner, packets chan json.RawMessage) {
	fmt.Print("> ")
	if scanner.Scan() {
		text := scanner.Text()
		if text == "exit" {
			return
		}
		SendMessage(text, packets)
	}

}
