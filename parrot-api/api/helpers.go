package api

import (
	"github.com/kataras/iris/v12"
	"github.com/iris-contrib/parrot/parrot-api/render"
)

var (
	validContentTypes = []string{
		"application/json",
		"application/json; charset=utf-8"}
)

// isValidContentType returns true if the provided content type
// is an allowed one.
func isValidContentType(ct string) bool {
	if ct == "" {
		return false
	}
	for _, v := range validContentTypes {
		if ct == v {
			return true
		}
	}
	return false
}

// ping is an API endpoint for checking if the API is up.
func ping(ctx iris.Context) {
	render.JSON(ctx, 200, map[string]interface{}{
		"message": "Parrot says hello.",
	})
}
