package main

import (
	"embed"
	"log/slog"
	"os"

	"github.com/DenisBytes/GoTemplHTMXLoginTemplate/db"
	"github.com/DenisBytes/GoTemplHTMXLoginTemplate/handler"
	"github.com/DenisBytes/GoTemplHTMXLoginTemplate/pkg/sb"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:embed public/*
var FS embed.FS

func main() {
	if err := initEverything(); err != nil {
		slog.Error("Init err: ", err)
	}

	router := echo.New()

	router.Use(middleware.Logger())
	router.Use(middleware.Recover())
	router.Use(handler.WithUser)

	// Handler for static files
	router.StaticFS("/*", FS)

	router.GET("/", handler.HandleHomeIndex)
	router.GET("/login", handler.HandleLoginIndex)
	router.GET("/signup", handler.HandleSignUpIndex)
	router.POST("/signup", handler.HandleSignUpPost)
	router.GET("/auth/callback", handler.HandleAuthCallback)
	router.POST("/logout", handler.HandleLogoutPost)
	router.POST("/login", handler.HandleLoginPost)
	router.GET("/login/provider/github", handler.HandleLoginWIthGithub)

	router.GET("/account/setup", handler.HandleAccountSetupIndex, handler.WithAuth)
	router.POST("/account/setup", handler.HandleAccountSetupPost, handler.WithAuth)

	router.GET("/settings", handler.HandleSettingsIndex, handler.WithAccountSetup, handler.WithAuth)
	router.PUT("/settings/account/profile", handler.HandleSettingsUsernameUpdate, handler.WithAccountSetup, handler.WithAuth)

	router.GET("/auth/reset-password", handler.HandleResetPasswordIndex, handler.WithAccountSetup, handler.WithAuth)
	router.POST("/auth/reset-password", handler.HandleResetPasswordPost, handler.WithAccountSetup, handler.WithAuth)
	router.PUT("/auth/reset-password", handler.HandleResetPasswordUpdate, handler.WithAccountSetup, handler.WithAuth)

	port := os.Getenv("HTTP_LISTEN_ADDR")

	if err := router.Start(port); err != nil {
		slog.Error("Echo run and serve error:", err)
	}
}

func initEverything() error {
	if err := godotenv.Load(); err != nil {
		slog.Error("Godotenv error:", err)
		return err
	}
	if err := db.Init(); err != nil {
		return err
	}
	return sb.Init()
}
