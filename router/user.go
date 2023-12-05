package router

import (
	"full-stack-server/model"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type CreateUserPayload struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Comment  string `json:"comment"`
}

type LoginUserPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	jwt.RegisteredClaims
}

func generateToken(u model.User) (string, error) {
	claims := &jwt.RegisteredClaims{
		Issuer:    u.ID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	}
	token := jwt.NewWithClaims((jwt.SigningMethodHS256), claims)
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}
	return t, nil
}

func GetUsersHandler(c echo.Context) error {
	users, err := model.GetUser()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, users)
}

func CreateUserHandler(c echo.Context) error {
	var payload CreateUserPayload
	err := c.Bind(&payload)
	user, err := model.CreateUser(payload.Name, payload.Email, payload.Password, payload.Comment)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	t, err := generateToken(user)
	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

func setCookie(c echo.Context, t string) {
	cookie := new(http.Cookie)
	cookie.Name = "jwt"
	cookie.Value = t
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.HttpOnly = true
	c.SetCookie(cookie)
}

func LoginUserHandler(c echo.Context) error {
	var payload LoginUserPayload

	err := c.Bind(&payload)
	user, err := model.Login(payload.Email, payload.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	t, err := generateToken(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	setCookie(c, t)
	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}
