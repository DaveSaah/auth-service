package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func logout(c echo.Context) error {
	cookie := http.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
		Path:     "/", // Ensure it applies to the whole domain
	}
	c.SetCookie(&cookie)

	c.Logger().Info("user logged out")

	return c.JSON(http.StatusOK, echo.Map{
		"message": "You have been logged out successfully.",
	})
}
