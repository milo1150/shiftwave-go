package handler

import (
	"net/http"
	"shiftwave-go/internal/auth"
	"shiftwave-go/internal/database"
	"shiftwave-go/internal/types"
	"shiftwave-go/internal/utils"
	v1repo "shiftwave-go/internal/v1/repository"
	"shiftwave-go/internal/validators"
	"time"

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
		return c.JSON(http.StatusBadRequest, "Invalid username or password")
	}

	// Validate password
	result := auth.ComparePassword(user.Password, payload.Password)
	if !result {
		return c.JSON(http.StatusBadRequest, "Invalid username or password")
	}

	// Generate jwt token
	token, err := auth.GenerateToken(app.ENV.JWT, user.Username, int(user.ID), user.Role, user.ActiveStatus)
	if err != nil {
		return err
	}

	// Store token in redis
	rdJwtKey := database.GetJwtKey(user.Username)
	r := app.RDB.Set(c.Request().Context(), rdJwtKey, token, 24*time.Hour)
	if r.Err() != nil {
		return echo.ErrInternalServerError
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

	// Initialize Validator with custom rules
	v := validator.New()
	v.RegisterValidation("userRole", validators.ValidateUserRole)

	// Validate login payload
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
