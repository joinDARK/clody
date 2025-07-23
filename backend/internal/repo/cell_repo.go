package repo

import "dater/backend/internal/domain"

type CellRepo interface {
	Create(base *domain.Cell) error
	GetByID(id int64) (*domain.Cell, error)
	Update(base *domain.Cell) error
	Delete(id int64) error
}
