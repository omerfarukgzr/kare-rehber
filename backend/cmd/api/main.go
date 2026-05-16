package main

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/koc-luk/backend/internal/auth"
	"github.com/koc-luk/backend/internal/config"
	"github.com/koc-luk/backend/internal/handler"
	"github.com/koc-luk/backend/internal/repository"
	"github.com/koc-luk/backend/internal/router"
	"github.com/koc-luk/backend/internal/service"
	smspkg "github.com/koc-luk/backend/internal/sms"
	apperrors "github.com/koc-luk/backend/pkg/errors"
	applogger "github.com/koc-luk/backend/pkg/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load error: %v", err)
	}

	lg := applogger.New(cfg.AppEnv)

	db, err := config.OpenDB(cfg)
	if err != nil {
		lg.Error("db open error", "err", err)
		log.Fatalf("db open error: %v", err)
	}

	jwtMgr := auth.NewManager(cfg.JWTSecret, cfg.JWTExpiresIn)
	validate := validator.New(validator.WithRequiredStructEnabled())

	smsProvider := smspkg.NewMockProvider(db)

	userRepo := repository.NewUserRepository(db)
	regRepo := repository.NewRegistrationRepository(db)
	assignmentRepo := repository.NewAssignmentRepository(db)
	weekRepo := repository.NewWeekRepository(db)
	evalRepo := repository.NewEvaluationRepository(db)
	messagingRepo := repository.NewMessagingRepository(db)
	smsRepo := repository.NewSmsRepository(db)

	authSvc := service.NewAuthService(userRepo, jwtMgr)
	regSvc := service.NewRegistrationService(db, regRepo, userRepo, smsProvider, cfg)
	userSvc := service.NewUserService(db, userRepo, smsProvider, cfg)
	assignmentSvc := service.NewAssignmentService(db, assignmentRepo, userRepo)
	weekSvc := service.NewWeekService(db, weekRepo)
	evalSvc := service.NewEvaluationService(db, evalRepo, weekRepo, assignmentRepo, userRepo)
	coachSvc := service.NewCoachService(db)
	messagingSvc := service.NewMessagingService(db, messagingRepo, userRepo)
	smsSvc := service.NewSmsService(smsProvider, userRepo, evalRepo, smsRepo)
	reportSvc := service.NewReportService(db)

	app := fiber.New(fiber.Config{
		AppName:      "kare-rehber-api",
		ErrorHandler: apperrors.Handler,
	})

	app.Use(recover.New())
	app.Use(fiberlogger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSAllowedOrigins,
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	}))

	deps := router.Deps{
		Health:        handler.NewHealthHandler(db),
		Auth:          handler.NewAuthHandler(authSvc, userRepo, validate),
		Registrations: handler.NewRegistrationHandler(regSvc, validate),
		Users:         handler.NewUserHandler(userSvc, validate),
		Assignments:   handler.NewAssignmentHandler(assignmentSvc, validate),
		Weeks:         handler.NewWeekHandler(weekSvc, validate),
		Evaluations:   handler.NewEvaluationHandler(evalSvc, coachSvc, validate),
		Messaging:     handler.NewMessagingHandler(messagingSvc, validate),
		Sms:           handler.NewSmsHandler(smsSvc, validate),
		Reports:       handler.NewReportHandler(reportSvc),
		JWT:           jwtMgr,
	}
	router.Register(app, deps)

	addr := ":" + cfg.HTTPPort
	lg.Info("server starting", "addr", addr, "env", cfg.AppEnv)
	if err := app.Listen(addr); err != nil {
		lg.Error("server failed", "err", err)
		log.Fatalf("server failed: %v", err)
	}
}
