package model

import (
	"github.com/google/uuid"
)

type HopeWhisperType string

const (
	HopeCorner  HopeWhisperType = "hopes"
	WhisperWall HopeWhisperType = "whispers"
)

func (hwt HopeWhisperType) String() string {
	return string(hwt)
}

func (hwt HopeWhisperType) Singular() string {
	if hwt == HopeCorner {
		return "hope"
	}
	return "whisper"
}

type CreateHopeWhisperRequest struct {
	Content  string `json:"content" binding:"required,max=2000,min=1"`
	IsPublic bool   `json:"is_public"`
}

type CreateHopeWhisperResponse struct {
	ID uuid.UUID `json:"id"`
}

type FindHopesWhispersLazyLoadRequest struct {
	AfterCreatedAt string `form:"after_created_at" binding:"required,rfc3339"`
	AfterId        string `form:"after_id" binding:"required,uuid"`
}

type FindHopeWhisperResponse struct {
	ID        uuid.UUID `json:"id"`
	Content   string    `json:"content"`
	CreatedAt string    `json:"created_at"`
}

type FindHopeWhisperAsAdminResponse struct {
	FindHopeWhisperResponse
	IsApproved *bool  `json:"is_approved"`
	IsPublic   bool   `json:"is_public"`
	UpdatedAt  string `json:"updated_at"`
}

type UpdateHopeWhisperRequest struct {
	Content    string `json:"content" binding:"omitempty,max=2000,min=1"`
	IsApproved *bool  `json:"is_approved"`
}
