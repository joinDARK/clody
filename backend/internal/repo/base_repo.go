package repo

import "dater/backend/internal/domain"

type BaseRepo interface {
	Create(base *domain.Base) error
	GetByID(id int64) (*domain.Base, error)
	Update(base *domain.Base) error
	Delete(id int64) error
}
