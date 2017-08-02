package api

import (
	"github.com/kataras/iris"

	apiErrors "github.com/iris-contrib/parrot/parrot-api/errors"
)

// Role is an identifier for a group of assigned grants.
type Role string

// Role grant is an identifier for the right to perform an action.
type RoleGrant string

// known roles
const (
	ownerRole     = "owner"
	editorRole    = "editor"
	viewerRole    = "viewer"
	clientRole    = "client"
	developerRole = "developer"
)

// known grants
const (
	canAssignProjectRoles = "CanAssignProjectRoles"
	canRevokeProjectRoles = "CanRevokeProjectRoles"
	canUpdateProjectRoles = "CanUpdateProjectRoles"
	canViewProjectRoles   = "CanViewProjectRoles"
	canUpdateProject      = "CanUpdateProject"
	canDeleteProject      = "CanDeleteProject"
	canViewProject        = "CanViewProject"
	canCreateLocales      = "CanCreateLocales"
	canUpdateLocales      = "CanUpdateLocales"
	canDeleteLocales      = "CanDeleteLocales"
	canViewLocales        = "CanViewLocales"
	canManageAPIClients   = "CanManageAPIClients"
	canExportLocales      = "CanExportLocales"
)

// permissions mapping of Roles to Grants.
var permissions = map[Role][]RoleGrant{
	ownerRole: []RoleGrant{
		canAssignProjectRoles,
		canRevokeProjectRoles,
		canUpdateProjectRoles,
		canViewProjectRoles,
		canUpdateProject,
		canDeleteProject,
		canViewProject,
		canCreateLocales,
		canUpdateLocales,
		canDeleteLocales,
		canViewLocales,
		canManageAPIClients,
		canExportLocales,
	},
	editorRole: []RoleGrant{
		canViewProjectRoles,
		canUpdateProject,
		canViewProject,
		canCreateLocales,
		canUpdateLocales,
		canDeleteLocales,
		canViewLocales,
		canExportLocales,
	},
	viewerRole: []RoleGrant{
		canViewProjectRoles,
		canViewProject,
		canViewLocales,
		canExportLocales,
	},
	clientRole: []RoleGrant{
		canExportLocales,
	},
	developerRole: []RoleGrant{
		canViewProjectRoles,
		canViewProject,
		canViewLocales,
		canExportLocales,
		canManageAPIClients,
	},
}

// isRole returns true if the provided string can be casted to a known role.
func isRole(r string) bool {
	v := Role(r)
	switch v {
	case ownerRole, editorRole, viewerRole, developerRole:
		return true
	}
	return false
}

// isAllowed returns true if the provided role has the provided grant.
func isAllowed(r Role, a RoleGrant) bool {
	actions, ok := permissions[r]
	if !ok {
		return false
	}
	for _, currentAction := range actions {
		if currentAction == a {
			return true
		}
	}
	return false
}

// getProjectUserRole returns the role a user has for a given project
// or returns an error (for example, the user was not found).
func getProjectUserRole(projID, userID string) (string, error) {
	user, err := store.GetProjectUser(projID, userID)
	if err != nil {
		return "", err
	}
	return user.Role, nil
}

// mustBeProjectClient returns an error if no client with provided clientID
// exists for a given project.
func mustBeProjectClient(projID, clientID string) error {
	client, err := store.GetProjectClient(projID, clientID)
	if err != nil || client == nil {
		return err
	}
	return nil
}

// mustAuthorize authorizes or denies requests based on required rights for action.
// Identifies if requesting subject is able to perform action on the particular project.
func mustAuthorize(action RoleGrant) iris.Handler {
	return func(ctx iris.Context) {
		projectID := ctx.Params().Get("projectID")
		if projectID == "" {
			handleError(ctx, apiErrors.ErrBadRequest)
			return
		}

		subType, err := getSubjectType(ctx)
		if err != nil {
			handleError(ctx, apiErrors.ErrBadRequest)
			return
		}
		requesterID, err := getSubjectID(ctx)
		if err != nil {
			handleError(ctx, apiErrors.ErrBadRequest)
			return
		}

		var requesterRole string

		switch subType {
		case userSubject:
			requesterRole, err = getProjectUserRole(projectID, requesterID)
			if err != nil {
				handleError(ctx, err)
				return
			}
		case clientSubject:
			err := mustBeProjectClient(projectID, requesterID)
			if err != nil {
				handleError(ctx, err)
				return
			}
			requesterRole = clientRole
		default:
			handleError(ctx, apiErrors.ErrBadRequest)
			return
		}

		if !isAllowed(Role(requesterRole), action) {
			handleError(ctx, apiErrors.ErrForbiden)
			return
		}

		ctx.Next()
	}
}
