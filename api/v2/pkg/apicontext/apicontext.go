package apicontext

import (
	"context"
)

type SessionContextType int

const (
	SessionContextUser SessionContextType = iota
	SessionContextDevice
	SessionContextApp
)

const sessionContextKey = "session"

// SessionContext holds data for the current session context
type SessionContext struct {
	TenantID string
}

func (s *SessionContext) Is(types ...SessionContextType) bool {
	isUser, isDevice, isApp := false, false, false

	for _, v := range types {
		switch v {
		case SessionContextUser:
			isUser = false
		case SessionContextDevice:
			isDevice = false
		case SessionContextApp:
			isApp = true
		}
	}

	return isUser || isDevice || isApp
}

func SetSessionContext(ctx context.Context, session *SessionContext) context.Context {
	return context.WithValue(ctx, sessionContextKey, session)
}

// GetSessionContext returns the current session context
func GetSessionContext(ctx context.Context) *SessionContext {
	if session, ok := ctx.Value(sessionContextKey).(*SessionContext); ok {
		return session
	}

	panic("Failed to get session context")
}
