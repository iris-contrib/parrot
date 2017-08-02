package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/kataras/iris"

	"github.com/iris-contrib/parrot/parrot-api/auth"
	apiErrors "github.com/iris-contrib/parrot/parrot-api/errors"
)

// subjectType is an internal identifier to know if the requesting entity
// is a project user or an application.
type subjectType string

const (
	userSubject   = "user"
	clientSubject = "client"
)

// tokenMiddleware guards against request without a valid token.
// Adds subject ID and subject type values to request context.
func tokenMiddleware(tp auth.TokenProvider) iris.Handler {
	return func(ctx iris.Context) {
		tokenString, err := getTokenString(ctx.Request())
		if err != nil {
			handleError(ctx, apiErrors.ErrUnauthorized)
			return
		}

		claims, err := tp.ParseAndVerifyToken(tokenString)
		if err != nil {
			handleError(ctx, apiErrors.ErrUnauthorized)
			return
		}

		subID := claims["sub"]
		if subID == nil || subID == "" {
			handleError(ctx, apiErrors.ErrInternal)
			return
		}

		subType := claims["subType"]
		if subType == nil || subType == "" {
			handleError(ctx, apiErrors.ErrInternal)
			return
		}

		ctx.Values().Set("subjectID", subID)
		ctx.Values().Set("subjectType", subType)
		ctx.Next()
	}
}

// getTokenString extracts the encoded token from HTTP Authorization Headers.
func getTokenString(r *http.Request) (string, error) {
	token := r.Header.Get("Authorization")
	if token == "" {
		return "", fmt.Errorf("no auth header")
	}

	if len(token) > 6 && strings.ToUpper(token[0:7]) == "BEARER " {
		token = token[7:]
	}

	return token, nil
}

// getSubjectID extract subject ID from context.
func getSubjectID(ctx iris.Context) (string, error) {
	v := ctx.Values().Get("subjectID")
	if v == nil {
		return "", apiErrors.ErrBadRequest
	}
	id, ok := v.(string)
	if id == "" || !ok {
		return "", apiErrors.ErrInternal
	}
	return id, nil
}

// getSubjectType extract user type from context.
func getSubjectType(ctx iris.Context) (subjectType, error) {
	subType := ctx.Values().Get("subjectType")
	if subType == nil {
		return "", apiErrors.ErrBadRequest
	}

	casted, ok := subType.(string)
	if !ok || casted == "" {
		return "", apiErrors.ErrBadRequest
	}

	return subjectType(casted), nil
}
