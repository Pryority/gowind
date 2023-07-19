package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// LoadTemplate is the handler for the catch-all route for templates.
func LoadTemplate(c *fiber.Ctx) error {
	routeName := c.Params("route_name")
	return c.Render("layouts/main", fiber.Map{"RouteName": routeName})
}
