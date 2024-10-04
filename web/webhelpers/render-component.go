package webhelpers

import (
	"bytes"
	"context"
	"fmt"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

func RenderHttpComponent(component templ.Component, c *fiber.Ctx, ctx context.Context) error {
	var buf bytes.Buffer
	if err := component.Render(ctx, &buf); err != nil {
		return fmt.Errorf("failed to render component: %w", err)
	}
	compStr := buf.String()
	c.Response().Header.SetContentType("text/html")
	if err := c.SendString(compStr); err != nil {
		return fmt.Errorf("failed to send response: %w", err)
	}
	return nil
}
