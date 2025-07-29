package service

import (
	"dater/backend/internal/domain"
	"dater/backend/internal/repo"

	"github.com/rs/zerolog"
)

type BaseService struct {
	repo repo.BaseRepo
}

func NewBaseService(repo repo.BaseRepo) *BaseService {
	return &BaseService{
		repo: repo,
	}
}

func (s *BaseService) GetBaseInfo(id int64, logger zerolog.Logger) (*domain.Base, error) {
	logger.Debug().Msg("[GetBaseInfo] Getting base info...")
	
	base, err := s.repo.GetByID(id)
	if err != nil {
		logger.Error().Err(err).Msg("[GetBaseInfo] Failed to get base info")
		return nil, err
	}
	
	logger.Debug().
		Int64("ID", base.ID).
		Str("Name", base.Name).
		Str("Description", *base.Description).
		Msg("[GetBaseInfo] Base info:")
	return base, nil
}