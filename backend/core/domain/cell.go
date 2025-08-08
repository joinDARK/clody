package domain

import (
	"time"

	"gorm.io/datatypes"
)

type Cell struct {
	ID        int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	FieldID   int64          `gorm:"not null;index" json:"field_id"`
	RecordID  int64          `gorm:"not null;index" json:"record_id"`
	Position  int            `gorm:"not null" json:"position"`
	Value     datatypes.JSON `gorm:"type:jsonb" json:"value"`
	CreatedAt time.Time      `gorm:"not null" json:"created_at"`

	Field  *Field  `gorm:"foreignKey:FieldID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"field,omitempty"`
	Record *Record `gorm:"foreignKey:RecordID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"record,omitempty"`
}
