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
	var entity domain.Base
	if err := r.db.First(&entity, id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *postgresBaseRepo) Create(base *domain.Base) error {
	return r.db.Create(base).Error
}

func (r *postgresBaseRepo) Read() ([]*domain.Base, error) {
	var bases []*domain.Base
	if err := r.db.Find(&bases).Error; err != nil {
		return nil, err
	}
	return bases, nil
}

func (r *postgresBaseRepo) Update(base *domain.Base) error {
	return r.db.Save(base).Error
}

func (r *postgresBaseRepo) Delete(id int64) error {
	return r.db.Delete(&domain.Base{}, id).Error
}
