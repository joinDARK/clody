package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"dater/backend/internal/config"
	"dater/backend/internal/handlers"
	"dater/backend/internal/router"
)

type Server struct {
	Router   *gin.Engine
	Database *gorm.DB
	Logger   zerolog.Logger
}

func NewServer(database *gorm.DB, logger zerolog.Logger) *Server {
	return &Server{
		Database: database,
		Logger:   logger,
	}
}

func (s *Server) InitRouter(cfg *config.Server) {
	s.Logger.Debug().Msg("Initializing router...")
	r := router.NewRouter(cfg)
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

	r.GET("/ping", handlers.Ping)
	{
		base := r.Group("/bases")
		base.POST("/", handlers.CreateBase(s.Database))
		base.GET("/", handlers.GetBase(s.Database))
		base.PATCH("/:id", handlers.UpdateBase(s.Database))
		base.DELETE("/:id", handlers.DeleteBase(s.Database))
	}
	{
		table := r.Group("/tables")
		table.POST("/", handlers.CreateTable(s.Database))
		table.GET("/", handlers.GetTable(s.Database))
		table.PATCH("/:id", handlers.UpdateTable(s.Database))
		table.DELETE("/:id", handlers.DeleteTable(s.Database))
	}

	s.Logger.Debug().Int("routes_count", len(r.Routes())).Msg("Router handlers are initialized")
	s.Logger.Info().Msg("Router handlers are initialized")
}
