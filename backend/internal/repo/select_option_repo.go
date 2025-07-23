package repo

import "dater/backend/internal/domain"

type SelectOptionRepo interface {
	Create(field *domain.SelectOption) error
	GetByID(id int64) (*domain.SelectOption, error)
	Update(field *domain.SelectOption) error
	Delete(id int64) error
}
