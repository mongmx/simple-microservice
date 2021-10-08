package handlers

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/mongmx/simple-microservice/services/auth/models"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) SignIn(c echo.Context) error {
	var reqBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := c.Bind(&reqBody)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if reqBody.Email == "" || reqBody.Password == "" {
		return c.JSON(http.StatusBadRequest, errors.New("input cannot empty"))
	}

	var user models.User
	tx := h.db.Where("EMAIL = ?", reqBody.Email).First(&user)
	if tx.Error != nil {
		return c.JSON(http.StatusOK, "Could not find current user")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.Password))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "Could not find current user")
	}
	claims := &JwtCustomClaims{
		user.ID,
		user.Email,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	sess, err := session.Get("pay9.sess", c)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	sess.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
	}
	sess.Values["jwt"] = tokenString
	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, user)
}
