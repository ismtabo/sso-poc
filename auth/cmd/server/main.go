package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/ismtabo/sso-poc/auth/pkg/cfg"
	"github.com/ismtabo/sso-poc/auth/pkg/controller"
	"github.com/ismtabo/sso-poc/auth/pkg/middleware"
	"github.com/ismtabo/sso-poc/auth/pkg/repository"
	"github.com/ismtabo/sso-poc/auth/pkg/service"
	"github.com/rs/zerolog"
	"github.com/sonyarouje/simdb"
)

func main() {
	log := getLogger()
	var config Config
	if err := cfg.Load("config.yml", &config); err != nil {
		log.Fatal().Msgf("Error loading configuration. %s", err)
	}
	log.Debug().Msgf("Configuration loaded: %+v", config)
	if err := configLogger(&config); err != nil {
		log.Fatal().Msgf("Error configuring the logger. %s", err)
	}
	driver := mustCreateSimdbDriver(&config, log)
	userRepo := repository.NewSimdbUserRepository(driver)
	fileRepo := repository.NewPageRepository(config.Web.Pages)
	authSvc := service.NewAuthService(userRepo)
	store := sessions.NewCookieStore([]byte(config.Web.Session.Key))
	sessionSvc := service.NewCookieSessionService(store, config.Web.Session.Cookie)
	authCtrl := controller.NewAuthController(authSvc, sessionSvc)
	pagesCtrl := controller.NewPagesController(sessionSvc, fileRepo)
	http.Handle("/", middleware.WithSession(sessionSvc, http.HandlerFunc(pagesCtrl.Home)))
	http.Handle("/login", middleware.WithSession(sessionSvc, http.HandlerFunc(pagesCtrl.Login)))
	http.HandleFunc("/auth/login", authCtrl.Login)
	http.HandleFunc("/auth/logout", authCtrl.Logout)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(config.Web.Static))))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal().Msg("Error starting server")
	}
}

func getLogger() *zerolog.Logger {
	zerolog.TimeFieldFormat = "2006-01-02T15:04:05.000Z07:00"
	zerolog.TimestampFieldName = "time"
	zerolog.LevelFieldName = "lvl"
	zerolog.MessageFieldName = "msg"
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return &log
}

func configLogger(config *Config) error {
	lvl, err := zerolog.ParseLevel(strings.ToLower(config.Log.Level))
	if err != nil {
		return err
	}
	zerolog.SetGlobalLevel(lvl)
	return nil
}

func mustCreateSimdbDriver(config *Config, log *zerolog.Logger) *simdb.Driver {
	driver, err := simdb.New(config.Database.File)
	if err != nil {
		log.Fatal().Msgf("Error creating simdb driver. %s", err)
	}
	return driver
}
