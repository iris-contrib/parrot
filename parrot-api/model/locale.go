// Package model holds the various types and interfaces for Parrot.
package model

import "github.com/iris-contrib/parrot/parrot-api/errors"

var (
	ErrInvalidLocaleIdent = &errors.Error{
		Type:    "InvalidLocaleIdent",
		Message: "invalid field locale ident"}
	ErrInvalidLocaleLanguage = &errors.Error{
		Type:    "InvalidLocaleLanguage",
		Message: "invalid field locale language"}
	ErrInvalidLocaleCountry = &errors.Error{
		Type:    "InvalidLocaleCountry",
		Message: "invalid field locale country"}
)

// LocaleStorer is the interface to store locales.
type LocaleStorer interface {
	CreateLocale(loc Locale) (*Locale, error)
	DeleteLocale(projID string, ident string) error
}

type Locale struct {
	ID        string            `db:"id" json:"id"`
	Ident     string            `db:"ident" json:"ident"`
	Language  string            `db:"language" json:"language"`
	Country   string            `db:"country" json:"country"`
	Pairs     map[string]string `db:"pairs" json:"pairs"`
	ProjectID string            `db:"project_id" json:"project_id"`
}

// Validate returns an error if the locale's data is invalid.
func (loc *Locale) Validate() error {
	var errs []errors.Error
	if !HasMinLength(loc.Ident, 2) {
		errs = append(errs, *ErrInvalidLocaleIdent)
	}
	if !HasMinLength(loc.Language, 1) {
		errs = append(errs, *ErrInvalidLocaleLanguage)
	}
	if !HasMinLength(loc.Country, 1) {
		errs = append(errs, *ErrInvalidLocaleCountry)
	}
	if errs != nil {
		return NewValidationError(errs)
	}
	return nil
}

// SyncKeys will add new keys from string slice t to document pairs.
func (loc *Locale) SyncKeys(t []string) {
	if loc.Pairs == nil {
		loc.Pairs = make(map[string]string)
	}

	temp := make(map[string]string)

	// Assign each key, if it's already there it will simply reassign to there
	// previous value, otherwise an empty string will be set
	for _, v := range t {
		temp[v] = loc.Pairs[v]
	}

	loc.Pairs = temp
}
