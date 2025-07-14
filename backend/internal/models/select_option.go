package models

type SelectOption struct {
	ID       int64  `json:"id"`
	FieldID  int64  `json:"field_id"`
	Label    string `json:"label"`
	Color    string `json:"color"`
	Position int    `json:"position"`
}

func NewSelectOption(id, fieldId int64, label string, color string, position int) *SelectOption {
	return &SelectOption{
		ID:       id,
		FieldID:  fieldId,
		Label:    label,
		Color:    color,
		Position: position,
	}
}
