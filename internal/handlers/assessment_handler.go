package handlers

import (
	"net/http"
	"shiftwave-go/internal/repositories"
	"shiftwave-go/internal/types"
	"shiftwave-go/internal/utils"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
)

func CreateAssessmentHandler(c echo.Context, app *types.App) error {
	payload := new(types.CreateAssessmentPayload)
	if err := c.Bind(payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
	}

	validate := validator.New()
	if err := validate.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMessagees := utils.ExtractErrorMessages(validationErrors)
		return c.JSON(http.StatusBadRequest, errorMessagees)
	}

	if result := repositories.CreateAssessment(app, payload); result != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": result.Error()})
	}

	return c.JSON(http.StatusOK, "OK")
}
