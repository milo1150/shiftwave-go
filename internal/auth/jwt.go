package auth

import (
	"shiftwave-go/internal/enum"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

// jwtCustomClaims are custom claims extending default ones.
// See https://github.com/golang-jwt/jwt for more examples
type JwtCustomClaims struct {
	Name                 string    `json:"name"`
	ID                   int       `json:"id"`
	ActiveStatus         bool      `json:"active_status"`
	Role                 enum.Role `json:"role"`
	jwt.RegisteredClaims           // struct embedding (in ts call extend interface)
}

func JwtConfig(secret string) echojwt.Config {
	return echojwt.Config{
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
			return new(JwtCustomClaims)
		},
	}
}

// Generate encoded jwt token for client
func GenerateToken(secret string, name string, id int, role enum.Role, status bool) (string, error) {
	// Set custom claims
	claims := &JwtCustomClaims{
		Name:         name,
		ID:           id,
		Role:         role,
		ActiveStatus: status,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
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
