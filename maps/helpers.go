package maps

import (
	"github.com/erikknave/go-code-oracle/types"
	"github.com/gofiber/contrib/websocket"
)

func AddUserConnection(conn *websocket.Conn) {
	user := conn.Locals("user").(types.User)
	userName := user.Username
	existingConnections, exists := UserWSConnections.Load(userName)

	if !exists {
		existingConnections = []*websocket.Conn{}
	}

	existingConnections = append(existingConnections, conn)
	UserWSConnections.Store(userName, existingConnections)
}

func RemoveUserConnection(conn *websocket.Conn) {
	UserWSConnections.Range(func(key string, value []*websocket.Conn) bool {
		connections := value
		for i, connection := range connections {
			if connection == conn {
				connections = append(connections[:i], connections[i+1:]...)
				UserWSConnections.Store(key, connections)
			}
		}
		return true
	})
}
