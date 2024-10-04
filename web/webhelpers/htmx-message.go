package webhelpers

import (
	"bytes"
	"context"
	"time"

	"github.com/erikknave/go-code-oracle/server/serverhelpers"
	"github.com/erikknave/go-code-oracle/types"
	"github.com/erikknave/go-code-oracle/web/templates"
)

func SendHtmxChatMessageToUser(ctx context.Context, user *types.User, message types.ChatMessage) {
	component := templates.MsgUpdateWrapper(message)
	var buf bytes.Buffer
	err := component.Render(ctx, &buf)
	if err != nil {
		panic(err)
	}
	compStr := buf.String()
	serverhelpers.SendStringToUser(user.Username, compStr)
}

func SendServerMessageToUser(ctx context.Context, user *types.User, messageString string) {
	serverMessage := types.ChatMessage{
		User:    *user,
		Content: messageString,
		Role:    "server",
		Date:    time.Now(),
	}

	SendHtmxChatMessageToUser(ctx, user, serverMessage)
}
