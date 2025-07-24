package repo

import "dater/backend/internal/domain"

type BaseRepo interface {
	Create(base *domain.Base) error
	GetByID(id int64) (*domain.Base, error)
	Update(base *domain.Base) error
	Delete(id int64) error
}

type TableRepo interface {
	Create(base *domain.Table) error
	GetByID(id int64) (*domain.Table, error)
	GetByBase(baseID int64) ([]*domain.Table, error)
	Update(base *domain.Table) error
	Delete(id int64) error
}

type RecordRepo interface {
	Create(record *domain.Record) error
	GetByID(id int64) (*domain.Record, error)
	Update(record *domain.Record) error
	Delete(id int64) error
}

type RecordRelationRepo interface {
	Create(field *domain.RecordRelation) error
	GetByID(id int64) (*domain.RecordRelation, error)
	Update(field *domain.RecordRelation) error
	Delete(id int64) error
}

type FieldRepo interface {
	Create(field *domain.Field) error
	GetByID(id int64) (*domain.Field, error)
	Update(field *domain.Field) error
	Delete(id int64) error
}

type CellRepo interface {
	Create(base *domain.Cell) error
	GetByID(id int64) (*domain.Cell, error)
	Update(base *domain.Cell) error
	Delete(id int64) error
}

type SelectOptionRepo interface {
	Create(field *domain.SelectOption) error
	GetByID(id int64) (*domain.SelectOption, error)
	Update(field *domain.SelectOption) error
	Delete(id int64) error
}