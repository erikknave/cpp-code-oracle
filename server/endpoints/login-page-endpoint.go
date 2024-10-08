package endpoints

import (
	"context"
	_ "embed"

	"github.com/erikknave/go-code-oracle/web/templates"
	"github.com/erikknave/go-code-oracle/web/webhelpers"
	"github.com/gofiber/fiber/v2"
)

func LoginPageEndPoint(c *fiber.Ctx) error {
	ctx := context.Background()
	return webhelpers.RenderHttpComponent(templates.LoginPage(""), c, ctx)
}
