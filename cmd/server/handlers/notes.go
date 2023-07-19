// cmd/server/handlers/notes.go
package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/pryority/gowind/cmd/server/models"
	"gorm.io/gorm"
)

func LoadNotes(c *fiber.Ctx) error {
	notes, err := GetAllNotes(db)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get notes from the database")
	}
	return c.Render("notes", fiber.Map{"RouteName": "notes", "Notes": notes}, "layouts/main")
}

func CreateNoteHandler(c *fiber.Ctx) error {
	var noteData models.CreateNoteSchema

	// Parse the request body and validate the data
	if err := c.BodyParser(&noteData); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Create the new note in the database

	if newNote, err := models.CreateNote(db, noteData); err != nil {
		log.Println(newNote)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create a new note")
	}

	// Load all notes again to get the updated list of notes
	notes, err := GetAllNotes(db)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get notes from the database")
	}

	// Render the "notes" template with the updated list of notes
	return c.Render("partials/note.html", fiber.Map{"Notes": notes})
}

func DeleteNoteHandler(c *fiber.Ctx) error {
	noteID := c.Params("id")
	id, err := uuid.Parse(noteID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid note ID")
	}

	if err := models.DeleteNote(db, id); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete the note")
	}

	return c.SendString("Note deleted successfully")
}

func GetAllNotes(db *gorm.DB) ([]models.Note, error) {
	var notes []models.Note
	if err := db.Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}
