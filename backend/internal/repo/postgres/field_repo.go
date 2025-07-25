package postgres

import (
	"dater/backend/internal/domain"
	"dater/backend/internal/repo"

	"gorm.io/gorm"
)

type postgresFieldRepo struct {
	db *gorm.DB
}

func NewFieldRepo(db *gorm.DB) repo.FieldRepo {
	return &postgresFieldRepo{db: db}
}

func (r *postgresFieldRepo) GetByTable(tableID int64) ([]*domain.Field, error) {
	var fields []*domain.Field
	if err := r.db.Where("table_id = ?", tableID).Find(&fields).Error; err != nil {
		return nil, err
	}
	return fields, nil
}

func (r *postgresFieldRepo) Create(field *domain.Field) error {
	return r.db.Create(field).Error
}

func (r *postgresFieldRepo) GetByID(id int64) (*domain.Field, error) {
	var field domain.Field
	if err := r.db.First(&field, id).Error; err != nil {
		return nil, err
	}
	return &field, nil
}

func (r *postgresFieldRepo) Update(field *domain.Field) error {
	return r.db.Save(field).Error
}

func (r *postgresFieldRepo) Delete(id int64) error {
	return r.db.Delete(&domain.Field{}, id).Error
}