package postgres

import (
	"dater/backend/internal/domain"
	"dater/backend/internal/repo"

	"gorm.io/gorm"
)

type postgresRecordRepo struct {
	db *gorm.DB
}

func NewRecordRepo(db *gorm.DB) repo.RecordRepo {
	return &postgresRecordRepo{db: db}
}

func (r *postgresRecordRepo) Create(record *domain.Record) error {
	return r.db.Create(record).Error
}

func (r *postgresRecordRepo) GetAllRecordsByTable(tableID int64) ([]*domain.Record, error) {
	var records []*domain.Record
	err := r.db.Where("table_id = ?", tableID).Find(&records).Error
	return records, err
}

func (r *postgresRecordRepo) GetAllRecordsBySourceID(recordID int64) ([]*domain.Record, error) {
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

func (r *postgresRecordRepo) GetAllRecordsByTargetID(recordID int64) ([]*domain.Record, error) {
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

func (r *postgresRecordRepo) GetAllCells(recordID int64) ([]*domain.Cell, error) {
	var cells []*domain.Cell
    err := r.db.Where("record_id = ?", recordID).Find(&cells).Error
    return cells, err
}

func (r *postgresRecordRepo) GetByID(id int64) (*domain.Record, error) {
	var record domain.Record
	err := r.db.Where("id = ?", id).First(&record).Error
	return &record, err
}

func (r *postgresRecordRepo) Update(record *domain.Record) error {
	return r.db.Save(record).Error
}

func (r *postgresRecordRepo) Delete(id int64) error {
	return r.db.Delete(&domain.Record{}, id).Error
}
