package model

import "github.com/iris-contrib/parrot/parrot-api/errors"
import "strings"

var (
	ErrInvalidEmail = &errors.Error{
		Type:    "InvalidEmail",
		Message: "invalid email"}
	ErrInvalidName = &errors.Error{
		Type:    "InvalidName",
		Message: "invalid name"}
	ErrInvalidPassword = &errors.Error{
		Type:    "InvalidPassword",
		Message: "invalid password"}
)

// UserStorer is the interface to store users.
type UserStorer interface {
	GetUserByID(string) (*User, error)
	GetUserByEmail(string) (*User, error)
	CreateUser(User) (*User, error)
	UpdateUserPassword(User) (*User, error)
	UpdateUserName(User) (*User, error)
	UpdateUserEmail(User) (*User, error)
}

type User struct {
	ID       string `db:"id" json:"id"`
	Name     string `db:"name" json:"name,omitempty"`
	Email    string `db:"email" json:"email,omitempty"`
	Password string `db:"password" json:"password,omitempty"`
}

func (u *User) Normalize() {
	u.Email = strings.ToLower(u.Email)
}

// Validate returns an error if the user's data is invalid.
// It will normalize the user data before validating
func (u *User) Validate() error {
	u.Normalize()

	var errs []errors.Error
	if !ValidEmail(u.Email) {
		errs = append(errs, *ErrInvalidEmail)
	}
	if !HasMinLength(strings.Trim(u.Name, " "), 1) {
		errs = append(errs, *ErrInvalidName)
	}
	if !HasMinLength(u.Password, 8) {
		errs = append(errs, *ErrInvalidPassword)
	}
	if errs != nil {
		return NewValidationError(errs)
	}
	return nil
}
