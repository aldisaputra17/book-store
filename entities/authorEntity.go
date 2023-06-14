package entities

import "github.com/google/uuid"

type Author struct {
	ID      uuid.UUID `gorm:"primaryKey" json:"id"`
	Name    string    `gorm:"type:varchar(255)" json:"name"`
	Country string    `json:"country"`
	BookID  []string  `json:"book_id" gorm:"-"`
	Books   []*Book   `gorm:"many2many:author_books;foreignKey:ID;joinForeignKey:AuthorID;References:ID;joinReferences:BookID" json:"books"`
}
