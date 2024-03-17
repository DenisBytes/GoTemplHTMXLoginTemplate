package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/DenisBytes/GoTemplHTMXLoginTemplate/pkg/sb"
	"github.com/DenisBytes/GoTemplHTMXLoginTemplate/pkg/validate"
	"github.com/DenisBytes/GoTemplHTMXLoginTemplate/view/auth"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/nedpals/supabase-go"
)

func HandleLoginIndex(e echo.Context) error {
	return auth.Login().Render(e.Request().Context(), e.Response().Writer)
}

func HandleSignUpIndex(e echo.Context) error {
	return auth.SignUp().Render(e.Request().Context(), e.Response().Writer)
}

func HandleSignUpPost(e echo.Context) error {
	params := auth.SignUpParams{
		Email:           e.Request().FormValue("email"),
		Password:        e.Request().FormValue("password"),
		ConfirmPassword: e.Request().FormValue("confirmPassword"),
	}

	var errors = auth.SignUpErrors{}
	if ok := validate.New(&params, validate.Fields{
		"email":           validate.Rules(validate.Email),
		"password":        validate.Rules(validate.Password),
		"confirmPassword": validate.Rules(validate.Equal(params.Password), validate.Message("Password don't match")),
	}).Validate(&errors); !ok {
		fmt.Println("here")
		return auth.SignUpForm(params, errors).Render(e.Request().Context(), e.Response().Writer)
	}
	user, err := sb.Client.Auth.SignUp(e.Request().Context(), supabase.UserCredentials{
		Email:    params.Email,
		Password: params.Password,
	})
	if err != nil {
		return err
	}
	return auth.SignUpSuccess(user.Email).Render(e.Request().Context(), e.Response().Writer)
}

func HandleAuthCallback(e echo.Context) error {
	accessToken := e.Request().URL.Query().Get("access_token")
	if len(accessToken) > 0 {
		return auth.CallbackScript().Render(e.Request().Context(), e.Response().Writer)
	}
	if err := setAuthSession(e.Response().Writer, e.Request(), accessToken); err != nil {
		return err
	}
	http.Redirect(e.Response().Writer, e.Request(), "/", http.StatusSeeOther)
	return nil
}

func setAuthSession(w http.ResponseWriter, r *http.Request, accessToken string) error {
	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	session, _ := store.Get(r, "user")
	session.Values["accessToken"] = accessToken
	return session.Save(r, w)
}

func HandleLogoutPost(e echo.Context) error {
	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	session, _ := store.Get(e.Request(), "user")
	session.Values["accessToken"] = ""
	session.Save(e.Request(), e.Response().Writer)
	http.Redirect(e.Response().Writer, e.Request(), "/login", http.StatusSeeOther)
	return nil
}
