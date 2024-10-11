package endpoints

import (
	_ "embed"
	"os"

	"github.com/gofiber/fiber/v2"
)

func HajaMessageEndpoint(c *fiber.Ctx) error {
	if os.Getenv("HAJA_AGENT_TOOL_KEY") != c.Get("Authorization") {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}
	param := c.Query("message")
	return c.SendString("Query parameter: " + param)
	// return c.SendString("POST request authorized")
}
