package log

import "context"

type loggerContextKey struct{}

// ToContext save logger into context
func ToContext(ctx context.Context, logger Logger) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, loggerContextKey{}, logger)
}

// FromContext get logger from context or create new logger
func FromContext(ctx context.Context) Logger {
	if ctx == nil {
		return New()
	}
	if l, ok := ctx.Value(loggerContextKey{}).(Logger); ok && l != nil {
		return l
	}
	return New()
}
