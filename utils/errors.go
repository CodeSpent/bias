package utils

import (
	"github.com/labstack/echo/v4"
)

func RespondWithError(c echo.Context, errMsg string, statusCode int) error {
	return c.JSON(statusCode, map[string]string{"error": errMsg})
}
