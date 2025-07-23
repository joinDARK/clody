package repo

import "dater/backend/internal/domain"

type TableRepo interface {
	Create(base *domain.Table) error
	GetByID(id int64) (*domain.Table, error)
	Update(base *domain.Table) error
	Delete(id int64) error
}
