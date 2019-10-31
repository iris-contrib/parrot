package api

import (
	"strings"

	apiErrors "github.com/iris-contrib/parrot/parrot-api/errors"
	"github.com/iris-contrib/parrot/parrot-api/model"
	"github.com/iris-contrib/parrot/parrot-api/render"

	"github.com/kataras/iris/v12"
)

type projectKeyPayload struct {
	Key string `json:"key"`
}

type projectKeyUpdatePayload struct {
	OldKey string `json:"oldKey"`
	NewKey string `json:"newKey"`
}

// createProject is an API endpoint for creating new projects.
func createProject(ctx iris.Context) {
	project := model.Project{}
	errs := decodeAndValidate(ctx, &project)
	if errs != nil {
		render.Error(ctx, iris.StatusUnprocessableEntity, errs)
		return
	}
	userID, err := getSubjectID(ctx)
	if err != nil {
		handleError(ctx, apiErrors.ErrBadRequest)
		return
	}

	// TODO: use a transaction for this
	result, err := store.CreateProject(project)
	if err != nil {
		handleError(ctx, err)
		return
	}
	pu := model.ProjectUser{ProjectID: result.ID, UserID: userID, Role: ownerRole}
	_, err = store.AssignProjectUser(pu)
	if err != nil {
		handleError(ctx, err)
		return
	}

	render.JSON(ctx, iris.StatusCreated, result)
}

// updateProjectName is an API endpoint for updating the name of a project.
func updateProjectName(ctx iris.Context) {
	projectID := ctx.Params().Get("projectID")
	if projectID == "" {
		handleError(ctx, apiErrors.ErrBadRequest)
		return
	}

	project := model.Project{}
	errs := decodeAndValidate(ctx, &project)
	if errs != nil {
		render.Error(ctx, iris.StatusUnprocessableEntity, errs)
		return
	}

	result, err := store.UpdateProjectName(projectID, project.Name)
	if err != nil {
		handleError(ctx, err)
		return
	}

	render.JSON(ctx, iris.StatusOK, result)
}

// addProjectKey is an API endpoint for adding keys ('strings') to a project.
func addProjectKey(ctx iris.Context) {
	projectID := ctx.Params().Get("projectID")
	if projectID == "" {
		handleError(ctx, apiErrors.ErrBadRequest)
		return
	}

	var data = projectKeyPayload{}
	if err := ctx.ReadJSON(&data); err != nil {
		handleError(ctx, err)
		return
	}

	if data.Key == "" {
		handleError(ctx, apiErrors.ErrUnprocessable)
		return
	}

	data.Key = strings.Trim(data.Key, " ")

	result, err := store.AddProjectKey(projectID, data.Key)
	if err != nil {
		handleError(ctx, err)
		return
	}

	render.JSON(ctx, iris.StatusOK, result)
}

// updateProjectKey is an API endpoint for renaming keys ('strings') in a project.
func updateProjectKey(ctx iris.Context) {
	projectID := ctx.Params().Get("projectID")
	if projectID == "" {
		handleError(ctx, apiErrors.ErrBadRequest)
		return
	}

	var data = projectKeyUpdatePayload{}

	if err := ctx.ReadJSON(&data); err != nil {
		handleError(ctx, err)
		return
	}

	if data.OldKey == "" || data.NewKey == "" {
		handleError(ctx, apiErrors.ErrUnprocessable)
		return
	}

	data.NewKey = strings.Trim(data.NewKey, "")

	project, localesAffected, err := store.UpdateProjectKey(projectID, data.OldKey, data.NewKey)
	if err != nil {
		handleError(ctx, err)
		return
	}

	result := map[string]interface{}{
		"localesAffected": localesAffected,
		"project":         project,
	}

	render.JSON(ctx, iris.StatusOK, result)
}

// deleteProjectKey is an API endpoint for deleting keys ('strings') from a project.
func deleteProjectKey(ctx iris.Context) {
	projectID := ctx.Params().Get("projectID")
	if projectID == "" {
		handleError(ctx, apiErrors.ErrBadRequest)
		return
	}

	var data = projectKeyPayload{}
	if err := ctx.ReadJSON(&data); err != nil {
		handleError(ctx, err)
		return
	}

	if data.Key == "" {
		handleError(ctx, apiErrors.ErrUnprocessable)
		return
	}

	result, err := store.DeleteProjectKey(projectID, data.Key)
	if err != nil {
		handleError(ctx, err)
		return
	}

	render.JSON(ctx, iris.StatusOK, result)
}

// showProject is an API endpoint for retrieving a particular project.
func showProject(ctx iris.Context) {
	projectID := ctx.Params().Get("projectID")
	if projectID == "" {
		handleError(ctx, apiErrors.ErrBadRequest)
		return
	}

	project, err := store.GetProject(projectID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	render.JSON(ctx, iris.StatusOK, project)
}

// deleteProject is an API endpoint for deleting a particular project.
func deleteProject(ctx iris.Context) {
	projectID := ctx.Params().Get("projectID")
	if projectID == "" {
		handleError(ctx, apiErrors.ErrBadRequest)
		return
	}

	err := store.DeleteProject(projectID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	render.JSON(ctx, iris.StatusNoContent, nil)
}
