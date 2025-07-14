package models

import (
	"time"
)

type Cell struct {
	Id        int64      `json:"id"`
	FieldId   int64      `json:"field_id"`
	RecordId  int64      `json:"record_id"`
	Position  int        `json:"position"`
	Value     *CellValue `json:"value"`
	CreatedAt time.Time  `json:"created_at"`
}

type CellValue struct {
	OptionIds   []int64 `json:"option_ids,omitempty"`
	RelationIds []int64 `json:"relation_ids,omitempty"`
	Value       any     `json:"value"`
}

func NewCell(id, fieldId, recordId int64, position int, value *CellValue, createdAt time.Time) *Cell {
	return &Cell{
		Id:        id,
		FieldId:   fieldId,
		RecordId:  recordId,
		Position:  position,
		Value:     value,
		CreatedAt: createdAt,
	}
}

func NewCellValue(optionIds, relationIds []int64, value any) *CellValue {
	return &CellValue{
		OptionIds:   optionIds,
		RelationIds: relationIds,
		Value:       value,
	}
}
