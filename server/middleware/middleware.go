package middleware

import (
	"github.com/erikknave/go-code-oracle/server/serverhelpers"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(ctx *fiber.Ctx) error {
	path := ctx.Path()
	nonProtectedPaths := []string{"/login", "/signup", "/haja-message"}
	for _, nonProtectedPath := range nonProtectedPaths {
		if path == nonProtectedPath {
			return ctx.Next()
		}
	}
	_, err := serverhelpers.GetUserFromCookie(ctx)
	if err != nil {
		return ctx.Redirect("/login", 302)
	}
	return ctx.Next()
}
