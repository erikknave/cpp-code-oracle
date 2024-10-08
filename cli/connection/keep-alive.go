package connection

import (
	"context"
	"time"

	"nhooyr.io/websocket"
)

func KeepAlive(c *websocket.Conn) {
	for {
		// fmt.Printf("Sending keep-alive message")
		err := c.Write(context.Background(), websocket.MessageText, []byte("ping"))
		if err != nil {
			return
		}
		time.Sleep(10 * time.Second)
	}
}
