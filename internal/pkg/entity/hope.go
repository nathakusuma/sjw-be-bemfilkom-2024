package entity

import (
	"database/sql"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Hope struct {
	gorm.Model
	ID         uuid.UUID `gorm:"orimaryKey; not null; type:uuid; unique"`
	Content    string    `gorm:"not null; type:varchar(2000)"`
	IsApproved sql.NullBool
}