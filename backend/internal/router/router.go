package router

import (
	"github.com/gofiber/fiber/v2"

	"github.com/koc-luk/backend/internal/auth"
	"github.com/koc-luk/backend/internal/domain"
	"github.com/koc-luk/backend/internal/handler"
	"github.com/koc-luk/backend/internal/middleware"
)

type Deps struct {
	Health        *handler.HealthHandler
	Auth          *handler.AuthHandler
	Registrations *handler.RegistrationHandler
	Users         *handler.UserHandler
	Assignments   *handler.AssignmentHandler
	Weeks         *handler.WeekHandler
	Evaluations   *handler.EvaluationHandler
	Messaging     *handler.MessagingHandler
	Sms           *handler.SmsHandler
	Reports       *handler.ReportHandler
	JWT           *auth.Manager
}

func Register(app *fiber.App, deps Deps) {
	v1 := app.Group("/api/v1")

	v1.Get("/health", deps.Health.Get)

	v1.Post("/auth/login", deps.Auth.Login)
	v1.Post("/registrations/student", deps.Registrations.ApplyStudent)
	v1.Post("/registrations/coach", deps.Registrations.ApplyCoach)

	authReq := v1.Group("", middleware.JWTAuth(deps.JWT))
	authReq.Get("/auth/me", deps.Auth.Me)

	admin := authReq.Group("/admin", middleware.RequireRole(domain.RoleAdmin))

	admin.Get("/registrations", deps.Registrations.List)
	admin.Post("/registrations/:id/decision", deps.Registrations.Decide)

	admin.Get("/users", deps.Users.List)
	admin.Post("/users", deps.Users.Create)
	admin.Get("/users/:id", deps.Users.Get)
	admin.Patch("/users/:id", deps.Users.Update)
	admin.Post("/users/:id/reset-password", deps.Users.ResetPassword)

	admin.Get("/assignments", deps.Assignments.List)
	admin.Post("/assignments", deps.Assignments.Assign)
	admin.Post("/assignments/unassign", deps.Assignments.Unassign)
	admin.Get("/assignments/students-without-coach", deps.Assignments.StudentsWithoutCoach)
	admin.Post("/assignments/set-coordinator", deps.Assignments.SetCoordinator)
	admin.Post("/assignments/set-parent", deps.Assignments.SetParent)

	admin.Get("/weeks", deps.Weeks.ListAll)
	admin.Post("/weeks", deps.Weeks.Create)
	admin.Patch("/weeks/:id", deps.Weeks.Update)
	admin.Post("/weeks/:id/open", deps.Weeks.Open)
	admin.Post("/weeks/:id/close", deps.Weeks.Close)
	admin.Post("/weeks/generate", deps.Weeks.Generate)

	authReq.Get("/weeks/open", deps.Weeks.ListOpen)

	authReq.Get("/evaluations", deps.Evaluations.List)
	authReq.Get("/evaluations/:id", deps.Evaluations.Get)

	coach := authReq.Group("/coach", middleware.RequireRole(domain.RoleCoach))
	coach.Get("/students", deps.Evaluations.MyStudents)
	coach.Post("/evaluations", deps.Evaluations.CoachCreate)
	coach.Patch("/evaluations/:id", deps.Evaluations.CoachUpdate)

	admin.Patch("/evaluations/:id", deps.Evaluations.AdminUpdate)
	admin.Post("/evaluations/:id/approve", deps.Evaluations.Approve)
	admin.Post("/evaluations/:id/revoke", deps.Evaluations.Revoke)
	admin.Get("/evaluations/:id/versions", deps.Evaluations.Versions)

	authReq.Get("/threads", deps.Messaging.ListThreads)
	authReq.Post("/threads", deps.Messaging.Create)
	authReq.Get("/threads/:id/messages", deps.Messaging.Messages)
	authReq.Post("/threads/:id/messages", deps.Messaging.Send)
	authReq.Post("/threads/:id/close", deps.Messaging.Close)

	admin.Get("/sms/templates", deps.Sms.Templates)
	admin.Post("/sms/send", deps.Sms.Send)
	admin.Get("/sms/logs", deps.Sms.Logs)
	admin.Get("/sms/missing-coaches", deps.Sms.MissingCoaches)

	admin.Get("/reports/summary", deps.Reports.Summary)
	admin.Get("/reports/coach-performance", deps.Reports.CoachPerformance)
	admin.Get("/reports/city-distribution", deps.Reports.CityDistribution)
	admin.Get("/reports/week-stats", deps.Reports.WeekStats)
	authReq.Get("/reports/student-trend", deps.Reports.StudentTrend)
}
