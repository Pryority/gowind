// cmd/server/handlers/home.go
package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// LoadMain renders the "index" template for the home page.
func LoadMain(c *fiber.Ctx) error {
	return c.Render("index", nil)
}
