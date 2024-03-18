package handler

import (
	"github.com/DenisBytes/GoTemplHTMXLoginTemplate/db"
	"github.com/DenisBytes/GoTemplHTMXLoginTemplate/pkg/validate"
	"github.com/DenisBytes/GoTemplHTMXLoginTemplate/view/settings"
	"github.com/labstack/echo/v4"
)

func HandleSettingsIndex(e echo.Context) error {
	user := GetAuthenticatedUser(e)
	return settings.Index(user).Render(e.Request().Context(), e.Response().Writer)
}

func HandleSettingsUsernameUpdate(e echo.Context) error {
	params := settings.ProfileParams{
		Username: e.Request().FormValue("username"),
	}
	errors := settings.ProfileErrors{}
	ok := validate.New(&params, validate.Fields{
		"Username": validate.Rules(validate.Min(3), validate.Max(40)),
	}).Validate(&errors)
	if !ok {
		return settings.ProfileForm(params, errors).Render(e.Request().Context(), e.Response().Writer)
	}

	user := GetAuthenticatedUser(e)
	user.Account.Username = params.Username
	if err := db.UpdateAccount(&user.Account); err != nil {
		return err
	}
	params.Success = true
	return settings.ProfileForm(params, settings.ProfileErrors{}).Render(e.Request().Context(), e.Response().Writer)
}
