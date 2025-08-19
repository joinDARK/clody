package service

import (
	"clody/core/domain"
	"clody/core/repo"
	"errors"

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

	if id <= 0 {
		logger.Error().Msg("[GetBaseInfo] ID cannot be less than or equal to 0")
		return nil, errors.New("ID cannot be less than or equal to 0")
	}

	base, err := s.repo.GetByID(id)
	if err != nil {
		logger.Error().Err(err).Msg("[GetBaseInfo] Failed to get base info")
		return nil, err
	}

	desc := ""
	if base.Description != nil {
		desc = *base.Description
	}

	logger.Debug().
		Int64("ID", base.ID).
		Str("Name", base.Name).
		Str("Description", desc).
		Msg("[GetBaseInfo] Base info:")
	return base, nil
}

func (s *BaseService) CreateBase(name string, desc *string, logger zerolog.Logger) (*domain.Base, error) {
	logger.Debug().Msg("[CreateBase] Creating base...")

	if name == "" {
		logger.Error().Msg("[CreateBase] Name cannot be empty")
		return nil, errors.New("Name cannot be empty")
	}

	base, err := s.repo.Create(&domain.Base{Name: name, Description: desc})
	if err != nil {
		logger.Error().Err(err).Msg("[CreateBase] Failed to create base")
		return nil, err
	}

	descLog := ""
	if base.Description != nil {
		descLog = *base.Description
	}

	logger.Debug().
		Int64("ID", base.ID).
		Str("Name", base.Name).
		Str("Description", descLog).
		Msg("[CreateBase] Base created:")
	return base, nil
}

func (s *BaseService) UpdateBase(id int64, name string, desc *string, logger zerolog.Logger) (*domain.Base, error) {
	logger.Debug().Msg("[UpdateBase] Updating base...")

	if id <= 0 {
		logger.Error().Msg("[UpdateBase] ID cannot be less than or equal to 0")
		return nil, errors.New("ID cannot be less than or equal to 0")
	}

	if name == "" {
		logger.Error().Msg("[UpdateBase] Name cannot be empty")
		return nil, errors.New("Name cannot be empty")
	}

	// FIXME: Улучшить обновление данных
	base, err := s.GetBaseInfo(id, logger)
	if err != nil {
		logger.Error().Err(err).Msg("[UpdateBase] Base not found")
		return nil, err
	}

	base.Name = name
	if desc != nil {
		base.Description = desc
	}

	base, err = s.repo.Update(base)
	if err != nil {
		logger.Error().Err(err).Msg("[UpdateBase] Failed to update base")
		return nil, err
	}

	descLog := ""
	if base.Description != nil {
		descLog = *base.Description
	}

	logger.Debug().
		Int64("ID", base.ID).
		Str("Name", base.Name).
		Str("Description", descLog).
		Msg("[UpdateBase] Base updated:")
	return base, nil
}

func (s *BaseService) DeleteBase(id int64, logger zerolog.Logger) (int64, error) {
	logger.Debug().Msg("[DeleteBase] Deleting base...")

	if id == 0 {
		logger.Error().Msg("[DeleteBase] ID cannot be 0")
		return 0, errors.New("ID cannot be 0")
	}

	deleteID, err := s.repo.Delete(id)
	if err != nil {
		logger.Error().Err(err).Msg("[DeleteBase] Failed to delete base")
		return 0, err
	}

	logger.Debug().
		Int64("ID", deleteID).
		Msg("[DeleteBase] Base deleted:")
	return deleteID, nil
}

func (s *BaseService) GetBaseTables(id int64, logger zerolog.Logger) ([]*domain.Table, error) {
	logger.Debug().Msg("[GetBaseTables] Getting base tables...")

	if id <= 0 {
		logger.Error().Msg("[UpdateBase] ID cannot be less than or equal to 0")
		return nil, errors.New("ID cannot be less than or equal to 0")
	}

	baseTables, err := s.repo.GetAllTables(id)
	if err != nil {
		logger.Error().Err(err).Msg("[GetBaseTables] Failed to get base tables")
		return nil, err
	}

	logger.Debug().
		Int64("ID", id).
		Int("Count", len(baseTables)).
		Msg("[GetBaseTables] Base tables retrieved:")
	return baseTables, nil
}
