package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

const (
	templateFolder = "./cmd/web/views"
	staticFolder   = "./cmd/web/public"
	localPort      = ":8080"
)

func main() {
	engine := html.New(templateFolder, ".html")

	// Create a new Fiber app
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Serve static files
	app.Static("/public", staticFolder)

	// Define routes
	app.Get("/", loadMain)
	app.Get("/:route_name", loadTemplate)

	log.Printf("Server started at http://localhost%s", localPort)

	// Start the server
	if err := app.Listen(localPort); err != nil {
		log.Fatalf("Error starting the server: %v", err)
	}
}

func loadMain(c *fiber.Ctx) error {
	return c.Render("index", nil)
}

func loadTemplate(c *fiber.Ctx) error {
	routeName := c.Params("route_name")

	// Render the template
	return c.Render("layout.tpl.html", fiber.Map{
		"RouteName": routeName,
	})
}
