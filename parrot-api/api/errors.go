package api

import (
	"github.com/kataras/iris/v12"

	datastoreErrors "github.com/iris-contrib/parrot/parrot-api/datastore/errors"
	apiErrors "github.com/iris-contrib/parrot/parrot-api/errors"
	"github.com/iris-contrib/parrot/parrot-api/render"
)

// handleError writes an error response.
// If the error is not a known API error, it will try to
// cast it or simply write an Internal error.
func handleError(ctx iris.Context, err error) {
	// Try to match store error
	var outErr *apiErrors.Error
	// If cast is successful, done, we got our error
	if castedErr, ok := err.(*apiErrors.Error); ok {
		outErr = castedErr
	} else {
		// Check if it is a datastore error
		switch err {
		case datastoreErrors.ErrNotFound:
			outErr = apiErrors.ErrNotFound
		case datastoreErrors.ErrAlreadyExists:
			outErr = apiErrors.ErrAlreadyExists
		default:
			ctx.Application().Logger().Errorf("%v", err)
			outErr = apiErrors.ErrInternal

		}
	}

	render.Error(ctx, outErr.Status, outErr)
}
