package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/dgrijalva/jwt-go"
)

func authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("jwt")
		if err != nil {
			c.Logger().Info("Unauthorized: No JWT found")
			return c.JSON(http.StatusUnauthorized, echo.Map{"authenticated": false})
		}

		claims := &jwt.StandardClaims{}
		_, err = jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (any,error) {
			return jwtKey, nil
		})

		if err != nil {
			c.Logger().Info("Unauthorized: Invalid JWT")
			return c.JSON(http.StatusUnauthorized, echo.Map{"authenticated": false})
		}

		// authentication successful, store user info in context
		c.Set("authenticated", true)
		c.Set("userID", claims.Issuer) 

		c.Logger().Info(echo.Map{"authenticated": true, "userID": claims.Issuer})
		return next(c) // Continue to the next middleware/handler
	}
}