package handler

import (
	"log"
	"net/http"
	"os"
	"shiftwave-go/internal/services"
	"shiftwave-go/internal/types"
	"shiftwave-go/internal/utils"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func GenerateQRCodeHandler(c echo.Context) error {
	q := &types.GeneratePdfParams{}
	if err := c.Bind(q); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid Query")
	}

	v := validator.New()
	if err := v.Struct(q); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMessages := utils.ExtractErrorMessages(validationErrors)
		return c.JSON(http.StatusBadRequest, errorMessages)
	}

	_, err := godotenv.Read()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	m := services.GenerateReviewQRcode(os.Getenv("BASE_URL"), q.BranchId)

	document, err := m.Generate()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to generate pdf.")
	}

	return c.Blob(http.StatusOK, "application/pdf", document.GetBytes())
}
