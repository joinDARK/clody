package domain

type RecordRelation struct {
	ID             int64 `gorm:"primaryKey;autoIncrement" json:"id"`
	FieldID        int64 `gorm:"index" json:"field_id"`
	SourceRecordID int64 `gorm:"not null;index" json:"source_record_id"`
	TargetRecordID int64 `gorm:"not null;index" json:"target_record_id"`

	Field        *Field  `gorm:"foreignKey:FieldID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"field,omitempty"`
	SourceRecord *Record `gorm:"foreignKey:SourceRecordID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"source_record,omitempty"`
	TargetRecord *Record `gorm:"foreignKey:TargetRecordID;references:ID;constraint:OnUpdate:SET NULL,OnDelete:CASCADE" json:"target_record,omitempty"`
}