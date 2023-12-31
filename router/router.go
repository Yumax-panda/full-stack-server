package router

import (
	"os"

	"github.com/labstack/echo/v4/middleware"

	_ "net/http"

	"github.com/labstack/echo/v4"
)

func SetRouter(e *echo.Echo) error {

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339_nano} ${host} ${method} ${uri} ${status} ${header}\n",
		Output: os.Stdout,
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	}))

	api := e.Group("/api/v1")
	{
		apiUsers := api.Group("/users")
		{
			apiUsers.GET("", GetUsers)
			apiUsers.POST("", CreateUser)
			apiUsers.POST("/login", Login)
			apiUsers.GET("/logout", Logout)
			apiMe := apiUsers.Group("/@me")
			{
				apiMe.GET("", GetCurrentUser)
			}
		}
	}

	err := e.Start(":8000")
	return err
}
