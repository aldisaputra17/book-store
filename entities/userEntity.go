package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"uniqueIndex;type:varchar(255)" json:"email" `
	Password  string    `gorm:"->;<-;not null" json:"-" validate:"required, min=6"`
	Token     string    `gorm:"-" json:"token,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
