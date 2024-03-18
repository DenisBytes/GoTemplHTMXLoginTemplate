package handler

import (
	"net/http"

	"github.com/DenisBytes/GoTemplHTMXLoginTemplate/types"
	"github.com/labstack/echo/v4"
)

func hxRedirect(w http.ResponseWriter, r *http.Request, path string) error {
	if len(r.Header.Get("HX-Request")) > 0 {
		w.Header().Set("HX-Redirect", path)
		w.WriteHeader(http.StatusSeeOther)
		return nil
	}
	http.Redirect(w, r, path, http.StatusSeeOther)
	return nil
}

func GetAuthenticatedUser(e echo.Context) types.User {
	user, ok := e.Request().Context().Value("user").(types.User)
	if !ok {
		return types.User{}
	}
	return user
}
