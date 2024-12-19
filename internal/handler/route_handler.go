package handler

import (
	"shiftwave-go/internal/types"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, app *types.App) {
	e.POST("/assessment", func(ctx echo.Context) error {
		return CreateAssessmentHandler(ctx, app)
	})
	e.GET("/assessments", func(ctx echo.Context) error {
		return GetAssessmentsHandler(ctx, app)
	})
	e.GET("/assessment/:id", func(ctx echo.Context) error {
		return GetAssessmentHandler(ctx, app)
	})
}
