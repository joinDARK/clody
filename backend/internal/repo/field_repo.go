package repo

import "dater/backend/internal/domain"

type FieldRepo interface {
	Create(field *domain.Field) error
	GetByID(id int64) (*domain.Field, error)
	Update(field *domain.Field) error
	Delete(id int64) error
}
