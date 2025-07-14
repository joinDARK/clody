package models

type RecordRelation struct {
	ID             int64 `json:"id"`
	FieldID        int64 `json:"field_id"`
	SourceRecordID int64 `json:"source_record_id"`
	TargetRecordID int64 `json:"target_record_id"`
}

func NewRecordRelation(id, fieldId, sourceRecordId, targetRecordId int64) *RecordRelation {
	return &RecordRelation{
		ID:             id,
		FieldID:        fieldId,
		SourceRecordID: sourceRecordId,
		TargetRecordID: targetRecordId,
	}
}
