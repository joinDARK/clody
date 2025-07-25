package postgres

import (
	"dater/backend/internal/domain"
	"dater/backend/internal/repo"

	"gorm.io/gorm"
)

type postgresCellRepo struct {
	db *gorm.DB
}

func NewPostgresCellRepo(db *gorm.DB) repo.CellRepo {
	return &postgresCellRepo{
		db: db,
	}
}

func (r *postgresCellRepo) GetByID(id int64) (*domain.Cell, error) {
	var cell domain.Cell
	if err := r.db.Where("id = ?", id).First(&cell).Error; err != nil {
		return nil, err
	}
	return &cell, nil
}

func (r *postgresCellRepo) Create(cell *domain.Cell) error {
	return r.db.Create(cell).Error
}

func (r *postgresCellRepo) Update(cell *domain.Cell) error {
	return r.db.Save(cell).Error
}

func (r *postgresCellRepo) Delete(id int64) error {
	return r.db.Delete(&domain.Cell{}, id).Error
}