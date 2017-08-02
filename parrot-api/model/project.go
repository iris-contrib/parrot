package model

import "github.com/iris-contrib/parrot/parrot-api/errors"

// ProjectStorer is the interface to store projects.
type ProjectStorer interface {
	GetProject(string) (*Project, error)
	CreateProject(Project) (*Project, error)
	UpdateProject(Project) (*Project, error)
	DeleteProject(string) error
	UpdateProjectName(projectID, name string) (*Project, error)
	AddProjectKey(projectID, key string) (*Project, error)
	UpdateProjectKey(projectID, oldKey, newKey string) (*Project, int, error)
	DeleteProjectKey(projectID, key string) (*Project, error)
}

// ProjectLocaleStorer is the interface to store project locales.
type ProjectLocaleStorer interface {
	UpdateLocalePairs(projID string, localeIdent string, pairs map[string]string) (*Locale, error)
	GetProjectLocaleByIdent(projID string, localeIdent string) (*Locale, error)
	GetProjectLocales(projID string, localeIdents ...string) ([]Locale, error)
}

var (
	ErrInvalidProjectName = &errors.Error{
		Type:    "InvalidProjectName",
		Message: "invalid field project name"}
)

type Project struct {
	ID   string   `db:"id" json:"id"`
	Name string   `db:"name" json:"name"`
	Keys []string `db:"keys" json:"keys"`
}

// SanitizeKeys removes empty and duplicate keys.
func (p *Project) SanitizeKeys() {
	var sk []string
	for _, key := range p.Keys {
		if key == "" {
			continue
		}
		if contains(sk, key) {
			continue
		}
		sk = append(sk, key)
	}

	p.Keys = sk
}

// contains returns true if parameter string 'str' is contained in parameter []string 'col'.
func contains(col []string, str string) bool {
	for _, v := range col {
		if v == str {
			return true
		}
	}
	return false
}

// Validate returns an error if the project's data is invalid.
func (p *Project) Validate() error {
	var errs []errors.Error
	if !HasMinLength(p.Name, 1) {
		errs = append(errs, *ErrInvalidProjectName)
	}
	if errs != nil {
		return NewValidationError(errs)
	}
	return nil
}
