package entities

import (
	"time"

	"github.com/google/uuid"
)

type Book struct {
	ID            uuid.UUID `gorm:"primaryKey" json:"id"`
	Title         string    `gorm:"type:varchar(255)" json:"title"`
	PublishedYear time.Time `json:"published_year"`
	Isbn          string    `json:"isbn"`
	AuthorID      []string  `gorm:"-" json:"author_id"`
	Authors       []*Author `gorm:"many2many:author_books;foreignKey:ID;joinForeignKey:BookID;References:ID;joinReferences:AuthorID" json:"authors"`
}
