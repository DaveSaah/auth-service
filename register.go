package main

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/DaveSaah/auth-service/db"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func registerValidation(c echo.Context) error {
	params := c.QueryParams() // get all query params

	if len(params) == 0 {
		return nil // skip validation if there are no params
	}

	conn, err := db.Init()
	if err != nil {
		return err
	}
	defer conn.Close()

	q := db.New(conn) // query interface
	ctx := context.Background()

	validationChecks := map[string]func(context.Context, string) (db.User, error){
		"username": q.GetUserByUsername,
		"email":    q.GetUserByEmail,
	}

	for key, checkFunc := range validationChecks {
		if value, exists := params[key]; exists {
			if _, err := checkFunc(ctx, value[0]); err != sql.ErrNoRows {
				return c.JSON(http.StatusPreconditionFailed, echo.Map{
					"message": key + " already exists",
				})
			}
		}
	}

	return nil
}

func register(c echo.Context) error {
	conn, err := db.Init()
	if err != nil {
		return err
	}
	defer conn.Close()

	q := db.New(conn) // query interface
	ctx := context.Background()
	var u db.User

	if err := c.Bind(&u); err != nil {
		c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid form data",
		})
	}

	hash_passwd, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = q.CreateUser(ctx, db.CreateUserParams{
		Username: u.Username,
		Email:    u.Email,
		Password: string(hash_passwd),
	})
	return err
}
