package endpoints

import (
	"context"
	"time"

	"github.com/erikknave/go-code-oracle/dbhelpers"
	"github.com/erikknave/go-code-oracle/web/templates"
	"github.com/erikknave/go-code-oracle/web/webhelpers"
	"github.com/gofiber/fiber/v2"
)

func SubmitLoginEndPoint(c *fiber.Ctx) error {
	ctx := context.Background()
	userName := c.FormValue("username")
	user, err := dbhelpers.LoadUserFromUserName(userName)
	if err != nil {
		return webhelpers.RenderHttpComponent(templates.UpdatedLoginView("User not found. Please try again."), c, ctx)
		// return c.SendString(fmt.Sprintf("Error loading user from query: %v\n", err))
	}
	if user.ID == 0 {
		return webhelpers.RenderHttpComponent(templates.UpdatedLoginView("User not found. Please try again."), c, ctx)
	}
	cookie := new(fiber.Cookie)
	cookie.Name = "username"
	cookie.Value = user.Username                         // replace with your actual username value
	cookie.Expires = time.Now().Add(24 * 31 * time.Hour) // Cookie expires in 24 hours
	c.Cookie(cookie)
	return c.Redirect("/", 302)
}
