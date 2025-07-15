package models

type SelectOption struct {
	ID       int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	FieldID  int64  `gorm:"not null;index" json:"field_id"`
	Label    string `gorm:"size:255;not null" json:"label"`
	Color    string `gorm:"size:255;not null" json:"color"`
	Position int    `gorm:"not null" json:"position"`
	Field    *Field `gorm:"foreignKey:FieldID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"field,omitempty"`
}
