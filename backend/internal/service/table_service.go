package service

import (
	"dater/backend/internal/domain"
	"dater/backend/internal/repo"

	"github.com/rs/zerolog"
)

type TableService struct {
	repo repo.TableRepo
}

func NewTableService(repo repo.TableRepo) *TableService {
	return &TableService{
		repo: repo,
	}
}

func (s *TableService) GetTableInfo(id int64, logger zerolog.Logger) (*domain.Table, error) {
	logger.Debug().Msg("[GetTableInfo] Getting table info...")
	
	table, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	
	logger.Debug().
		Int64("ID", table.ID).
		Str("Name", table.Name).
		Str("Description", *table.Description).
		Msg("[GetTableInfo] Table info:")
	return table, nil
}
