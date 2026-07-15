// Package middleware verifies the JWT minted by the web app (better-auth
// session → jose SignJWT) and exposes the skalpai convention `GetUser(c)`.
// The core never signs tokens; it only verifies the HS256 `token` cookie and
// reads the `sub` (better-auth user id, TEXT), `email` and `name` claims.
package middleware

import (
	"errors"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type User struct {
	ID    string
	Email string
	Name  string
}

type claims struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	jwt.RegisteredClaims
}

// RequireAuth verifies the `token` cookie (or Bearer header) against JWT_SECRET
// and stores the resolved User in the echo context. JWT_SECRET MUST match the
// web app's JWT_SECRET.
func RequireAuth() echo.MiddlewareFunc {
	secret := []byte(os.Getenv("JWT_SECRET"))
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			raw := tokenFromRequest(c)
			if raw == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "unauthenticated")
			}
			u, err := parse(raw, secret)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "unauthenticated")
			}
			c.Set("user", u)
			return next(c)
		}
	}
}

func tokenFromRequest(c echo.Context) string {
	if cookie, err := c.Cookie("token"); err == nil && cookie.Value != "" {
		return cookie.Value
	}
	if auth := c.Request().Header.Get("Authorization"); len(auth) > 7 && auth[:7] == "Bearer " {
		return auth[7:]
	}
	return ""
}

func parse(raw string, secret []byte) (User, error) {
	var cl claims
	tok, err := jwt.ParseWithClaims(raw, &cl, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	},
		// Belt-and-suspenders alg guard: reject anything but HS256 at the
		// parser level, before the keyfunc runs, so an alg-confusion token
		// (alg:none, RS256→HS256) can't slip through a future keyfunc edit.
		jwt.WithValidMethods([]string{"HS256"}),
		// Reject tokens without an exp claim outright, instead of treating a
		// missing expiry as "never expires".
		jwt.WithExpirationRequired(),
	)
	if err != nil || !tok.Valid || cl.Subject == "" {
		return User{}, errors.New("invalid token")
	}
	return User{ID: cl.Subject, Email: cl.Email, Name: cl.Name}, nil
}

func GetUser(c echo.Context) (User, bool) {
	v := c.Get("user")
	if v == nil {
		return User{}, false
	}
	u, ok := v.(User)
	return u, ok
}
