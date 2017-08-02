package api

import (
	"github.com/kataras/iris"

	apiErrors "github.com/iris-contrib/parrot/parrot-api/errors"
	"github.com/iris-contrib/parrot/parrot-api/model"
	"github.com/iris-contrib/parrot/parrot-api/render"
)

// getUserProjects is an API endpoint for retrieving all projects that a user
// has access to.
func getUserProjects(ctx iris.Context) {
	id, err := getSubjectID(ctx)
	if err != nil {
		handleError(ctx, apiErrors.ErrBadRequest)
		return
	}

	projects, err := store.GetUserProjects(id)
	if err != nil {
		handleError(ctx, err)
		return
	}

	render.JSON(ctx, iris.StatusOK, projects)
}

// getProjectUsers is an API endpoint for retrieving all users with access to a project.
func getProjectUsers(ctx iris.Context) {
	projectID := ctx.Params().Get("projectID")
	if projectID == "" {
		handleError(ctx, apiErrors.ErrBadRequest)
		return
	}

	projectUsers, err := store.GetProjectUsers(projectID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	// Remove self user from slice
	id, err := getSubjectID(ctx)
	if err != nil {
		handleError(ctx, apiErrors.ErrBadRequest)
		return
	}

	result := make([]model.ProjectUser, 0)
	for _, pu := range projectUsers {
		if pu.UserID == id {
			continue
		}
		result = append(result, pu)
	}

	render.JSON(ctx, iris.StatusOK, result)
}

// assignProjectUser is an API endpoint for giving an already registered user
// rights to access a project.
func assignProjectUser(ctx iris.Context) {
	projectID := ctx.Params().Get("projectID")
	if projectID == "" {
		handleError(ctx, apiErrors.ErrBadRequest)
		return
	}

	// TODO: decode and validate only required fields. Whitelisting?
	var pu model.ProjectUser
	if err := ctx.ReadJSON(&pu); err != nil {
		handleError(ctx, apiErrors.ErrBadRequest)
		return
	}

	// Don't allow self editing
	id, err := getSubjectID(ctx)
	if err != nil {
		handleError(ctx, apiErrors.ErrBadRequest)
		return
	}
	if id == pu.UserID {
		handleError(ctx, apiErrors.ErrForbiden)
		return
	}

	// Validate that the url of the request matches the body data
	if projectID != pu.ProjectID {
		handleError(ctx, apiErrors.ErrForbiden)
		return
	}
	// If neither email nor user id is provided, there's nothing we can do
	if pu.Email == "" && pu.UserID == "" {
		handleError(ctx, apiErrors.ErrBadRequest)
		return
	}

	// If email is provided, but no user id, find the user by email
	// Otherwise we already have the id, and no need to fetch data before the grant operation
	if pu.Email != "" && pu.UserID == "" {
		user, err := store.GetUserByEmail(pu.Email)
		if err != nil {
			handleError(ctx, err)
			return
		}
		pu.UserID = user.ID
	}

	result, err := store.AssignProjectUser(pu)
	if err != nil {
		handleError(ctx, err)
		return
	}

	render.JSON(ctx, iris.StatusOK, result)
}

// updateProjectUserRole is an API endpoint for changing a user's role in a project.
func updateProjectUserRole(ctx iris.Context) {
	projectID := ctx.Params().Get("projectID")
	if projectID == "" {
		handleError(ctx, apiErrors.ErrBadRequest)
		return
	}
	userID := ctx.Params().Get("userID")
	if userID == "" {
		handleError(ctx, apiErrors.ErrBadRequest)
		return
	}

	// Get updated role
	data := struct {
		Role string `json:"role"`
	}{}
	if err := ctx.ReadJSON(&data); err != nil {
		handleError(ctx, apiErrors.ErrBadRequest)
		return
	}
	if !isRole(data.Role) {
		handleError(ctx, apiErrors.ErrBadRequest)
		return
	}

	pu := model.ProjectUser{UserID: userID, ProjectID: projectID, Role: data.Role}

	result, err := store.UpdateProjectUser(pu)
	if err != nil {
		handleError(ctx, err)
		return
	}

	render.JSON(ctx, iris.StatusOK, result)
}

// revokeProjectUser is an API endpoint for removing a user's role from a project.
func revokeProjectUser(ctx iris.Context) {
	projectID := ctx.Params().Get("projectID")
	if projectID == "" {
		handleError(ctx, apiErrors.ErrBadRequest)
		return
	}
	userID := ctx.Params().Get("userID")
	if userID == "" {
		handleError(ctx, apiErrors.ErrBadRequest)
		return
	}
	pu := model.ProjectUser{UserID: userID, ProjectID: projectID}

	err := store.RevokeProjectUser(pu)
	if err != nil {
		handleError(ctx, err)
		return
	}

	render.JSON(ctx, iris.StatusNoContent, nil)
}
