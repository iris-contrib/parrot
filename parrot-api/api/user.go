package api

import (
	"errors"

	"github.com/kataras/iris"

	apiErrors "github.com/iris-contrib/parrot/parrot-api/errors"
	"github.com/iris-contrib/parrot/parrot-api/model"
	"github.com/iris-contrib/parrot/parrot-api/render"
	"golang.org/x/crypto/bcrypt"
)

type userSelfPayload struct {
	*model.User
	ProjectRoles  projectRoles  `json:"projectRoles,omitempty"`
	ProjectGrants projectGrants `json:"projectGrants,omitempty"`
}

type projectGrants map[string][]RoleGrant

type projectRoles map[string]string

// getUserSelf is an API endpoint for getting the requesting user's details.
func getUserSelf(ctx iris.Context) {
	id, err := getSubjectID(ctx)
	if err != nil {
		handleError(ctx, apiErrors.ErrBadRequest)
		return
	}

	user, err := store.GetUserByID(id)
	if err != nil {
		handleError(ctx, err)
		return
	}
	// Hide password
	user.Password = ""

	payload := userSelfPayload{user, nil, nil}

	include := ctx.Request().URL.Query().Get("include")
	if include != "" {
		switch include {
		case "projectRoles":
			projectUsers, err := store.GetUserProjectRoles(user.ID)
			if err != nil {
				handleError(ctx, err)
				return
			}

			result := make(projectRoles)
			for _, pu := range projectUsers {
				result[pu.ProjectID] = pu.Role
			}

			payload.ProjectRoles = result

		case "projectGrants":
			projectUsers, err := store.GetUserProjectRoles(user.ID)
			if err != nil {
				handleError(ctx, err)
				return
			}

			grants := make(projectGrants)
			for _, pu := range projectUsers {
				role := Role(pu.Role)
				grants[pu.ProjectID] = permissions[role]
			}
			payload.ProjectGrants = grants
		}
	}

	render.JSON(ctx, iris.StatusOK, payload)
}

// createUser is an API endpoint for registering new users.
func createUser(ctx iris.Context) {
	user := model.User{}
	errs := decodeAndValidate(ctx, &user)
	if errs != nil {
		render.Error(ctx, iris.StatusUnprocessableEntity, errs)
		return
	}

	existingUser, err := store.GetUserByEmail(user.Email)
	if err == nil && existingUser.Email == user.Email {
		handleError(ctx, apiErrors.ErrAlreadyExists)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		handleError(ctx, err)
		return
	}

	user.Password = string(hashed)

	result, err := store.CreateUser(user)
	if err != nil {
		handleError(ctx, err)
		return
	}

	// Hide password
	result.Password = ""
	render.JSON(ctx, iris.StatusCreated, result)
}

// updateUserPassword is an API endpoint for changing a user's password.
func updateUserPassword(ctx iris.Context) {
	payload := updatePasswordPayload{}
	err := decodePayloadAndValidate(ctx, &payload)
	if err != nil {
		handleError(ctx, apiErrors.ErrUnprocessable)
		return
	}

	// Validate requesting user matches requested user to be updated
	err = mustMatchContextUser(ctx, payload.UserID)
	if err != nil {
		handleError(ctx, apiErrors.ErrForbiden)
		return
	}

	claimedUser, err := store.GetUserByID(payload.UserID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(claimedUser.Password), []byte(payload.OldPassword)); err != nil {
		handleError(ctx, apiErrors.ErrForbiden)
		return
	}

	claimedUser.Password = payload.NewPassword
	errs := claimedUser.Validate()
	if errs != nil {
		render.Error(ctx, iris.StatusUnprocessableEntity, errs)
		return
	}

	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(claimedUser.Password), bcrypt.DefaultCost)
	if err != nil {
		handleError(ctx, err)
		return
	}

	claimedUser.Password = string(newPasswordHash)

	result, err := store.UpdateUserPassword(*claimedUser)
	if err != nil {
		handleError(ctx, err)
		return
	}

	// Hide password
	result.Password = ""
	render.JSON(ctx, iris.StatusOK, result)
}

// updateUserName is an API endpoint for changing a user's name.
func updateUserName(ctx iris.Context) {
	payload := updateUserNamePayload{}
	err := decodePayloadAndValidate(ctx, &payload)
	if err != nil {
		handleError(ctx, apiErrors.ErrUnprocessable)
		return
	}

	err = mustMatchContextUser(ctx, payload.UserID)
	if err != nil {
		handleError(ctx, apiErrors.ErrForbiden)
		return
	}

	claimedUser, err := store.GetUserByID(payload.UserID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	claimedUser.Name = payload.Name
	errs := claimedUser.Validate()
	if errs != nil {
		render.Error(ctx, iris.StatusUnprocessableEntity, errs)
		return
	}

	result, err := store.UpdateUserName(*claimedUser)
	if err != nil {
		handleError(ctx, err)
		return
	}

	// Hide password
	result.Password = ""
	render.JSON(ctx, iris.StatusOK, result)
}

// updateUserEmail is an API endpoint for changing a user's email.
func updateUserEmail(ctx iris.Context) {
	payload := updateUserEmailPayload{}
	err := decodePayloadAndValidate(ctx, &payload)
	if err != nil {
		handleError(ctx, apiErrors.ErrUnprocessable)
		return
	}

	err = mustMatchContextUser(ctx, payload.UserID)
	if err != nil {
		handleError(ctx, apiErrors.ErrForbiden)
		return
	}

	claimedUser, err := store.GetUserByID(payload.UserID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	claimedUser.Email = payload.Email
	errs := claimedUser.Validate()
	if errs != nil {
		render.Error(ctx, iris.StatusUnprocessableEntity, errs)
		return
	}

	result, err := store.UpdateUserEmail(*claimedUser)
	if err != nil {
		handleError(ctx, err)
		return
	}

	// Hide password
	result.Password = ""
	render.JSON(ctx, iris.StatusOK, result)
}

// decodeAndValidate decodes a model that implements the Validatable interface
// and calls the Validate function on it, returning any errros if something went wrong.
func decodeAndValidate(ctx iris.Context, m model.Validatable) error {
	if err := ctx.ReadJSON(m); err != nil {
		return apiErrors.ErrBadRequest
	}
	return m.Validate()
}

// mustMatchContextUser returns an error if the provided userID does not match
// the userID placed in the request context.
func mustMatchContextUser(ctx iris.Context, userID string) error {
	id, err := getSubjectID(ctx)
	if err != nil {
		return err
	}

	// Validate requesting user is the user being updated
	if userID != id {
		return errors.New("context user does not match request user id")
	}

	return nil
}
