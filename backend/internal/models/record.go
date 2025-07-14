package models

import (
	"time"
)

type Record struct {
	Id        int64     `json:"id"`
	TableId   int64     `json:"table_id"`
	Position  int       `json:"position"`
	CreatedAt time.Time `json:"created_at"`
}

func NewRecord(id, tableId int64, position int, createdAt time.Time) *Record {
	return &Record{
		Id:        id,
		TableId:   tableId,
		Position:  position,
		CreatedAt: createdAt,
	}
}
