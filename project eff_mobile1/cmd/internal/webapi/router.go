package webapi

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"effective_mobile_service/internal/repository/postgres"
	"effective_mobile_service/internal/usecase"
	"effective_mobile_service/internal/webapi/handler"
	"effective_mobile_service/internal/webapi/middleware"
	"effective_mobile_service/pkg/logger"

	_ "effective_mobile_service/docs"
)

type Server struct {
	router *gin.Engine
	log    logger.Logger
	db     *sql.DB
}

func NewServer(log logger.Logger, db *sql.DB) *Server {
	router := gin.New()

	server := &Server{
		router: router,
		log:    log,
		db:     db,
	}

	server.setupRoutes()

	return server
}

func (s *Server) setupRoutes() {
	// Middleware
	s.router.Use(middleware.LoggerMiddleware(s.log))
	s.router.Use(middleware.RecoveryMiddleware(s.log))
	s.router.Use(gin.Recovery())

	// Swagger
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check
	s.router.GET("/health", func(c *gin.Context) {
		if err := s.db.Ping(); err != nil {
			c.JSON(500, gin.H{"status": "unhealthy"})
			return
		}
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// API routes
	api := s.router.Group("/api/v1")
	{
		s.setupSubscriptionRoutes(api)
	}
}

func (s *Server) setupSubscriptionRoutes(rg *gin.RouterGroup) {
	// Initialize dependencies
	repo := postgres.NewSubscriptionRepository(s.db)
	useCase := usecase.NewSubscriptionUseCase(repo)
	handler := handler.NewSubscriptionHandler(useCase, s.log)

	subscriptions := rg.Group("/subscriptions")
	{
		subscriptions.POST("", handler.CreateSubscription)
		subscriptions.GET("", handler.ListSubscriptions)
		subscriptions.GET("/total-cost", handler.GetTotalCost)
		subscriptions.GET("/:id", handler.GetSubscription)
		subscriptions.PUT("/:id", handler.UpdateSubscription)
		subscriptions.DELETE("/:id", handler.DeleteSubscription)
	}
}

func (s *Server) Run(addr string) error {
	s.log.Info("Starting server", "address", addr)
	return s.router.Run(addr)
}
