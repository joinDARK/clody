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
	return &recordRelation, err
}

func (r *postgresRecordRelationRepo) GetByFieldID(fieldID int64) ([]*domain.RecordRelation, error) {
	var recordRelations []*domain.RecordRelation
	err := r.db.Where("field_id = ?", fieldID).Find(&recordRelations).Error
	return recordRelations, err
}

func (r *postgresRecordRelationRepo) GetBySourceRecordID(recordID int64) ([]*domain.RecordRelation, error) {
	var recordRelations []*domain.RecordRelation
	err := r.db.Where("source_record_id = ?", recordID).Find(&recordRelations).Error
	return recordRelations, err
}

func (r *postgresRecordRelationRepo) GetByTargetRecordID(recordID int64) ([]*domain.RecordRelation, error) {
	var recordRelations []*domain.RecordRelation
	err := r.db.Where("target_record_id = ?", recordID).Find(&recordRelations).Error
	return recordRelations, err
}
