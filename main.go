package main

import (
	"github.com/labstack/echo/v4"
)

func main() {
	api := echo.New()

	api.POST("/register", register)
	api.GET("/register", registerValidation)

	api.Logger.Fatal(api.Start(":3130"))
}
