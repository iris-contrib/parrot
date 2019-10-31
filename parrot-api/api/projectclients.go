package api

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/kataras/iris/v12"

	apiErrors "github.com/iris-contrib/parrot/parrot-api/errors"
	"github.com/iris-contrib/parrot/parrot-api/model"
	"github.com/iris-contrib/parrot/parrot-api/render"
)

var (
	clientSecretBytes = 32
)

// getProjectClients is an API endpoint for retrieving all clients ('applications') for a project.
func getProjectClients(ctx iris.Context) {
	projectID := ctx.Params().Get("projectID")
	if projectID == "" {
		handleError(ctx, apiErrors.ErrBadRequest)
		return
	}

	result, err := store.GetProjectClients(projectID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	render.JSON(ctx, iris.StatusOK, result)
}

// getProjectClient is an API endpoint for retrieving a project client.
func getProjectClient(ctx iris.Context) {
	clientID := ctx.Params().Get("clientID")
	if clientID == "" {
		handleError(ctx, apiErrors.ErrBadRequest)
		return
	}
	projectID := ctx.Params().Get("projectID")
	if projectID == "" {
		handleError(ctx, apiErrors.ErrBadRequest)
		return
	}

	result, err := store.GetProjectClient(projectID, clientID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	render.JSON(ctx, iris.StatusOK, result)
}

// deleteProjectClient is an API endpoint for deleting a project client.
func deleteProjectClient(ctx iris.Context) {
	projectID := ctx.Params().Get("projectID")
	if projectID == "" {
		handleError(ctx, apiErrors.ErrBadRequest)
		return
	}
	clientID := ctx.Params().Get("clientID")
	if clientID == "" {
		handleError(ctx, apiErrors.ErrBadRequest)
		return
	}

	err := store.DeleteProjectClient(projectID, clientID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	render.JSON(ctx, iris.StatusNoContent, nil)
}

// createProjectClient is an API endpoint for registering a new project client.
func createProjectClient(ctx iris.Context) {
	projectID := ctx.Params().Get("projectID")
	if projectID == "" {
		handleError(ctx, apiErrors.ErrBadRequest)
		return
	}

	pc := model.ProjectClient{}
	errs := decodeAndValidate(ctx, &pc)
	if errs != nil {
		render.Error(ctx, iris.StatusUnprocessableEntity, errs)
		return
	}
	secret, err := generateClientSecret(clientSecretBytes)
	if err != nil {
		handleError(ctx, apiErrors.ErrInternal)
		return
	}
	pc.Secret = secret
	pc.ProjectID = projectID

	result, err := store.CreateProjectClient(pc)
	if err != nil {
		handleError(ctx, err)
		return
	}

	render.JSON(ctx, iris.StatusCreated, result)
}

// updateProjectClientName is an API endpoint for updating a project client's name.
func updateProjectClientName(ctx iris.Context) {
	projectID := ctx.Params().Get("projectID")
	if projectID == "" {
		handleError(ctx, apiErrors.ErrBadRequest)
		return
	}
	clientID := ctx.Params().Get("clientID")
	if clientID == "" {
		handleError(ctx, apiErrors.ErrBadRequest)
		return
	}

	pc := model.ProjectClient{}
	errs := decodeAndValidate(ctx, &pc)
	if errs != nil {
		render.Error(ctx, iris.StatusUnprocessableEntity, errs)
		return
	}
	pc.ProjectID = projectID
	pc.ClientID = clientID

	result, err := store.UpdateProjectClientName(pc)
	if err != nil {
		handleError(ctx, err)
		return
	}

	render.JSON(ctx, iris.StatusOK, result)
}

// resetProjectClientSecret is an API endpoint for regenerating a project client's secret.
func resetProjectClientSecret(ctx iris.Context) {
	projectID := ctx.Params().Get("projectID")
	if projectID == "" {
		handleError(ctx, apiErrors.ErrBadRequest)
		return
	}
	clientID := ctx.Params().Get("clientID")
	if clientID == "" {
		handleError(ctx, apiErrors.ErrBadRequest)
		return
	}
	secret, err := generateClientSecret(clientSecretBytes)
	if err != nil {
		handleError(ctx, apiErrors.ErrInternal)
		return
	}

	pc := model.ProjectClient{
		ClientID:  clientID,
		ProjectID: projectID,
		Secret:    secret}

	result, err := store.UpdateProjectClientSecret(pc)
	if err != nil {
		handleError(ctx, err)
		return
	}

	render.JSON(ctx, iris.StatusOK, result)
}

// generateClientSecret generates a cryptographically secure pseudorandom string.
func generateClientSecret(bytes int) (string, error) {
	b := make([]byte, bytes)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
