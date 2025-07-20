package table

import (
	"time"

	"dater/backend/internal/field"
	"dater/backend/internal/record"
)

type Table struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	BaseID      int64     `gorm:"not null;index" json:"base_id"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	Description *string   `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time `gorm:"not null" json:"created_at"`

	Fields  []*field.Field   `gorm:"foreignKey:TableID" json:"fields,omitempty"`
	Records []*record.Record  `gorm:"foreignKey:TableID" json:"records,omitempty"`
}
