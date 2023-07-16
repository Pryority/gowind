package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/pryority/gowind/cmd/server/models"
	"github.com/wpcodevo/golang-fiber/initializers"
)

const (
	templateFolder = "./cmd/web/views"
	staticFolder   = "./cmd/web/public"
	localPort      = ":8080"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatalln("Failed to load environment variables:", err.Error())
	}
	initializers.ConnectDB(&config)
}

func main() {
	engine := html.New(templateFolder, ".html")

	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "layouts/main",
	})

	// Serve static files
	app.Static("/public", staticFolder)

	// Group related routes with subrouter
	api := app.Group("/api")

	// Define routes
	app.Get("/", loadMain)
	app.Get("/notes", loadNotes)

	// Handle API routes
	api.Post("/notes", createNoteHandler)
	api.Get("/notes", getAllNotesHandler)

	// Catch-all route for templates
	app.Get("/:route_name", loadTemplate)

	log.Printf("🌩 Server started at http://localhost%s", localPort)

	// Start the server
	if err := app.Listen(localPort); err != nil {
		log.Fatalf("Error starting the server: %v", err)
	}
}

func loadMain(c *fiber.Ctx) error {
	return c.Render("index", nil)
}

func loadNotes(c *fiber.Ctx) error {
	notes, err := models.GetAllNotes(initializers.DB)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get notes from the database")
	}
	return c.Render("notes", fiber.Map{"RouteName": "notes", "Notes": notes}, "layouts/main")
}

func loadTemplate(c *fiber.Ctx) error {
	routeName := c.Params("route_name")
	return c.Render("layouts/main", fiber.Map{"RouteName": routeName})
}

func createNoteHandler(c *fiber.Ctx) error {
	var noteData models.CreateNoteSchema

	if err := c.BodyParser(&noteData); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	newNote, err := models.CreateNote(initializers.DB, noteData)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create a new note")
	}
	return c.JSON(newNote)
}

func getAllNotesHandler(c *fiber.Ctx) error {
	notes, err := models.GetAllNotes(initializers.DB)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get notes from the database")
	}

	log.Println("Number of notes:", len(notes))
	log.Println(notes)
	return c.JSON(notes)
}
