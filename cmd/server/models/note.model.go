package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Note struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	Title     string    `gorm:"varchar(255);uniqueIndex;not null" json:"title,omitempty"`
	Content   string    `gorm:"not null" json:"content,omitempty"`
	Category  string    `gorm:"varchar(100)" json:"category,omitempty"`
	Published bool      `gorm:"default:false;not null" json:"published"`
	CreatedAt time.Time `gorm:"not null" json:"createdAt,omitempty"`
	UpdatedAt time.Time `gorm:"not null" json:"updatedAt,omitempty"`
}

var validate = validator.New()
var db *gorm.DB

type ErrorResponse struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value,omitempty"`
}

func ValidateStruct[T any](payload T) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.Field = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

type CreateNoteSchema struct {
	Title     string `json:"title" validate:"required"`
	Content   string `json:"content" validate:"required"`
	Category  string `json:"category,omitempty"`
	Published bool   `json:"published,omitempty"`
}

type UpdateNoteSchema struct {
	Title     string `json:"title,omitempty"`
	Content   string `json:"content,omitempty"`
	Category  string `json:"category,omitempty"`
	Published *bool  `json:"published,omitempty"`
}

func CreateNote(db *gorm.DB, noteData CreateNoteSchema) (*Note, error) {
	newNote := Note{
		ID:        uuid.New(),
		Title:     noteData.Title,
		Content:   noteData.Content,
		Category:  noteData.Category,
		Published: noteData.Published,
	}
	if err := db.Create(&newNote).Error; err != nil {
		return nil, err
	}
	return &newNote, nil
}

func DeleteNote(db *gorm.DB, id uuid.UUID) error {
	var note Note
	if err := db.First(&note, id).Error; err != nil {
		return err
	}

	if err := db.Delete(&note).Error; err != nil {
		return err
	}

	return nil
}
