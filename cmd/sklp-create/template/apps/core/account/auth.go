package account

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// AuthLogout clears the auth session cookie. The token itself is minted and
// validated elsewhere (web better-auth → JWT the core verifies); logging out on
// the core side is just dropping the `token` cookie the browser sends.
//
// AuthLogout godoc
// @Summary  Clear the current auth session cookie
// @Tags     auth
// @Produce  json
// @Success  200  {object}  map[string]string
// @Router   /auth/session [delete]
// @ID       authLogout
func (s *Service) AuthLogout(c echo.Context) error {
	clearTokenCookie(c)
	return c.JSON(http.StatusOK, map[string]string{"status": "logged_out"})
}

// clearTokenCookie expires the `token` cookie. Secure stays false so it works
// over plain HTTP in local dev; front the core with TLS in production (the web
// app sets the cookie Secure there) and this expiry still applies.
func clearTokenCookie(c echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
}

// RegisterRoutes mounts the account/auth endpoints. Logout is mounted without
// RequireAuth: clearing a cookie must work even once the token is already
// expired or invalid.
func (s *Service) RegisterRoutes(e *echo.Echo) {
	e.DELETE("/auth/session", s.AuthLogout)
}
