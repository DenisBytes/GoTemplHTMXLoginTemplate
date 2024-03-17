package handler

import (
	"github.com/DenisBytes/GoTemplHTMXLoginTemplate/view/home"
	"github.com/labstack/echo/v4"
)

func HandleHomeIndex(ctx echo.Context) error {
	return home.Index().Render(ctx.Request().Context(), ctx.Response().Writer)
}
