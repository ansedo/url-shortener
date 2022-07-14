package helpers

import (
	"context"
	"github.com/ansedo/url-shortener/internal/middlewares"
)

func GetUIDFromCtx(ctx context.Context) string {
	if uid := ctx.Value(middlewares.CookieCtxName); uid != nil {
		return uid.(string)
	}
	return ""
}
