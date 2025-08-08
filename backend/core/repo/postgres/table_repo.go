package postgres

import (
	"clody/core/domain"
	"clody/core/repo"

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

func (r *postgresTableRepo) Create(table *domain.Table) (*domain.Table, error) {
	err := r.db.Create(table).Error
	return table, err
}

func (r *postgresTableRepo) Update(table *domain.Table) (*domain.Table, error) {
	err := r.db.Save(table).Error
	return table, err
}

func (r *postgresTableRepo) Delete(id int64) (int64, error) {
	err := r.db.Delete(&domain.Table{}, id).Error
	return id, err
}