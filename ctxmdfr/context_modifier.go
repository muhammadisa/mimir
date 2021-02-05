package ctxmdfr

import "context"

func ContextWithTZ(parentCtx context.Context, tzValue string) context.Context {
	return context.WithValue(parentCtx, "tz", tzValue)
}

func GetTZFromCtx(ctx context.Context) string {
	return ctx.Value("tz").(string)
}
