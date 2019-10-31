package api

import (
	"github.com/kataras/iris/v12"

	apiErrors "github.com/iris-contrib/parrot/parrot-api/errors"
	"github.com/iris-contrib/parrot/parrot-api/model"
	"github.com/iris-contrib/parrot/parrot-api/render"
)

// createLocale is an API endpoint for creating a new project locale.
func createLocale(ctx iris.Context) {
	projectID := ctx.Params().Get("projectID")
	if projectID == "" {
		handleError(ctx, apiErrors.ErrBadRequest)
		return
	}

	loc := model.Locale{}
	errs := decodeAndValidate(ctx, &loc)
	if errs != nil {
		render.Error(ctx, iris.StatusUnprocessableEntity, errs)
		return
	}
	loc.ProjectID = projectID

	proj, err := store.GetProject(projectID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	loc.SyncKeys(proj.Keys)

	result, err := store.CreateLocale(loc)
	if err != nil {
		handleError(ctx, err)
		return
	}

	render.JSON(ctx, iris.StatusCreated, result)
}

// showLocale is an API endpoint for retrieving a project locale by ident.
func showLocale(ctx iris.Context) {
	projectID := ctx.Params().Get("projectID")
	if projectID == "" {
		handleError(ctx, apiErrors.ErrBadRequest)
		return
	}
	ident := ctx.Params().Get("localeIdent")
	if ident == "" {
		handleError(ctx, apiErrors.ErrBadRequest)
		return
	}

	loc, err := store.GetProjectLocaleByIdent(projectID, ident)
	if err != nil {
		handleError(ctx, err)
		return
	}

	proj, err := store.GetProject(projectID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	loc.SyncKeys(proj.Keys)

	render.JSON(ctx, iris.StatusOK, loc)
}

// findLocales is an API endpoint for retrieving project locales and filtering by ident.
func findLocales(ctx iris.Context) {
	projectID := ctx.Params().Get("projectID")
	if projectID == "" {
		handleError(ctx, apiErrors.ErrBadRequest)
		return

	}
	localeIdents := ctx.Request().URL.Query()["ident"]

	locs, err := store.GetProjectLocales(projectID, localeIdents...)
	if err != nil {
		handleError(ctx, err)
		return
	}

	project, err := store.GetProject(projectID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	for i := range locs {
		locs[i].SyncKeys(project.Keys)
	}

	render.JSON(ctx, iris.StatusOK, locs)
}

// updateLocalePairs is an API endpoint for updating a locale's key value pairs.
func updateLocalePairs(ctx iris.Context) {
	ident := ctx.Params().Get("localeIdent")
	if ident == "" {
		handleError(ctx, apiErrors.ErrBadRequest)
		return
	}
	projectID := ctx.Params().Get("projectID")
	if projectID == "" {
		handleError(ctx, apiErrors.ErrBadRequest)
		return
	}

	loc := &model.Locale{}

	if err := ctx.ReadJSON(&loc.Pairs); err != nil {
		handleError(ctx, apiErrors.ErrUnprocessable)
		return
	}

	project, err := store.GetProject(projectID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	loc.SyncKeys(project.Keys)

	result, err := store.UpdateLocalePairs(projectID, ident, loc.Pairs)
	if err != nil {
		handleError(ctx, err)
		return
	}

	render.JSON(ctx, iris.StatusOK, result)
}

// deleteLocale is an API endpoint for deleting a project's locale.
func deleteLocale(ctx iris.Context) {
	ident := ctx.Params().Get("localeIdent")
	if ident == "" {
		handleError(ctx, apiErrors.ErrBadRequest)
		return
	}
	projectID := ctx.Params().Get("projectID")
	if projectID == "" {
		handleError(ctx, apiErrors.ErrBadRequest)
		return
	}

	err := store.DeleteLocale(projectID, ident)
	if err != nil {
		handleError(ctx, err)
		return
	}

	render.JSON(ctx, iris.StatusNoContent, nil)
}
