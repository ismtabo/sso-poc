package main

import (
	"context"
	"log"

	"github.com/ismtabo/sso-poc/auth/pkg/service"
	"github.com/rs/zerolog"
)

type contextKey string

const (
	appCtxKey = contextKey("app")
)

type appContext struct {
	log     *zerolog.Logger
	authSvc service.AuthService
}

func withAppCtx(ctx context.Context, log *zerolog.Logger, authSvc service.AuthService) context.Context {
	if app := ctx.Value(appCtxKey); app == nil {
		return context.WithValue(ctx, appCtxKey, &appContext{log: log, authSvc: authSvc})
	}
	return ctx
}

func getAppCtx(ctx context.Context) *appContext {
	if app := ctx.Value(appCtxKey); app != nil {
		return app.(*appContext)
	}
	log.Fatal("missing app context. may forgive to initialize using withAppCtx before reading")
	return nil
}
