package middleware

import (
	"net/http"
	"time"

	"shiftwave-go/internal/types"
	"shiftwave-go/internal/utils"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

// jwtCustomClaims are custom claims extending default ones.
// See https://github.com/golang-jwt/jwt for more examples
type jwtCustomClaims struct {
	Name                 string `json:"name"`
	Admin                bool   `json:"admin"`
	jwt.RegisteredClaims        // struct embedding (in ts call extend interface)
}

func ConfigJWT(e *echo.Echo, secret string) echo.MiddlewareFunc {
	// WithConfig returns a JSON Web Token (JWT) auth middleware or panics if configuration is invalid.
	//
	// For valid token, it sets the user in context and calls next handler.
	// For invalid token, it returns "401 - Unauthorized" error.
	// For missing token, it returns "400 - Bad Request" error.
	//
	// See: https://jwt.io/introduction
	return echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(secret),

		// Context key to store user information from the token into context.
		// Optional. Default value "user".
		ContextKey: "user",

		// Extract jwt context in header into jwtCustomClaims struct
		// Example output:
		// Claims: (*middleware.jwtCustomClaims)(0x40003c2280)({
		//  Name: (string) (len=8) "John Doe",
		//  Admin: (bool) false,
		//  RegisteredClaims: (jwt.RegisteredClaims) {
		//   Issuer: (string) "",
		//   Subject: (string) (len=10) "1234567890",
		//   Audience: (jwt.ClaimStrings) <nil>,
		//   ExpiresAt: (*jwt.NumericDate)(<nil>),
		//   NotBefore: (*jwt.NumericDate)(<nil>),
		//   IssuedAt: (*jwt.NumericDate)(0x400000e090)(2018-01-18 01:30:22 +0000 UTC),
		//   ID: (string) ""
		//  }
		// })
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
	})
}

func JWT(e *echo.Echo, env types.Env) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Run JWT Extraction.
			// Call nested function [(next)(c)] until get an error.
			jwtMiddleware := ConfigJWT(e, env.JWT)
			if err := jwtMiddleware(next)(c); err != nil {
				return err
			}

			user := c.Get("user").(*jwt.Token)
			spew.Dump(user)

			return next(c)
		}
	}
}

func GenerateToken(secret string, name string, admin bool) (string, error) {
	// Set custom claims
	claims := &jwtCustomClaims{
		Name:  name,
		Admin: admin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Second)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return t, nil
}

func LoginHandler(c echo.Context, app *types.App) error {
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

	// TODO: Query -> Find user in db

	// Generate encoded token
	t, err := GenerateToken(app.ENV.JWT, "Min", true)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"token": t})
}

func GenerateUser(app *types.App) error {
	return nil
}

// TODO:
func HashPassword(password string) (string, error) {
	return "", nil
}
