package api

import (
	"errors"

	"github.com/kataras/iris"
)

var (
	ErrInvalidPayload = errors.New("invalid payload")
)

type ValidatablePayload interface {
	Validate() error
}

type updatePasswordPayload struct {
	UserID      string `json:"userId"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

func (p *updatePasswordPayload) Validate() error {
	if p.NewPassword == "" || p.OldPassword == "" || p.UserID == "" {
		return ErrInvalidPayload
	}
	return nil
}

type updateUserNamePayload struct {
	UserID string `json:"userId"`
	Name   string `json:"name"`
}

func (p *updateUserNamePayload) Validate() error {
	if p.Name == "" || p.UserID == "" {
		return ErrInvalidPayload
	}
	return nil
}

type updateUserEmailPayload struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
}

func (p *updateUserEmailPayload) Validate() error {
	if p.Email == "" || p.UserID == "" {
		return ErrInvalidPayload
	}
	return nil
}

func decodePayloadAndValidate(ctx iris.Context, p ValidatablePayload) error {
	err := ctx.ReadJSON(&p)
	if err != nil {
		return err
	}

	if err := p.Validate(); err != nil {
		return err
	}

	return nil
}
