package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/DaveSaah/auth-service/db"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// Secret key used for signing the token
var jwtKey = []byte(os.Getenv("SECRET"))

func login(c echo.Context) error {
	conn, err := db.Init()
	if err != nil {
		return err
	}
	defer conn.Close()

	q := db.New(conn)
	ctx := context.Background()
	var u db.User

	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid form data",
		})
	}

	user, err := q.GetUserByEmail(ctx, u.Email)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusPreconditionFailed, echo.Map{
			"message": "user does not exist",
		})
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password)); err != nil {
		return c.JSON(http.StatusPreconditionFailed, echo.Map{
			"message": "incorrect password",
		})
	}

	// user authentication successful

	// add jwt token
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // 1 day,
	})

	// generate JWT token
	token, err := claims.SignedString([]byte(jwtKey))
	if err != nil {
		return err
	}

	// store token in a cookie
	cookie := http.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	}
	c.SetCookie(&cookie)

	return c.JSON(http.StatusOK, echo.Map{
		"username": user.Username,
	})
}

func isAuthenticated(c echo.Context) error {
	// default is false
	authMsg := echo.Map{"authenticated": false}
	cookie, err := c.Cookie("jwt")

	if err != nil {
		c.JSON(http.StatusUnauthorized, authMsg)
		c.Logger().Info(authMsg)
		return err
	}

	_, err = jwt.ParseWithClaims(
		cookie.Value,
		&jwt.StandardClaims{},
		func(t *jwt.Token) (any, error) {
			return []byte(jwtKey), nil
		})
	if err != nil {
		c.JSON(http.StatusUnauthorized, authMsg)
		c.Logger().Info(authMsg)
		return err
	}

	// JWT valid
	authMsg["authenticated"] = true
	c.Logger().Info(authMsg)
	return c.JSON(http.StatusOK, authMsg)
}
