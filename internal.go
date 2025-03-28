package main

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/DaveSaah/auth-service/db"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// Secret key used for signing the token
var jwtKey = []byte(os.Getenv("SECRET"))

func getExpiryDate() time.Time {
	return time.Now().Add(24 * time.Hour) // Token expires in 24 hours
}

func generateJWT(u db.User) (string, error) {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(getExpiryDate()),
		Issuer:    strconv.Itoa(int(u.ID)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func setJWTInCookie(c echo.Context, token string) {
	cookie := new(http.Cookie)
	cookie.Name = "jwt"
	cookie.Value = token
	cookie.Expires = getExpiryDate()
	cookie.HttpOnly = true
	cookie.Secure = true
	cookie.Path = "/"
	cookie.SameSite = http.SameSiteNoneMode
	c.SetCookie(cookie)
}

func setJWT(c echo.Context, u db.User) error {
	token, err := generateJWT(u)
	if err != nil {
		return err
	}
	setJWTInCookie(c, token)
	return nil
}
