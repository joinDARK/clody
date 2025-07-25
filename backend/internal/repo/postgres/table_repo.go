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

func (r *postgresTableRepo) GetByBase(baseID int64) ([]*domain.Table, error) {
	var tables []*domain.Table
	if err := r.db.Where("base_id = ?", baseID).Find(&tables).Error; err != nil {
		return nil, err
	}
	return tables, nil
}

func (r *postgresTableRepo) GetByID(id int64) (*domain.Table, error) {
	var table domain.Table
	if err := r.db.Where("id = ?", id).First(&table).Error; err != nil {
		return nil, err
	}
	return &table, nil
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