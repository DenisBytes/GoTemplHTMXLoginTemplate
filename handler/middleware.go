package handler

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/DenisBytes/GoTemplHTMXLoginTemplate/db"
	"github.com/DenisBytes/GoTemplHTMXLoginTemplate/pkg/sb"
	"github.com/DenisBytes/GoTemplHTMXLoginTemplate/types"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

func WithUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(e echo.Context) error {
		if strings.Contains(e.Request().URL.Path, "/public") {
			return next(e)
		}
		store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
		session, err := store.Get(e.Request(), "user")
		if err != nil {
			return next(e)
		}
		accessToken := session.Values["accessToken"]
		if accessToken == nil {
			return next(e)
		}
		resp, err := sb.Client.Auth.User(e.Request().Context(), accessToken.(string))
		if err != nil {
			return next(e)
		}
		user := types.User{
			ID:          uuid.MustParse(resp.ID),
			Email:       resp.Email,
			LoggedIn:    true,
			AccessToken: accessToken.(string),
		}
		ctx := context.WithValue(e.Request().Context(), "user", user)
		e.SetRequest(e.Request().WithContext(ctx))
		return next(e)
	}
}

func WithAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(e echo.Context) error {
		if strings.Contains(e.Request().URL.Path, "/public") {
			return next(e)
		}
		user := GetAuthenticatedUser(e)
		if !user.LoggedIn {
			path := e.Request().URL.Path
			http.Redirect(e.Response().Writer, e.Request(), "/login?to="+path, http.StatusSeeOther)
			return nil
		}
		return next(e)
	}
}

func WithAccountSetup(next echo.HandlerFunc) echo.HandlerFunc {
	return func(e echo.Context) error {
		user := GetAuthenticatedUser(e)
		account, err := db.GetAccountByUserID(user.ID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Redirect(e.Response().Writer, e.Request(), "/account/setup", http.StatusSeeOther)
				return nil
			}
			return next(e)
		}
		user.Account = account

		ctx := context.WithValue(e.Request().Context(), "user", user)
		e.SetRequest(e.Request().WithContext(ctx))
		return next(e)
	}
}
