package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/haytamxp/redlab/backend/internal/auth"
	"github.com/haytamxp/redlab/backend/internal/database"
	"github.com/haytamxp/redlab/backend/internal/dto"
	"github.com/haytamxp/redlab/backend/internal/handlers"
	"github.com/haytamxp/redlab/backend/internal/permissions"
	"github.com/haytamxp/redlab/backend/internal/reporting"
	"github.com/haytamxp/redlab/backend/internal/repository"
	"github.com/haytamxp/redlab/backend/internal/services"
)

func (r *Router) RegisterRoutes(
	authHandler *handlers.AuthHandler,
	agentHandler *handlers.AgentHandler,
	ldapService *services.LDAPService,
	jwtSecret string,
) {
	/*
	 * Public health.
	 */

	r.Engine.GET(
		"/health",
		handlers.Health,
	)

	api := r.Engine.Group(
		"/api",
	)

	v1 := api.Group(
		"/v1",
	)

	v1.GET(
		"/health",
		handlers.Health,
	)

	/*
	 * Authentication.
	 */

	authGroup :=
		v1.Group("/auth")

	authGroup.POST(
		"/register",
		authHandler.Register,
	)

	authGroup.POST(
		"/login",
		authHandler.Login,
	)

	/*
	 * Services.
	 */

	agentRepository :=
		repository.NewAgentRepository(
			database.DB,
		)

	agentService :=
		services.NewAgentService(
			agentRepository,
		)

	taskRepository :=
		repository.NewTaskRepository(
			database.DB,
		)

	findingRepository :=
		repository.NewFindingRepository(
			database.DB,
		)

	findingService :=
		services.NewFindingService(
			findingRepository,
		)

	taskService :=
		services.NewTaskService(
			taskRepository,
			agentService,
		)

	taskHandler :=
		handlers.NewTaskHandler(
			taskService,
			agentService,
			findingService,
		)

	assessmentRepository :=
		repository.NewAssessmentRepository(
			database.DB,
		)

	assessmentService :=
		services.NewAssessmentService(
			assessmentRepository,
			agentService,
		)

	assessmentHandler :=
		handlers.NewAssessmentHandler(
			assessmentService,
		)

	reportingService :=
		reporting.NewService(
			database.DB,
		)

	reportingHandler :=
		reporting.NewHandler(
			reportingService,
		)

	/*
	 * Agent-authenticated routes.
	 */

	v1.POST(
		"/agents/:id/tasks/next",
		taskHandler.Next,
	)

	v1.GET(
		"/agents/:id/tasks",
		taskHandler.ListForAgent,
	)

	v1.POST(
		"/agents/:id/tasks/:taskId/result",
		taskHandler.Complete,
	)

	v1.POST(
		"/agents/:id/heartbeat",
		agentHandler.Heartbeat,
	)

	/*
	 * Human-authenticated routes.
	 */

	protected :=
		v1.Group("/")

	protected.Use(
		auth.JWTMiddleware(
			jwtSecret,
		),
	)

	/*
	 * Profile.
	 */

	protected.GET(
		"/profile",
		func(c *gin.Context) {
			c.JSON(
				http.StatusOK,
				gin.H{
					"success": true,
					"message": "Authenticated",
					"user_id": c.GetString(
						"userID",
					),
					"role": c.GetString(
						"role",
					),
				},
			)
		},
	)

	/*
	 * Dashboard.
	 */

	protected.GET(
		"/dashboard/stats",
		handlers.DashboardStatsHandler,
	)

	directory := protected.Group("/directory")
	directory.Use(auth.RequirePermission(permissions.LDAPSync))
	directory.GET(
		"/users",
		handlers.NewDirectoryHandler(ldapService).ListUsers,
	)

	/*
	 * Agents.
	 */

	agents :=
		protected.Group(
			"/agents",
		)

	agents.Use(
		auth.RequirePermission(
			permissions.ManageAgents,
		),
	)

	agents.POST(
		"",
		agentHandler.Create,
	)

	agents.GET(
		"",
		agentHandler.List,
	)

	agents.GET(
		"/:id",
		agentHandler.Get,
	)

	/*
	 * Tasks.
	 */

	tasks :=
		protected.Group(
			"/tasks",
		)

	tasks.Use(
		auth.RequirePermission(
			permissions.ManageAgents,
		),
	)

	tasks.POST(
		"",
		taskHandler.Create,
	)

	tasks.GET(
		"",
		taskHandler.ListAll,
	)

	tasks.DELETE(
		"/:id",
		taskHandler.Delete,
	)
	tasks.PATCH(
		"/:id/review",
		taskHandler.Review,
	)

	/*
	 * Assessments.
	 */

	assessments :=
		protected.Group(
			"/assessments",
		)

	assessments.Use(
		auth.RequirePermission(
			permissions.ManageAgents,
		),
	)

	assessments.POST(
		"",
		assessmentHandler.Create,
	)

	assessments.GET(
		"",
		assessmentHandler.List,
	)

	assessments.GET(
		"/:id/report",
		reportingHandler.Get,
	)

	assessments.GET(
		"/:id",
		assessmentHandler.Get,
	)

	assessments.PATCH(
		"/:id/status",
		assessmentHandler.UpdateStatus,
	)

	/*
	 * Findings.
	 */

	findings :=
		protected.Group(
			"/findings",
		)

	findings.Use(
		auth.RequirePermission(
			permissions.ManageAgents,
		),
	)

	findings.GET(
		"",
		handlers.NewFindingHandler(
			findingService,
		).List,
	)

	findings.GET(
		"/:id",
		handlers.NewFindingHandler(
			findingService,
		).Get,
	)

	_ = dto.TaskResponse{}
}
