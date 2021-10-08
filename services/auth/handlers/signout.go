package handlers

import (
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func (h *Handler) SignOut(c echo.Context) error {
	sess, err := session.Get("pay9.sess", c)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	sess.Options.MaxAge = -1
	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, "{}")
}
