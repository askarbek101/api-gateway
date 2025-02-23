package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"api-gateway/config"
	"api-gateway/internal/handlers"
	"api-gateway/internal/middleware"
	"api-gateway/internal/models/responses"
	"api-gateway/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Server represents the HTTP server
type Server struct {
	engine      *gin.Engine
	config      *config.Config
	userHandler *handlers.UserHandler
	testHandler *handlers.TestHandler
	httpServer  *http.Server
}

// New creates a new server instance with middleware
func New(cfg *config.Config) (*Server, error) {

	// Initialize services
	userService := services.NewUserService(cfg.ExternalServices.UserService)
	testService := services.NewTestService()

	// Create handlers
	userHandler := handlers.NewUserHandler(userService)
	testHandler := handlers.NewTestHandler(testService)

	// Set Gin mode

	// Initialize Redis client with cache config
	middleware.InitRedis(fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port), &cfg.Cache)

	// Initialize rate limiter with config
	middleware.InitRateLimit(&cfg.RateLimit)

	// Create server instance
	s := &Server{
		engine:      gin.New(),
		config:      cfg,
		userHandler: userHandler,
		testHandler: testHandler,
	}
	s.initRoutes()

	// Initialize databas
	return s, nil
}

// initDatabase initializes the database connection and sets up repositories and services
func (s *Server) initRoutes() error {
	gin.SetMode(s.config.Server.Mode)

	// Configure trusted proxies
	s.engine.SetTrustedProxies([]string{s.config.Server.TrustedProxy})

	// Initialize JWT secret
	middleware.SetJWTSecret(s.config.JWT.Secret)

	// Configure CORS
	s.engine.Use(cors.New(cors.Config{
		AllowOrigins:     s.config.Server.AllowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "Cache-Control", "If-None-Match"},
		ExposeHeaders:    []string{"Content-Length", "ETag"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Add middlewares
	s.engine.Use(middleware.Recovery()) // Custom recovery middleware
	s.engine.Use(middleware.Logger())
	s.engine.Use(middleware.RateLimit()) // Add rate limiting middleware
	s.engine.Use(middleware.Cache())     // Apply Redis cache middleware globally
	s.registerHttpRoutes()

	return nil
}

// registerRoutes sets up all the routes for the server
func (s *Server) registerHttpRoutes() {
	// Swagger documentation
	s.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("http://localhost:"+s.config.Server.Port+"/swagger/doc.json")))

	// Public routes
	// @Summary Health check
	// @Description Check if the API is up and running
	// @Tags health
	// @Accept json
	// @Produce json
	// @Success 200 {object} responses.HealthResponse
	// @Router /health [get]
	s.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, responses.HealthResponse{Status: "ok"})
	})

	s.engine.GET("/test", s.testHandler.Test)

	// API routes group
	api := s.engine.Group("/api")
	{

		// User routes
		users := api.Group("/users")
		{
			users.POST("", s.userHandler.CreateUser)
			users.GET("", s.userHandler.ListUsers)
		}

		// Protected user routes
		protected := users.Use(middleware.JWTAuth())
		{
			protected.GET("/:id", s.userHandler.GetUser)
			protected.PUT("/:id", s.userHandler.UpdateUser)
			protected.DELETE("/:id", s.userHandler.DeleteUser)
		}
	}
}

// Start starts the HTTP server with graceful shutdown
func (s *Server) Start() error {
	addr := ":" + s.config.Server.Port

	// Check if port is available
	if !isPortAvailable(s.config.Server.Port) {
		return fmt.Errorf("port %s is already in use", s.config.Server.Port)
	}

	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.engine,
	}

	// Start server in a goroutine
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	return nil
}

// Stop gracefully shuts down the server
func (s *Server) Stop() error {
	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown the HTTP server
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %v", err)
	}

	return nil
}

// Router returns the gin router instance
func (s *Server) Router() *gin.Engine {
	return s.engine
}

// isPortAvailable checks if a port is available
func isPortAvailable(port string) bool {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return false
	}
	ln.Close()
	return true
}
