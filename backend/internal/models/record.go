package models

import (
	"time"
)

type Record struct {
	Id        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	TableId   int64     `gorm:"not null;index" json:"table_id"`
	Position  int       `gorm:"not null" json:"position"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	Table     *Table    `gorm:"foreignKey:TableID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"table,omitempty"`
	Cells     []*Cell   `gorm:"foreignKey:RecordID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"cells,omitempty"`
}
