package auth

import (
	"time"

	"shiftwave-go/internal/types"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

// jwtCustomClaims are custom claims extending default ones.
// See https://github.com/golang-jwt/jwt for more examples
type jwtCustomClaims struct {
	Name string `json:"name"`
	// Admin                bool   `json:"admin"` // TODO: permission ?
	jwt.RegisteredClaims // struct embedding (in ts call extend interface)
}

func configJWT(secret string) echo.MiddlewareFunc {
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

func Jwt(e *echo.Echo, env types.Env) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Run JWT Extraction.
			// Call nested function [(next)(c)] until get an error.
			jwtMiddleware := configJWT(env.JWT)
			if err := jwtMiddleware(next)(c); err != nil {
				return err
			}

			// Save for debug
			// user := c.Get("user").(*jwt.Token)
			// spew.Dump(user)

			return next(c)
		}
	}
}

// Generate encoded jwt token for client
func GenerateToken(secret string, name string) (string, error) {
	// Set custom claims
	claims := &jwtCustomClaims{
		Name: name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(20 * time.Second)),
		},
	}

	// Create token with claims
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token
	token, err := jwtToken.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return token, nil
}
