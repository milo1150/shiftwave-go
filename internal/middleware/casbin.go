package middleware

import (
	"log"
	"net/http"
	"shiftwave-go/internal/auth"

	"github.com/casbin/casbin/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func RoutePermission(secretJwt string, e *casbin.Enforcer) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Check header
			user, ok := c.Get("user").(*jwt.Token)
			if !ok {
				return echo.ErrForbidden
			}

			// Extract value
			claim := user.Claims.(*auth.JwtCustomClaims)
			path := c.Request().URL.Path
			role := claim.Role

			// Validate route
			allowed, err := e.Enforce(role, path)
			if err != nil {
				log.Println("Enfore error:", err) // TODO: use zap
				return echo.ErrInternalServerError
			}

			// Throw error if permission not allowed
			if !allowed {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "Permission denied"})
			}

			return next(c)
		}
	}
}
