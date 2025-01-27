package middleware

import (
	"net/http"
	"shiftwave-go/internal/auth"
	"shiftwave-go/internal/database"
	"shiftwave-go/internal/types"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

// Validate JWT token and bind in echo.Context
func ValidateJwt(e *echo.Echo, app *types.App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cfg := auth.JwtConfig(app.ENV.JWT)
			jwtHandler := echojwt.WithConfig(cfg)
			if err := jwtHandler(next)(c); err != nil {
				return err
			}
			return nil
		}
	}
}

func CheckExistedJwtToken(app *types.App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cfg := auth.JwtConfig(app.ENV.JWT)

			jwt, ok := c.Get(cfg.ContextKey).(*jwt.Token)
			if !ok {
				return c.JSON(http.StatusUnauthorized, "Token not found")
			}

			currentToken := jwt.Raw
			username := jwt.Claims.(*auth.JwtCustomClaims).Name
			key := database.GetJwtKey(username)

			// Find cache token in redis
			existedToken, err := app.RDB.Get(c.Request().Context(), key).Result()
			if err != nil {
				return c.JSON(http.StatusUnauthorized, "Unauthorized token")
			}

			if currentToken != existedToken {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Your token is expired. Please login again."})
			}

			return next(c)
		}
	}

}
