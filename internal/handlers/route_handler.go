package handlers

import (
	"shiftwave-go/internal/types"

	"github.com/labstack/echo"
)

func SetupRoutes(e *echo.Echo, app *types.App) {
	e.POST("/assessment", func(ctx echo.Context) error {
		return CreateAssessmentHandler(ctx, app)
	})
	e.GET("/assessment", func(ctx echo.Context) error {
		return GetAssessmentsHandler(ctx, app)
	})
}
