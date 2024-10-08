package maps

import (
	"github.com/erikknave/go-code-oracle/types"
	"github.com/gofiber/contrib/websocket"
)

var UserChatSessions = SafeMap[string, []types.ChatMessage]{}
var AgentDescriptions = SafeMap[string, types.AgentDescription]{}
var UserWSConnections = SafeMap[string, []*websocket.Conn]{}
