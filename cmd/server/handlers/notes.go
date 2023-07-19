// cmd/server/handlers/notes.go
package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pryority/gowind/cmd/server/models"
)

// LoadNotes renders the "notes" template with the list of notes.
func LoadNotes(c *fiber.Ctx) error {
	notes, err := models.GetAllNotes(db)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get notes from the database")
	}
	return c.Render("notes", fiber.Map{"RouteName": "notes", "Notes": notes}, "layouts/main")
}

// CreateNoteHandler handles the creation of a new note.
// func CreateNoteHandler(c *fiber.Ctx) error {
// 	var noteData models.CreateNoteSchema

// 	if err := c.BodyParser(&noteData); err != nil {
// 		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
// 	}

// 	newNote, err := handlers.CreateNote(initializers.DB, noteData)
// 	if err != nil {
// 		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create a new note")
// 	}

// 	return c.Render("note", fiber.Map{"Note": newNote})
// }

func GetAllNotes(c *fiber.Ctx) error {
	var notes []models.Note
	if err := db.Find(&notes).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get notes from the database")
	}
	return c.JSON(notes)

}
