// Package middleware exposes the JWT user-extraction helper used by API
// handlers. The skalpai convention is `middleware.GetUser(c)` — keep
// this signature stable so handlers stay portable across services.
package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type User struct {
	ID    uuid.UUID
	Email string
}

func GetUser(c echo.Context) (User, bool) {
	v := c.Get("user")
	if v == nil {
		return User{}, false
	}
	u, ok := v.(User)
	return u, ok
}
