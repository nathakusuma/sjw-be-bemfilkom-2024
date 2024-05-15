package entity

import (
	"database/sql"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Hope struct {
	ID         uuid.UUID `gorm:"orimaryKey; not null; type:uuid; unique"`
	Content    string    `gorm:"not null; type:varchar(2000)"`
	IsApproved sql.NullBool
	CreatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}
