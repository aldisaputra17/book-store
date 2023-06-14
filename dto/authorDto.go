package dto

import "github.com/google/uuid"

type CreateAuthorRequest struct {
	Name    string `json:"name" form:"name" binding:"required"`
	Country string `json:"country" form:"country" binding:"required"`
}

type AuthorResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
}

type ReadAuthorResponse struct {
	ID      string                `json:"id"`
	Name    string                `json:"name"`
	Country string                `json:"country"`
	Book    []*CreateBookResponse `json:"book"`
}

type UpdateAuthorRequest struct {
	ID      uuid.UUID `json:"id" form:"id" binding:"required"`
	Name    string    `json:"name" form:"name"`
	Country string    `json:"country" form:"country"`
}

type UpdateAuthorResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
}
