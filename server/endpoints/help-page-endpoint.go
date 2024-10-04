package endpoints

import (
	"context"
	_ "embed"

	"github.com/erikknave/go-code-oracle/web/templates"
	"github.com/erikknave/go-code-oracle/web/webhelpers"
	"github.com/gofiber/fiber/v2"
)

func HelpPageEndPoint(c *fiber.Ctx) error {
	ctx := context.Background()
	return webhelpers.RenderHttpComponent(templates.HelpPage(), c, ctx)
}

func HelpViewWrapperEndPoint(c *fiber.Ctx) error {
	ctx := context.Background()
	return webhelpers.RenderHttpComponent(templates.HelpViewWrapper(), c, ctx)
}
