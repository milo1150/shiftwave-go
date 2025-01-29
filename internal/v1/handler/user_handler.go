package handler

import (
	"net/http"
	"shiftwave-go/internal/auth"
	"shiftwave-go/internal/database"
	"shiftwave-go/internal/types"
	"shiftwave-go/internal/utils"
	"shiftwave-go/internal/v1/dto"
	v1repo "shiftwave-go/internal/v1/repository"
	"shiftwave-go/internal/validators"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
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
	token, err := auth.GenerateToken(app.ENV.JWT, user.Username, user.Uuid, user.Role, user.ActiveStatus)
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
func CreateUserHandler(c echo.Context, app *types.App) error {
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

func GetAllUsersHandler(c echo.Context, app *types.App) error {
	// Find users
	users, err := v1repo.GetAllUsers(app.DB)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, http.StatusInternalServerError)
	}

	// Transform Users
	userDtos := dto.TransformUserModels(*users)

	return c.JSON(http.StatusOK, userDtos)
}

func GetUserHandler(c echo.Context, app *types.App) error {
	// Extract user token from context
	jwtToken, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
	}

	// Extract claims from token
	claims, ok := jwtToken.Claims.(*auth.JwtCustomClaims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token claims"})
	}

	// Find users
	user, err := v1repo.GetUserByUUID(app.DB, claims.Uuid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, http.StatusInternalServerError)
	}

	// Transform Users
	userDto := dto.TransformUserModel(*user)

	return c.JSON(http.StatusOK, userDto)
}

func UpdateUsersHandler(c echo.Context, app *types.App) error {
	// Extract request payload
	payload := &[]types.UpdateUserPayload{}
	if err := c.Bind(payload); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid payload")
	}

	// Initialize Validator with custom rules
	v := validator.New()
	v.RegisterValidation("userRole", validators.ValidateUserRole)

	// Validate slice payload
	if err := validators.ValidateSlicePayload(c, v, payload); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, http.StatusOK)
}
