package server

import (
	"fmt"

	"clody/core/config"
	"clody/core/transport/api"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Server struct {
	Router   *gin.Engine
	Config   config.Config
	Database *gorm.DB
	Logger   zerolog.Logger
}

func NewServer(cfg *config.Config) *Server {
	logger := InitLogger(&cfg.Logger)
	
	database, err := NewDB(&cfg.Database, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize database")
	}
	
	router := api.SetupRouter(&cfg.Server, logger, database)
	
	return &Server{
		Database: database,
		Logger:   logger,
		Config:   *cfg,
		Router:   router,
	}
}

func (s *Server) Start() {
	if s.Router == nil {
		s.Logger.Fatal().Msg("Router is not initialized")
	}
	s.Logger.Debug().Msg("Starting server...")

	addr := fmt.Sprintf("%s:%d", s.Config.Server.Host, s.Config.Server.Port)

	s.Logger.Debug().Msg("Server is started: http://" + addr)
	s.Logger.Info().Msg("Server is started: http://" + addr)
	s.Router.Run(addr)
}
