package postgres

import (
	"clody/core/domain"
	"clody/core/repo"

	"gorm.io/gorm"
)

type postgresFieldRepo struct {
	db *gorm.DB
}

func NewPostgresFieldRepo(db *gorm.DB) repo.FieldRepo {
	return &postgresFieldRepo{db: db}
}

func (r *postgresFieldRepo) GetAllRelationsRecords(fieldID int64) ([]*domain.Record, error) {
	// Предполагаем, что есть связь через RecordRelation
	var relations []domain.RecordRelation
	err := r.db.Where("field_id = ?", fieldID).Find(&relations).Error
	if err != nil {
		return nil, err
	}

	var records []*domain.Record
	if len(relations) == 0 {
		return records, nil
	}

	recordIDs := make([]int64, len(relations))
	for i, relation := range relations {
		recordIDs[i] = relation.TargetRecordID
	}

	err = r.db.Where("id IN ?", recordIDs).Find(&records).Error
	return records, err
}

func (r *postgresFieldRepo) GetAllCells(fieldID int64) ([]*domain.Cell, error) {
	var cells []*domain.Cell
	err := r.db.Where("field_id = ?", fieldID).Find(&cells).Error
	return cells, err
}

func (r *postgresFieldRepo) GetAllSelectOptions(fieldID int64) ([]*domain.SelectOption, error) {
	var options []*domain.SelectOption
	err := r.db.Where("field_id = ?", fieldID).Find(&options).Error
	return options, err
}

func (r *postgresFieldRepo) Create(field *domain.Field) error {
	return r.db.Create(field).Error
}

func (r *postgresFieldRepo) GetByID(id int64) (*domain.Field, error) {
	var field domain.Field
	err := r.db.First(&field, id).Error
	return &field, err
}

func (r *postgresFieldRepo) Update(field *domain.Field) error {
	return r.db.Save(field).Error
}

func (r *postgresFieldRepo) Delete(id int64) (int64, error) {
	err := r.db.Delete(&domain.Field{}, id).Error
	return id, err
}
