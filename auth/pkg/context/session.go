package context

import "context"

type ContextKey string

const sessionCtxKey = ContextKey("session")

type SessionContext struct {
	Authenticated bool
}

func WithSessionContext(ctx context.Context) context.Context {
	if session := ctx.Value(sessionCtxKey); session == nil {
		return context.WithValue(ctx, sessionCtxKey, &SessionContext{})
	}
	return ctx
}

func GetContextSession(ctx context.Context) *SessionContext {
	if session := ctx.Value(sessionCtxKey); session != nil {
		return session.(*SessionContext)
	}
	return nil
}
