package repo

import "dater/backend/internal/domain"

type BaseRepo interface {
	Create(base *domain.Base) (*domain.Base, error)
	Update(base *domain.Base) (*domain.Base, error)
	Delete(id int64) (int64, error)
	GetByID(id int64) (*domain.Base, error)
	GetAllTables(baseID int64) ([]*domain.Table, error)
}

type TableRepo interface {
	Create(table *domain.Table) (*domain.Table, error)
	GetByID(id int64) (*domain.Table, error)
	Update(table *domain.Table) (*domain.Table, error)
	Delete(id int64) (int64, error)
	GetAllRecords(tableID int64) ([]*domain.Record, error)
	GetAllFields(tableID int64) ([]*domain.Field, error)
}

type RecordRepo interface {
	Create(record *domain.Record) error
	GetByID(id int64) (*domain.Record, error)
	Update(record *domain.Record) error
	Delete(id int64) error
	GetAllRecordsBySourceID(recordID int64) ([]*domain.Record, error)
	GetAllRecordsByTargetID(recordID int64) ([]*domain.Record, error)
	GetAllCells(recordID int64) ([]*domain.Cell, error)
}

type RecordRelationRepo interface {
	Create(recordRelation *domain.RecordRelation) error
	GetByID(id int64) (*domain.RecordRelation, error)
	Update(recordRelation *domain.RecordRelation) error
	Delete(id int64) error
	GetAllRecordsBySourceID(recordID int64) ([]*domain.Record, error)
	GetAllRecordsByTargetID(recordID int64) ([]*domain.Record, error)
	GetAllFields(fieldID int64) ([]*domain.Field, error)
}

type FieldRepo interface {
	Create(field *domain.Field) error
	GetByID(id int64) (*domain.Field, error)
	Update(field *domain.Field) error
	Delete(id int64) (int64, error)
	GetAllRelationsRecords(fieldID int64) ([]*domain.Record, error)
	GetAllCells(fieldID int64) ([]*domain.Cell, error)
	GetAllSelectOptions(fieldID int64) ([]*domain.SelectOption, error)
}

type CellRepo interface {
	Create(cell *domain.Cell) error
	GetByID(id int64) (*domain.Cell, error)
	Update(cell *domain.Cell) error
	Delete(id int64) error
	GetAllCellsByRecordID(recordID int64) ([]*domain.Cell, error)
	GetAllCellsByFieldID(fieldID int64) ([]*domain.Cell, error)
}

type SelectOptionRepo interface {
	Create(selectOption *domain.SelectOption) error
	GetByID(id int64) (*domain.SelectOption, error)
	Update(selectOption *domain.SelectOption) error
	Delete(id int64) error
}
