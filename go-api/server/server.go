package server

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-api/config"
	"go-api/controllers"
	"go-api/repositories"
	"go-api/usecases"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Server struct {
	app *fiber.App
	cfg *config.Config
}

func New(cfg *config.Config) *Server {
	// Inyección de dependencias.
	authUsecase := usecases.NewAuthUseCase(cfg.JWTSecret, cfg.AuthUser, cfg.AuthPass, cfg.TokenExpiry())
	authController := controllers.NewAuthHandler(authUsecase)

	qrUsecase := usecases.NewQRUseCase()
	statsRepository := repositories.NewStatsRepository(cfg.NodeAPIURL)
	processUsecase := usecases.NewProcessUseCase(qrUsecase, statsRepository)
	matrixController := controllers.NewMatrixHandler(processUsecase)

	// App Fiber + middlewares globales.
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORSOrigins,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,OPTIONS",
	}))

	registerRoutes(app, cfg, authController, matrixController)

	return &Server{app: app, cfg: cfg}
}

// Start arranca el servidor y lo apaga de forma ordenada al recibir SIGINT/SIGTERM.
func (s *Server) Start() error {
	go func() {
		addr := ":" + s.cfg.ServerPort
		log.Printf("Servidor Go (QR + estadísticas) escuchando en %s (Node API: %s)", addr, s.cfg.NodeAPIURL)
		if err := s.app.Listen(addr); err != nil {
			log.Fatalf("Error al iniciar el servidor: %v", err)
		}
	}()

	// Espera señal del sistema para apagar.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Apagando el servidor...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.app.ShutdownWithContext(ctx)
}
