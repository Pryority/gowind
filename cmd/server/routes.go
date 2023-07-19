package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pryority/gowind/cmd/server/handlers"
)

const LocalPort = ":8080"

// SetupRoutes defines all routes and their corresponding handlers.
func SetupRoutes(app *fiber.App) {
	// Serve static files
	app.Static("/public", "./cmd/web/public")

	// Group related routes with subrouter
	api := app.Group("/api")

	// Define routes
	app.Get("/", handlers.LoadMain)
	app.Get("/notes", handlers.LoadNotes)

	// Handle API routes
	api.Post("/notes", handlers.CreateNoteHandler)
	// api.Get("/notes", handlers.GetAllNotes)
	api.Delete("/notes/:id", handlers.DeleteNoteHandler)

	// Catch-all route for templates
	app.Get("/:route_name", handlers.LoadTemplate)
	// app.Get("/refresh/notes", refreshNotes)
}
