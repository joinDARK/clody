package postgres

import (
	"clody/core/domain"
	"clody/core/repo"

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

func (r *postgresSelectOptionRepo) GetAllByFieldID(fieldID int64) ([]*domain.SelectOption, error) {
	var selectOptions []*domain.SelectOption
	err := r.db.Where("field_id = ?", fieldID).Find(&selectOptions).Error
	if err != nil {
		return nil, err
	}
	return selectOptions, nil
}
