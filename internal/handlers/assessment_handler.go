package handlers

import (
	"net/http"
	"shiftwave-go/internal/repositories"
	"shiftwave-go/internal/types"

	"github.com/labstack/echo"
)

func CreateAssessment(c echo.Context, app *types.App) error {
	result := repositories.CreateAssessment(app)
	if result != nil {
		return c.JSON(http.StatusInternalServerError, "error bruh")
	}
	return c.JSON(http.StatusOK, "")
}
