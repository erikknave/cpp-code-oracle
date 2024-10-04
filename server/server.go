package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/erikknave/go-code-oracle/agents/primaryagent"
	"github.com/erikknave/go-code-oracle/dbhelpers"
	"github.com/erikknave/go-code-oracle/server/endpoints"
	"github.com/erikknave/go-code-oracle/server/middleware"
	"github.com/erikknave/go-code-oracle/server/serverhelpers"
	"github.com/erikknave/go-code-oracle/test"
	"github.com/erikknave/go-code-oracle/types"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/redirect"
)

func ServerInit() {
	app := fiber.New()
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(redirect.New(redirect.Config{
		Rules: map[string]string{
			"/": "/stats",
		},
		StatusCode: 301,
	}))
	app.Use(middleware.AuthMiddleware)

	// app.Get("/chat", adaptor.HTTPHandler(templ.Handler(templates.Layout("Erik"))))
	app.Get("/chat", endpoints.ChatPageEndPoint)
	app.Get("/ws/message", websocket.New(endpoints.WSMessage))
	app.Post("/send-message", endpoints.SendMessageEndPoint)
	app.Post("/perform-search", endpoints.PerformSearchEndPoint)
	app.Post("/perform-package-search", endpoints.PerformPackageSearchEndPoint)
	app.Post("/perform-file-search", endpoints.PerformFileSearchEndPoint)
	app.Post("/perform-entity-search", endpoints.PerformEntitySearchEndPoint)
	app.Get("/repository", endpoints.RepositoryPageEndPoint)
	app.Get("/package", endpoints.PackagePageEndPoint)
	app.Get("/file", endpoints.FilePageEndPoint)
	app.Get("/search", endpoints.SearchPageEndPoint)
	app.Get("/stats", endpoints.StatsPageEndPoint)
	app.Post("/send-command", endpoints.CommandEndPoint)
	app.Get("/help", endpoints.HelpPageEndPoint)
	app.Get("/login", endpoints.LoginPageEndPoint)
	app.Post("/login", endpoints.SubmitLoginEndPoint)
	app.Get("/signup", endpoints.SignupPageEndPoint)
	app.Post("/signup", endpoints.SubmitSignupEndPoint)
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {

		context := context.Background()
		// var err error
		// _, err = dbhelpers.LoadUserFromWSGetQuery(c)
		// if err != nil {
		// 	c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Incorect User Name"))
		// 	c.Close()
		// 	return
		// }
		dbhelpers.LoadChatMessagesForWSUser(c)
		// serverhelpers.SendInitialMessage(c)
		defer func() {
			c.Close()
		}()

		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			user := c.Locals("user").(types.User)
			msgString := string(msg)
			if msgString == "/clear" {
				dbhelpers.ClearChatMessagesForUser(&user)
				continue
			}
			if msgString == "/test" {
				test.PerformTest(c)
				continue
			}

			if strings.HasPrefix(msgString, "/search") {
				serverhelpers.PerformWSSearch(msgString, c)
				continue
			}
			messageHistory, err := dbhelpers.LoadChatMessagesForWSUser(c)
			if err != nil {
				messageHistory = []types.ChatMessage{}
			}
			// messageHistory := c.Locals("chatMessages").([]types.ChatMessage)
			primaryAgent := primaryagent.PrimaryAgent{}
			primaryAgent.Init(nil, &user, context)
			var messages []types.ChatMessage
			_, messages, err = primaryAgent.Invoke(msgString, messageHistory, &user)
			if err != nil {
				log.Println("Error invoking primary agent:", err)
				break
			}
			// dbhelpers.AddChatMessage(&userMessage, c)
			// formattedMsgString := fmt.Sprintf("Echo from '%s': %s", user.Username, string(msg))

			log.Printf("recv from %s: %s", user.Username, msg)
			messages = dbhelpers.SetChatMessagesForUser(c, messages)
			// dbhelpers.AddChatMessage(&assistantMessage, c)
			chatPacket := types.ChatMessagePacket{
				Message:         messages[len(messages)-1],
				UserInputStatus: "can_respond",
				Type:            "ChatMessagePacket",
			}
			// err = c.WriteMessage(mt, []byte(formattedMsgString))
			err = c.WriteJSON(chatPacket)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	}))
	listenAddr := os.Getenv("HTTP_LISTEN_ADDR")

	log.Fatal(app.Listen(fmt.Sprintf(":%s", listenAddr)))
}
