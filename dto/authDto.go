package dto

type AuthRequest struct {
	Email    string `json:"email" form:"email" binding:"required"`
	Password string `json:"password" form:"password" binding:"required,min=6"`
}
