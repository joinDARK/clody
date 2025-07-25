package postgres

import (
	"dater/backend/internal/domain"
	"dater/backend/internal/repo"

	"gorm.io/gorm"
)

type postgresRecordRelationRepo struct {
	db *gorm.DB
}

func NewRecordRelationRepo(db *gorm.DB) repo.RecordRelationRepo {
	return &postgresRecordRelationRepo{db: db}
}

func (r *postgresRecordRelationRepo) Create(recordRelation *domain.RecordRelation) error {
	return r.db.Create(recordRelation).Error
}

func (r *postgresRecordRelationRepo) Update(recordRelation *domain.RecordRelation) error {
	return r.db.Save(recordRelation).Error
}

func (r *postgresRecordRelationRepo) Delete(id int64) error {
	return r.db.Delete(&domain.RecordRelation{}, id).Error
}

func (r *postgresRecordRelationRepo) GetByID(id int64) (*domain.RecordRelation, error) {
	var recordRelation domain.RecordRelation
	err := r.db.First(&recordRelation, id).Error
	if err != nil {
		return nil, err
	}
	return &recordRelation, nil
}