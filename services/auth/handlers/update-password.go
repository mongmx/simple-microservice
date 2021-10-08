package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mongmx/simple-microservice/services/auth/models"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) UpdatePassword(c echo.Context) error {
	var reqBody struct {
		OldPassword     string `json:"oldPassword"`
		NewPassword     string `json:"newPassword"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	err := c.Bind(&reqBody)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if reqBody.OldPassword == "" || reqBody.NewPassword == "" || reqBody.ConfirmPassword == "" {
		return c.JSON(http.StatusBadRequest, errors.New("input cannot empty"))
	}

	if reqBody.NewPassword != reqBody.ConfirmPassword {
		return c.JSON(http.StatusBadRequest, errors.New("passwords do not match"))
	}
	currentUser, ok := c.Get("currentUser").(*JwtCustomClaims)
	if !ok {
		return c.JSON(http.StatusOK, "Could not find current user")
	}
	var user models.User
	tx := h.db.Where("ID = ?", currentUser.ID).First(&user)
	if tx.Error != nil {
		return c.JSON(http.StatusOK, "Could not find current user")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.OldPassword))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Old password is incorrect")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqBody.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	user.Password = string(hashedPassword)
	tx = h.db.Save(&user)
	if tx.Error != nil {
		return c.JSON(http.StatusBadRequest, "Could not save user")
	}
	return c.JSON(http.StatusCreated, user)
}
