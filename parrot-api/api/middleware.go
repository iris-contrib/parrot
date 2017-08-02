package api

import (
	"github.com/kataras/iris"

	apiErrors "github.com/iris-contrib/parrot/parrot-api/errors"
)

// enforceContentTypeJSON only allows requests that have the
// Content-Type header set to a valid JSON mime type, unless
// the body is empty (useful for 'verb' or 'action' requests).
func enforceContentTypeJSON(ctx iris.Context) {
	switch ctx.Method() {
	case "POST", "PUT", "PATCH":
		ct := ctx.GetHeader("Content-Type")
		if !isValidContentType(ct) && ctx.Request().ContentLength > 0 {
			handleError(ctx, apiErrors.ErrUnsupportedMediaType)
			return
		}
	}

	ctx.Next()
}
