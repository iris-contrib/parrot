// Package auth handles the creation of an Auth Provider and its routes.
package auth

import (
	"github.com/kataras/iris"
)

// NewRouter creates and configures all routes for the parameter authentication provider.
func NewRouter(ds AuthStore, tp TokenProvider) iris.Configurator {
	return func(app *iris.Application) {
		app.Post("/api/v1/auth/token", IssueToken(tp, ds))
	}
}
