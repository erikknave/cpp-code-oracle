package connection

import (
	"context"
	"encoding/json"
	"log"

	"github.com/erikknave/go-code-oracle/cli/cliglobals"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func Init(packets chan json.RawMessage) {
	ctx := context.Background()

	conn, _, err := websocket.Dial(ctx, "ws://localhost:8080/ws?user=test", nil)
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	cliglobals.Conn = conn
	defer conn.Close(websocket.StatusInternalError, "Connection closed")
	// go KeepAlive(conn)

	for {
		var rawMessagePacket json.RawMessage
		err := wsjson.Read(ctx, conn, &rawMessagePacket)
		if err != nil {
			log.Print("Failed to read message:", err)
			return
		}
		packets <- rawMessagePacket
	}
}
