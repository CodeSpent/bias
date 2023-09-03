package handlers

import (
	"github.com/labstack/echo/v4"
)

type BaseHandler struct{}

func (bh *BaseHandler) ErrorResponse(c echo.Context, status int, err error) error {
	return c.JSON(status, map[string]string{"error": err.Error()})
}

func (bh *BaseHandler) SuccessResponse(c echo.Context, status int, data interface{}) error {
	return c.JSON(status, data)
}
