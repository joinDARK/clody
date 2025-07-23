package repo

import "dater/backend/internal/domain"

type RecordRelationRepo interface {
	Create(field *domain.RecordRelation) error
	GetByID(id int64) (*domain.RecordRelation, error)
	Update(field *domain.RecordRelation) error
	Delete(id int64) error
}
