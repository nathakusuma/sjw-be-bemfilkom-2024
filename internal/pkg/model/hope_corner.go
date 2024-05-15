package model

import (
	"database/sql"
	"github.com/google/uuid"
)

type CreateHopeRequest struct {
	Content string `json:"content" binding:"required,max=2000,min=1"`
}

type CreateHopeResponse struct {
	ID uuid.UUID `json:"id"`
}

type GetHopesLazyLoadRequest struct {
	AfterCreatedAt string `form:"after_created_at" binding:"required,rfc3339"`
	AfterId        string `form:"after_id" binding:"required,uuid"`
}

type GetHopeResponse struct {
	ID      uuid.UUID `json:"id"`
	Content string    `json:"content"`
}

type GetHopeAsAdminResponse struct {
	GetHopeResponse
	IsApproved sql.NullBool `json:"is_approved"`
}
