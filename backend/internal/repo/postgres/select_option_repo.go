package postgres

import (
	"dater/backend/internal/domain"
	"dater/backend/internal/repo"

	"gorm.io/gorm"
)

type postgresSelectOptionRepo struct {
	db *gorm.DB
}

func NewPostgresSelectOptionRepo(db *gorm.DB) repo.SelectOptionRepo {
	return &postgresSelectOptionRepo{db: db}
}

func (r *postgresSelectOptionRepo) GetByID(id int64) (*domain.SelectOption, error) {
	var selectOption domain.SelectOption
	err := r.db.First(&selectOption, id).Error
	if err != nil {
		return nil, err
	}
	return &selectOption, nil
}

func (r *postgresSelectOptionRepo) Create(selectOption *domain.SelectOption) error {
	return r.db.Create(selectOption).Error
}

func (r *postgresSelectOptionRepo) Update(selectOption *domain.SelectOption) error {
	return r.db.Save(selectOption).Error
}

func (r *postgresSelectOptionRepo) Delete(id int64) error {
	return r.db.Delete(&domain.SelectOption{}, id).Error
}
