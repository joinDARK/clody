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

func (r *postgresRecordRepo) GetTableRecord(tableID int64) (*domain.Record, error) {
	var record domain.Record
	if err := r.db.Where("table_id = ?", tableID).First(&record).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *postgresRecordRepo) GetByID(id int64) (*domain.Record, error) {
	var record domain.Record
	if err := r.db.Where("id = ?", id).First(&record).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *postgresRecordRepo) Update(record *domain.Record) error {
	return r.db.Save(record).Error
}

func (r *postgresRecordRepo) Delete(id int64) error {
	return r.db.Delete(&domain.Record{}, id).Error
}
