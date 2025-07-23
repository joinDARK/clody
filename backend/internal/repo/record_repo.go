package repo

import "dater/backend/internal/domain"

type RecordRepo interface {
	Create(record *domain.Record) error
	GetByID(id int64) (*domain.Record, error)
	Update(record *domain.Record) error
	Delete(id int64) error
}
