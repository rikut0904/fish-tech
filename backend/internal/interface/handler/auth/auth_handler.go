package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// AuthHandler handles authenticated requests (Firebase tokens are verified
// by middleware before handlers are called).

type AuthHandler struct {
}

// NewAuthHandler constructs an AuthHandler.
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

// GetProfile returns basic information about the authenticated user.
func (h *AuthHandler) GetProfile(c echo.Context) error {
	// the middleware stores the decoded token under "firebaseUser"
	decoded := c.Get("firebaseUser")
	if decoded == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "user information missing")
	}

	return c.JSON(http.StatusOK, decoded)
}
