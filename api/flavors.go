package api

import (
	"net/http"

	"github.com/labstack/echo"
)

func getServiceFlavors(c echo.Context) error {
	return c.JSON(http.StatusOK, []struct {
		name        string
		description string
	}{})
}

func getInstanceFlavors(c echo.Context) error {
	return getServiceFlavors(c)
}
