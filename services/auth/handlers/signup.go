package handlers

import (
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

func (h *Handler) SignUp(c echo.Context) error {
	var reqBody struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := c.Bind(&reqBody)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if reqBody.Name == "" || reqBody.Email == "" || reqBody.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "input cannot empty")
	}

	user := models.User{
		Name:  reqBody.Name,
		Email: reqBody.Email,
	}
	var count int64
	h.db.Model(&models.User{}).Where("email = ?", user.Email).Count(&count)
	if count > 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "อีเมล์นี้ถูกใช้แล้ว กรุณาใช้อีเมล์อื่น")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	user.Password = string(hashedPassword)
	user.Avatar = "img-url"
	h.db.Create(&user)
	// User created publish
	err = h.ec.Publish("user:created", user)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	//Email created publish
	msg := map[string]string{
		"Email":   user.Email,
		"Subject": "Thank you for registering an account!",
		"Text":    "Hello " + user.Name + ". Thank you for registering an account with pay9.co!",
	}
	err = h.ec.Publish("email:created", msg)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
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
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	sess, err := session.Get("pay9.sess", c)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	sess.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
	}
	sess.Values["jwt"] = tokenString
	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, user)
}
