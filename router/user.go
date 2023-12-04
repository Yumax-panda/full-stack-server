package router

import (
	"full-stack-server/model"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetUsersHandler(c echo.Context) error {
	users, err := model.GetUser()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, users)
}

type CreateUserPayload struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Comment  string `json:"comment"`
}

func CreateUserHandler(c echo.Context) error {
	var payload CreateUserPayload
	err := c.Bind(&payload)
	user, err := model.CreateUser(payload.Name, payload.Email, payload.Password, payload.Comment)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, user)
}

type LoginUserPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginUserHandler(c echo.Context) error {
	var payload LoginUserPayload
	err := c.Bind(&payload)
	user, err := model.Login(payload.Email, payload.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, user)
}
