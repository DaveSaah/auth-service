package main

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/DaveSaah/auth-service/db"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

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
	// set jwt token
	err = setJWT(c, u)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"username": user.Username,
	})
}

