package postgres

import (
	"dater/backend/internal/domain"
	"dater/backend/internal/repo"

	"gorm.io/gorm"
)

type postgresBaseRepo struct {
	db *gorm.DB
}

func NewPostgresBaseRepo(db *gorm.DB) repo.BaseRepo {
	return &postgresBaseRepo{
		db: db,
	}
}

func (r *postgresBaseRepo) GetByID(id int64) (*domain.Base, error) {
	var base domain.Base
	err := r.db.First(&base, id).Error
	return &base, err
}

func (r *postgresBaseRepo) Create(base *domain.Base) error {
	return r.db.Create(base).Error
}

func (r *postgresBaseRepo) GetAllTables(baseID int64) ([]*domain.Table, error) {
	var tables []*domain.Table
	err := r.db.Where("base_id = ?", baseID).Find(&tables).Error;
	return tables, err
}

func (r *postgresBaseRepo) Update(base *domain.Base) error {
	return r.db.Save(base).Error
}

func (r *postgresBaseRepo) Delete(id int64) error {
	return r.db.Delete(&domain.Base{}, id).Error
}
