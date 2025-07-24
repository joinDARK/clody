package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"dater/backend/internal/config"
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
	
	return &Server{
		Database: database,
		Logger:   logger,
		Config:   *cfg,
	}
}

func (s *Server) InitRouter(cfg *config.Server) {
	s.Logger.Debug().Msg("Initializing router...")

	
	r := NewRouter(cfg, s.Logger)
	s.InitHandlers(r)
	s.Router = r
	
	s.Logger.Debug().Msg("Router is initialized")
	s.Logger.Info().Msg("Router is initialized")
}

func (s *Server) Start(cfg *config.Server) {
	s.Logger.Debug().Msg("Starting server...")

	s.InitRouter(cfg)
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	s.Logger.Debug().Msg("Server is started: http://" + addr)
	s.Logger.Info().Msg("Server is started: http://" + addr)
	s.Router.Run(addr)
}

func (s *Server) InitHandlers(r *gin.Engine) {
	s.Logger.Debug().Msg("Initializing router handlers...")

	r.GET("/ping", s.Ping)
	
	base.InitBaseRoutes(r, s.Database, s.Logger)
	table.InitTableRoutes(r, s.Database, s.Logger)
	field.InitFieldRoutes(r, s.Database, s.Logger)

	s.Logger.Debug().Int("routes_count", len(r.Routes())).Msg("Router handlers are initialized")
	s.Logger.Info().Msg("Router handlers are initialized")
}

func (s *Server) Ping(c *gin.Context) {
	s.Logger.Info().Msg("Ping returned a response, the server is running and connection is successful")
	c.JSON(200, gin.H{"message": "connection is successful"})
}
