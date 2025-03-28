package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	api := echo.New()

	api.Logger.SetLevel(log.DEBUG) // initialize inline logger
	api.Use(middleware.Recover())   // recover
	api.Use(middleware.RequestID()) // add request ID
	api.Use(middleware.Logger())
	api.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}","id":"${id}","remote_ip":"${remote_ip}",` +
			`"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
			`"status":${status},"error":"${error}","latency_human":"${latency_human}"}` + "\n",
	})) // initialise application-wide logger

	// Enable CORS
	api.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{echo.GET, echo.POST, echo.PATCH, echo.DELETE},
		AllowHeaders:     []string{echo.HeaderAuthorization, echo.HeaderContentType},
		AllowCredentials: true,
	}))

	api.POST("/register", register)
	api.POST("/login", login)
	api.GET("/register", registerValidation)
	api.GET("/authenticated", isAuthenticated)

	api.Start(":3130")
}
