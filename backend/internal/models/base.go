package models

import (
	"time"
)

type Base struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

func NewBase(id int64, name, description string, createdAt time.Time) *Base {
	return &Base{
		Id:          id,
		Name:        name,
		Description: description,
		CreatedAt:   createdAt,
	}
}
