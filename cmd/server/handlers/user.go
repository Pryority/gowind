// cmd/server/handlers/user.go
package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// Login handles the login request and redirects to the home page.
func Login(c *fiber.Ctx) error {
	// Handle user login logic here
	// ...

	// Redirect to the home page after successful login
	return c.Redirect("/", fiber.StatusFound)
}

// Signup handles the user registration request and redirects to the home page.
func Signup(c *fiber.Ctx) error {
	// Handle user signup logic here
	// ...

	// Redirect to the home page after successful signup
	return c.Redirect("/", fiber.StatusFound)
}
