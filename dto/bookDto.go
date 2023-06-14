package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateBookRequest struct {
	Title    string   `json:"title" form:"title" binding:"required"`
	AuthorID []string `json:"author_id" form:"author_id" binding:"required"`
}

type ReadBookResponse struct {
	ID            string            `json:"id"`
	Title         string            `json:"title"`
	PublishedYear time.Time         `json:"published_year"`
	Isbn          string            `json:"isbn"`
	Author        []*AuthorResponse `json:"author"`
}

type CreateBookResponse struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	PublishedYear time.Time `json:"published_year"`
	Isbn          string    `json:"isbn"`
}

type UpdateBookRequest struct {
	ID    uuid.UUID `json:"id" form:"id" binding:"required"`
	Title string    `json:"title" form:"title" binding:"required"`
}

type UpdateBookResponse struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}
