package field

import (
	"time"

	"gorm.io/datatypes"
)

type Field struct {
	ID          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	TableID     int64          `gorm:"not null;index" json:"table_id"`
	Name        string         `gorm:"size:255;not null" json:"name"`
	Type        FieldType      `gorm:"type:field_type;default:'text';not null" json:"type"`
	Position    int            `gorm:"not null" json:"position"`
	Description *string        `gorm:"type:text" json:"description,omitempty"`
	Options     datatypes.JSON `gorm:"type:jsonb" json:"options,omitempty"`
	CreatedAt   time.Time      `gorm:"not null" json:"created_at"`

	SelectOptions []*SelectOption `gorm:"foreignKey:FieldID" json:"select_options,omitempty"`
}

type SelectOption struct {
	ID       int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	FieldID  int64  `gorm:"not null;index" json:"field_id"`
	Label    string `gorm:"size:255;not null" json:"label"`
	Color    string `gorm:"size:255;not null" json:"color"`
	Position int    `gorm:"not null" json:"position"`
	Field    *Field `gorm:"foreignKey:FieldID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"field,omitempty"`
}

type FieldType string

const (
	FieldTypeText         FieldType = "text"
	FieldTypeNumber       FieldType = "number"
	FieldTypeCheckbox     FieldType = "checkbox"
	FieldTypeDate         FieldType = "date"
	FieldTypeSingleSelect FieldType = "single_select"
	FieldTypeMultiSelect  FieldType = "multi_select"
	FieldTypeMany2One     FieldType = "many2one"
	FieldTypeOne2Many     FieldType = "one2many"
	FieldTypeMany2Many    FieldType = "many2many"
	FieldTypeFormula      FieldType = "formula"
)

func (ft FieldType) IsValid() bool {
	switch ft {
	case
		FieldTypeText,
		FieldTypeNumber,
		FieldTypeCheckbox,
		FieldTypeDate,
		FieldTypeSingleSelect,
		FieldTypeMultiSelect,
		FieldTypeMany2One,
		FieldTypeOne2Many,
		FieldTypeMany2Many,
		FieldTypeFormula:
		return true

	default:
		return false
	}
}
