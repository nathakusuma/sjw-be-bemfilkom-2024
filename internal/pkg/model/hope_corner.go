package model

import "github.com/google/uuid"

type CreateHopeRequest struct {
	Content string `json:"content" binding:"required,max=2000,min=1"`
}

type CreateHopeResponse struct {
	ID uuid.UUID `json:"id"`
}
