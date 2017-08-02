package model

import "github.com/iris-contrib/parrot/parrot-api/errors"

var (
	ErrInvalidClientName = &errors.Error{
		Type:    "InvalidClientName",
		Message: "invalid field name"}
	ErrInvalidProjectID = &errors.Error{
		Type:    "InvalidProjectID",
		Message: "invalid field project_id"}
)

type ProjectClient struct {
	ClientID  string `db:"client_id" json:"client_id"`
	Name      string `db:"name" json:"name"`
	Secret    string `db:"secret" json:"secret,omitempty"`
	ProjectID string `db:"project_id" json:"project_id"`
}

// ProjectClientStorer is the interface to store project clients.
type ProjectClientStorer interface {
	FindOneClient(string) (*ProjectClient, error)
	GetProjectClients(string) ([]ProjectClient, error)
	GetProjectClient(projectID, clientID string) (*ProjectClient, error)
	CreateProjectClient(ProjectClient) (*ProjectClient, error)
	UpdateProjectClientSecret(ProjectClient) (*ProjectClient, error)
	UpdateProjectClientName(ProjectClient) (*ProjectClient, error)
	DeleteProjectClient(projectID, clientID string) error
}

// Validate returns an error if the project client's data is invalid.
func (p *ProjectClient) Validate() error {
	var errs []errors.Error
	if !HasMinLength(p.Name, 1) {
		errs = append(errs, *ErrInvalidClientName)
	}
	if errs != nil {
		return NewValidationError(errs)
	}
	return nil
}
