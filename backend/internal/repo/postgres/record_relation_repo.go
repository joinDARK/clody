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

func (r *postgresRecordRelationRepo) GetAllRecordsBySourceID(recordID int64) ([]*domain.Record, error) {
	var relations []domain.RecordRelation
	err := r.db.Where("source_record_id = ?", recordID).Find(&relations).Error
	if err != nil {
		return nil, err
	}

	var records []*domain.Record
	if len(relations) == 0 {
		return records, nil
	}

	recordIDs := make([]int64, len(relations))
	for i, relation := range relations {
		recordIDs[i] = relation.TargetRecordID
	}

	err = r.db.Where("id IN ?", recordIDs).Find(&records).Error
	return records, err
}

func (r *postgresRecordRelationRepo) GetAllRecordsByTargetID(recordID int64) ([]*domain.Record, error) {
	var relations []domain.RecordRelation
    err := r.db.Where("target_record_id = ?", recordID).Find(&relations).Error
    if err != nil {
        return nil, err
    }

    var records []*domain.Record
    if len(relations) == 0 {
        return records, nil
    }

    recordIDs := make([]int64, len(relations))
    for i, relation := range relations {
        recordIDs[i] = relation.SourceRecordID
    }

    err = r.db.Where("id IN ?", recordIDs).Find(&records).Error
    return records, err
}

func (r *postgresRecordRelationRepo) GetAllFields(fieldID int64) ([]*domain.Field, error) {
	var fields []*domain.Field
    err := r.db.Where("id = ?", fieldID).Find(&fields).Error
    return fields, err
}

func (r *postgresRecordRelationRepo) GetByID(id int64) (*domain.RecordRelation, error) {
	var recordRelation domain.RecordRelation
	err := r.db.First(&recordRelation, id).Error
	return &recordRelation, err
}
