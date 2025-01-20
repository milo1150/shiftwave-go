package handler

import (
	"net/http"
	"shiftwave-go/internal/auth"
	"shiftwave-go/internal/types"
	"shiftwave-go/internal/utils"
	v1repo "shiftwave-go/internal/v1/repository"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
)

func LoginHandler(c echo.Context, app *types.App) error {
	// Extract json payload
	payload := &types.LoginPayload{}
	if err := c.Bind(payload); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid payload")
	}

	// Validate login payload
	v := validator.New()
	if err := v.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMessages := utils.ExtractErrorMessages(validationErrors)
		return c.JSON(http.StatusBadRequest, errorMessages)
	}

	// Find user in db
	user, err := v1repo.FindUser(app.DB, payload.Username)
	if err != nil {
		return err
	}

	result := auth.ComparePassword(user.Password, payload.Password)
	if !result {
		return c.JSON(http.StatusUnauthorized, "WTF man")
	}

	// Generate jwt token
	token, err := auth.GenerateToken(app.ENV.JWT, user.Username)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}

// CreateUser creates a new user
// @Summary Create a new user
// @Description API endpoint to create a new user
// @Tags Users
// @Accept json
// @Produce json
// @Param payload body types.CreateUserPayload true "Payload Create User"
// @Router /user [post]
func CreateUser(c echo.Context, app *types.App) error {
	// Extract json payload
	payload := &types.CreateUserPayload{}
	if err := c.Bind(payload); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid payload")
	}

	// Validate login payload
	v := validator.New()
	if err := v.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMessages := utils.ExtractErrorMessages(validationErrors)
		return c.JSON(http.StatusBadRequest, errorMessages)
	}

	// Handle create user and error
	err := v1repo.CreateUser(app.DB, payload)
	if pgErr, ok := err.(*pgconn.PgError); ok { // see type assertion in go spec
		return c.JSON(http.StatusConflict, pgErr.Detail)
	}
	if err != nil {
		return c.JSON(http.StatusConflict, err.Error())
	}

	return c.JSON(http.StatusCreated, http.StatusCreated)
}
