package model

import (
	"github.com/google/uuid"
)

type CreateHopeRequest struct {
	Content string `json:"content" binding:"required,max=2000,min=1"`
}

type CreateHopeResponse struct {
	ID uuid.UUID `json:"id"`
}

type FindHopesLazyLoadRequest struct {
	AfterCreatedAt string `form:"after_created_at" binding:"required,rfc3339"`
	AfterId        string `form:"after_id" binding:"required,uuid"`
}

type FindHopeResponse struct {
	ID      uuid.UUID `json:"id"`
	Content string    `json:"content"`
}

type FindHopeAsAdminResponse struct {
	FindHopeResponse
	IsApproved *bool `json:"is_approved"`
}

type UpdateHopeRequest struct {
	Content    string `json:"content" binding:"omitempty,max=2000,min=1"`
	IsApproved *bool  `json:"is_approved"`
}
