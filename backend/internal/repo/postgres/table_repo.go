package postgres

import (
	"dater/backend/internal/domain"
	"dater/backend/internal/repo"

	"gorm.io/gorm"
)

type postgresTableRepo struct {
	db *gorm.DB
}

func NewPostgresTableRepo(db *gorm.DB) repo.TableRepo {
	return &postgresTableRepo{
		db: db,
	}
}

func (r *postgresTableRepo) GetAllTablesByBase(baseID int64) ([]*domain.Table, error) {
	var tables []*domain.Table
	err := r.db.Where("base_id = ?", baseID).Find(&tables).Error
	return tables, err
}

func (r *postgresTableRepo) GetAllRecords(tableID int64) ([]*domain.Record, error) {
	var records []*domain.Record
	err := r.db.Where("table_id = ?", tableID).Find(&records).Error
	return records, err
}

func (r *postgresTableRepo) GetAllFields(tableID int64) ([]*domain.Field, error) {
	var fields []*domain.Field
	err := r.db.Where("table_id = ?", tableID).Find(&fields).Error
	return fields, err
}

func (r *postgresTableRepo) GetByID(id int64) (*domain.Table, error) {
	var table domain.Table
	err := r.db.Where("id = ?", id).First(&table).Error
	return &table, err
}

func (r *postgresTableRepo) Create(table *domain.Table) error {
	return r.db.Create(table).Error
}

func (r *postgresTableRepo) Update(table *domain.Table) error {
	return r.db.Save(table).Error
}

func (r *postgresTableRepo) Delete(id int64) error {
	return r.db.Delete(&domain.Table{}, id).Error
}