package router

import (
	"full-stack-server/model"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
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

func GetUsers(c echo.Context) error {
	users, err := model.GetUsers()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, users)
}

func CreateUser(c echo.Context) error {
	var payload CreateUserPayload
	err := c.Bind(&payload)
	user, err := model.CreateUser(payload.Name, payload.Email, payload.Password, payload.Comment)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, user)
}

func setCookie(c echo.Context, t string) {
	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    t,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	}
	c.SetCookie(cookie)
}

func Login(c echo.Context) error {
	var payload LoginUserPayload
	err := c.Bind(&payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := model.GetUserByEmail(payload.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid password")
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

func Logout(c echo.Context) error {
	newCookie := &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}
	c.SetCookie(newCookie)
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Logged out",
	})
}

func GetCurrentUser(c echo.Context) error {
	cookie, err := c.Cookie("jwt")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "No cookie")
	}
	token, err := jwt.ParseWithClaims(cookie.Value, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized: ", err.Error())
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized: ", err.Error())
	}
	user, err := model.GetUserByID(claims.Issuer)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized: ", err.Error())
	}
	return c.JSON(http.StatusOK, user)
}
