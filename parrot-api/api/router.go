// Package api handles the creation and configuration of all API resource routes.
package api

import (
	"github.com/kataras/iris/v12"
	"github.com/iris-contrib/parrot/parrot-api/auth"
	"github.com/iris-contrib/parrot/parrot-api/datastore"
)

// TODO: inject store via closures instead of keeping global var
var store datastore.Store

// NewRouter creates an API router based on the parameter datastore and token provider.
// It registers and configures all necessary routes.
func NewRouter(ds datastore.Store, tp auth.TokenProvider) iris.Configurator {
	store = ds
	mustHaveValidToken := tokenMiddleware(tp)

	return func(app *iris.Application) {
		app.PartyFunc("/api/v1", func(router iris.Party) {
			router.Use(enforceContentTypeJSON)

			router.Get("/ping", ping)
			router.Post("/users/register", createUser)

			router.PartyFunc("/users", func(r1 iris.Party) {
				// Past this point, all routes will require a valid token
				r1.Use(mustHaveValidToken)

				r1.PartyFunc("/self", func(r2 iris.Party) {
					r2.Get("/", getUserSelf)
					r2.Patch("/name", updateUserName)
					r2.Patch("/email", updateUserEmail)
					r2.Patch("/password", updateUserPassword)
				})
			})

			router.PartyFunc("/projects", func(r1 iris.Party) {
				// Past this point, all routes will require a valid token
				r1.Use(mustHaveValidToken)

				r1.Get("/", getUserProjects)
				r1.Post("/", createProject)

				r1.PartyFunc("/{projectID}", func(r2 iris.Party) {
					r2.Get("/", mustAuthorize(canViewProject), showProject)
					r2.Delete("/", mustAuthorize(canDeleteProject), deleteProject)

					r2.Patch("/name", mustAuthorize(canUpdateProject), updateProjectName)

					r2.Post("/keys", mustAuthorize(canUpdateProject), addProjectKey)
					r2.Patch("/keys", mustAuthorize(canUpdateProject), updateProjectKey)
					r2.Delete("/keys", mustAuthorize(canUpdateProject), deleteProjectKey)

					r2.PartyFunc("/users", func(r3 iris.Party) {
						r3.Get("/", mustAuthorize(canViewProjectRoles), getProjectUsers)
						r3.Post("/", mustAuthorize(canAssignProjectRoles), assignProjectUser)
						r3.Patch("/{userID}/role", mustAuthorize(canUpdateProjectRoles), updateProjectUserRole)
						r3.Delete("/{userID}", mustAuthorize(canRevokeProjectRoles), revokeProjectUser)
					})

					r2.PartyFunc("/clients", func(r3 iris.Party) {
						r3.Get("/", mustAuthorize(canManageAPIClients), getProjectClients)
						r3.Get("/{clientID}", mustAuthorize(canManageAPIClients), getProjectClient)
						r3.Post("/", mustAuthorize(canManageAPIClients), createProjectClient)
						r3.Patch("/{clientID}/resetSecret", mustAuthorize(canManageAPIClients), resetProjectClientSecret)
						r3.Patch("/{clientID}/name", mustAuthorize(canManageAPIClients), updateProjectClientName)
						r3.Delete("/{clientID}", mustAuthorize(canManageAPIClients), deleteProjectClient)
					})

					r2.PartyFunc("/locales", func(r3 iris.Party) {
						r3.Get("/", mustAuthorize(canViewLocales), findLocales)
						r3.Post("/", mustAuthorize(canCreateLocales), createLocale)

						r3.PartyFunc("/{localeIdent}", func(r4 iris.Party) {
							r4.Get("/", mustAuthorize(canViewLocales), showLocale)
							r4.Patch("/pairs", mustAuthorize(canUpdateLocales), updateLocalePairs)
							r4.Delete("/", mustAuthorize(canDeleteLocales), deleteLocale)

							r4.Get("/export/{type}", mustAuthorize(canExportLocales), exportLocale)
						})
					})

				})

			})

		})
	}
}
