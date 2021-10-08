package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mongmx/simple-microservice/services/auth/models"
)

func (h *Handler) CurrentUser(c echo.Context) error {
	currentUser, ok := c.Get("currentUser").(*JwtCustomClaims)
	if !ok {
		return c.JSON(http.StatusOK, "Could not find current user")
	}
	var user models.User
	tx := h.db.Where("ID = ?", currentUser.ID).First(&user)
	if tx.Error != nil {
		return c.JSON(http.StatusOK, "Could not find current user")
	}
	return c.JSON(http.StatusOK, map[string]models.User{"currentUser": user})
}
