package endpoints

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/erikknave/go-code-oracle/dbhelpers"
	"github.com/erikknave/go-code-oracle/web/templates"
	"github.com/erikknave/go-code-oracle/web/webhelpers"
	"github.com/gofiber/fiber/v2"
)

func SignupPageEndPoint(c *fiber.Ctx) error {
	ctx := context.Background()
	return webhelpers.RenderHttpComponent(templates.SignupPage(""), c, ctx)
}

func SubmitSignupEndPoint(c *fiber.Ctx) error {
	ctx := context.Background()
	userName := c.FormValue("username")
	user, err := dbhelpers.CreateUser(userName)
	if err != nil {
		return webhelpers.RenderHttpComponent(templates.UpdatedSignupView(fmt.Sprintf("User %s already exist", userName)), c, ctx)
	}
	return webhelpers.RenderHttpComponent(templates.LoginView(fmt.Sprintf("User %s created", user.Username)), c, ctx)
}
