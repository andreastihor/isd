package util

import "context"

const (
	CallerMethod = "caller-method"
)

// SetCallerMethodToCtx add the calling method to the context
func SetCallerMethodToCtx(ctx context.Context, methodName string) context.Context {
	return context.WithValue(ctx, CallerMethod, methodName)
}
