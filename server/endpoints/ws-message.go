package endpoints

import (
	"github.com/erikknave/go-code-oracle/maps"
	"github.com/erikknave/go-code-oracle/server/serverhelpers"
	"github.com/erikknave/go-code-oracle/types"
	"github.com/gofiber/contrib/websocket"
)

func WSMessage(c *websocket.Conn) {
	maps.AddUserConnection(c)
	defer serverhelpers.Cleanup(c)
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			break
		}
		user := c.Locals("user").(types.User)
		msgString := string(msg)
		serverhelpers.SendStringToUser(user.Username, msgString)
	}
}
