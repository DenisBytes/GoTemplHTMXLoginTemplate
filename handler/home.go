package handler

import (
	"fmt"

	"github.com/DenisBytes/GoTemplHTMXLoginTemplate/view/home"
	"github.com/labstack/echo/v4"
)

func HandleHomeIndex(e echo.Context) error {
	user := GetAuthenticatedUser(e)
	fmt.Printf("%+v\n", user)
	return home.Index().Render(e.Request().Context(), e.Response().Writer)
}
