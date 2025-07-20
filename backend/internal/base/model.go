package base

import (
	"dater/backend/internal/table"
	"time"
)

type Base struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	Description *string   `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time `gorm:"not null" json:"created_at"`

	Tables []*table.Table `gorm:"foreignKey:BaseID" json:"tables,omitempty"`
}
