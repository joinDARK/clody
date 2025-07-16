package models

import (
	"time"
)

type Table struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	BaseID      int64     `gorm:"not null;index" json:"base_id"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	Description *string   `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time `gorm:"not null" json:"created_at"`

	Base    *Base     `gorm:"foreignKey:BaseID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION" json:"base"`
	Fields  []*Field  `gorm:"foreignKey:TableID" json:"fields,omitempty"`
	Records []*Record `gorm:"foreignKey:TableID" json:"records,omitempty"`
}
