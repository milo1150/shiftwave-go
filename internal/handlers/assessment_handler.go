package handlers

import (
	"fmt"
	"net/http"
	"shiftwave-go/internal/repositories"
	"shiftwave-go/internal/types"
	"shiftwave-go/internal/utils"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
)

func CreateAssessmentHandler(c echo.Context, app *types.App) error {
	payload := new(types.CreateAssessmentPayload)
	if err := c.Bind(payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
	}

	v := validator.New()
	if err := v.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMessagees := utils.ExtractErrorMessages(validationErrors)
		return c.JSON(http.StatusBadRequest, errorMessagees)
	}

	if result := repositories.CreateAssessment(app, payload); result != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": result.Error()})
	}

	return c.JSON(http.StatusOK, "OK")
}

func GetAssessmentsHandler(c echo.Context, app *types.App) error {
	q := &types.AssessmentQueryParams{}
	if err := c.Bind(q); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid Query")
	}

	if q.Page != nil && *q.Page > 0 {
		fmt.Println("queryparams -> page = ", *q.Page)
	}
	if q.PageSize != nil && *q.PageSize > 0 {
		fmt.Println("queryparams -> page_size = ", *q.PageSize)
	}

	v := validator.New()
	if err := v.Struct(q); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMessages := utils.ExtractErrorMessages(validationErrors)
		return c.JSON(http.StatusBadRequest, errorMessages)
	}

	result, err := repositories.GetAssessments(app, q)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Query error")
	}

	return c.JSON(http.StatusOK, result)
}

func GetAssessment(c echo.Context, app *types.App) error {
	param := c.Param("id")

	id, err := strconv.Atoi(param)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid param")
	}

	result, _ := repositories.GetAssessment(app, id)

	return c.JSON(http.StatusOK, result)
}
