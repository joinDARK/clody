package repo

import "dater/backend/internal/domain"

type BaseRepo interface {
	Create(base *domain.Base) error
	Read() ([]*domain.Base, error)
	Update(base *domain.Base) error
	Delete(id int64) error
	
	GetByID(id int64) (*domain.Base, error)
}

type TableRepo interface {
	Create(table *domain.Table) error
	GetByID(id int64) (*domain.Table, error)
	GetByBase(baseID int64) ([]*domain.Table, error)
	Update(table *domain.Table) error
	Delete(id int64) error
}

type RecordRepo interface {
	Create(record *domain.Record) error
	GetTableRecord(tableID int64) (*domain.Record, error)
	GetByID(id int64) (*domain.Record, error)
	Update(record *domain.Record) error
	Delete(id int64) error
}

type RecordRelationRepo interface {
	Create(recordRelation *domain.RecordRelation) error
	GetByID(id int64) (*domain.RecordRelation, error)
	Update(recordRelation *domain.RecordRelation) error
	Delete(id int64) error
}

type FieldRepo interface {
	Create(field *domain.Field) error
	GetByID(id int64) (*domain.Field, error)
	GetByTable(tableID int64) ([]*domain.Field, error)
	Update(field *domain.Field) error
	Delete(id int64) error
}

type CellRepo interface {
	Create(cell *domain.Cell) error
	GetByID(id int64) (*domain.Cell, error)
	Update(cell *domain.Cell) error
	Delete(id int64) error
}

type SelectOptionRepo interface {
	Create(selectOption *domain.SelectOption) error
	GetByID(id int64) (*domain.SelectOption, error)
	Update(selectOption *domain.SelectOption) error
	Delete(id int64) error
}