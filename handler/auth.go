package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/DenisBytes/GoTemplHTMXLoginTemplate/db"
	"github.com/DenisBytes/GoTemplHTMXLoginTemplate/pkg/sb"
	"github.com/DenisBytes/GoTemplHTMXLoginTemplate/pkg/validate"
	"github.com/DenisBytes/GoTemplHTMXLoginTemplate/types"
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
	if len(accessToken) == 0 {
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

func HandleLoginPost(e echo.Context) error {
	params := supabase.UserCredentials{
		Email:    e.Request().FormValue("email"),
		Password: e.Request().FormValue("password"),
	}
	resp, err := sb.Client.Auth.SignIn(e.Request().Context(), params)
	if err != nil {
		slog.Error("Login error", "err", err)
		return auth.LoginForm(params, auth.LoginErrors{
			InvalidCredentials: "Invalid credentials",
		}).Render(e.Request().Context(), e.Response().Writer)
	}
	if err := setAuthSession(e.Response().Writer, e.Request(), resp.AccessToken); err != nil {
		return err
	}
	redirectURL := "/"
	toParam := e.QueryParam("to")
	if toParam != "" {
		redirectURL = toParam
	}
	return hxRedirect(e.Response().Writer, e.Request(), redirectURL)
}

func HandleLoginWIthGithub(e echo.Context) error {
	resp, err := sb.Client.Auth.SignInWithProvider(supabase.ProviderSignInOptions{
		Provider:   "github",
		RedirectTo: "http://localhost:3000/auth/callback",
	})
	if err != nil {
		return err
	}
	http.Redirect(e.Response().Writer, e.Request(), resp.URL, http.StatusSeeOther)
	return nil
}

func HandleAccountSetupIndex(e echo.Context) error {
	return auth.AccountSetup().Render(e.Request().Context(), e.Response().Writer)
}

func HandleAccountSetupPost(e echo.Context) error {
	params := auth.AccountSetupParams{
		Username: e.Request().FormValue("username"),
	}
	var errors auth.AccountSetupErrors
	ok := validate.New(&params, validate.Fields{
		"Username": validate.Rules(validate.Min(2), validate.Max(50)),
	}).Validate(&errors)
	if !ok {
		return auth.AccountSetupForm(params, errors).Render(e.Request().Context(), e.Response().Writer)
	}
	user := GetAuthenticatedUser(e)
	account := types.Account{

		UserID:   user.ID,
		Username: params.Username,
	}
	if err := db.CreateAccount(&account); err != nil {
		return err
	}
	return hxRedirect(e.Response().Writer, e.Request(), "/")
}

func HandleResetPasswordIndex(e echo.Context) error {
	return auth.ResetPassword().Render(e.Request().Context(), e.Response().Writer)
}

func HandleResetPasswordPost(e echo.Context) error {
	user := GetAuthenticatedUser(e)
	params := map[string]any{
		"email":      user.Email,
		"redirectTo": "http://localhost:3000/auth/reset-password",
	}
	b, err := json.Marshal(params)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s", sb.BaseAuthURL), bytes.NewReader(b))
	req.Header.Set("apikey", os.Getenv("SUPABASE_SECRET"))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("supabase password recovery responded with a non 200status code: %d", resp.StatusCode)
	}

	return auth.ResetPasswordInitiated(user.Email).Render(e.Request().Context(), e.Response().Writer)
}

func HandleResetPasswordUpdate(e echo.Context) error {
	user := GetAuthenticatedUser(e)
	params := map[string]any{
		"password": e.Request().FormValue("password"),
	}
	resp, err := sb.Client.Auth.UpdateUser(e.Request().Context(), user.AccessToken, params)
	if err != nil {
		return auth.ResetPasswordForm(auth.ResetPasswordErrors{NewPassword: "Please enter a valid password"}).Render(e.Request().Context(), e.Response().Writer)
	}
	_ = resp
	return hxRedirect(e.Response().Writer, e.Request(), "/")
}
