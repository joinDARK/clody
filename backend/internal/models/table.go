package models

import (
	"time"
)

type Table struct {
	Id          int64     `json:"id"`
	BaseId      int64     `json:"base_id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

func NewTable(id, baseId int64, name string, description string) *Table {
	return &Table{
		Id:          id,
		BaseId:      baseId,
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
	}
}
