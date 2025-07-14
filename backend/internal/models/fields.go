package models

import (
	"time"
)

type Field struct {
	Id          int64     `json:"id"`
	TableId     int64     `json:"table_id"`
	Name        string    `json:"name"`
	Type        FieldType `json:"type"`
	Position    int       `json:"position"`
	Description string    `json:"description,omitempty"`
	Options     *Option   `json:"options,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type Option struct {
	OptionIds   []int64 `json:"option_ids,omitempty"`
	RelationIds []int64 `json:"relation_ids,omitempty"`
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

func NewField(id, tableId int64, name string, typeStr FieldType, position int, description string, options *Option, createdAt time.Time) *Field {
	return &Field{
		Id:          id,
		TableId:     tableId,
		Name:        name,
		Type:        typeStr,
		Position:    position,
		Description: description,
		Options:     options,
		CreatedAt:   createdAt,
	}
}

func NewOption(optionIds, relationIds []int64) *Option {
	return &Option{
		OptionIds:   optionIds,
		RelationIds: relationIds,
	}
}

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
