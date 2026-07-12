package dto

import "github.com/codenomdev/viona/pkg/validator"

type (
	LoginWithEmailRequest struct {
		Email      string `json:"email" validate:"required,email"`
		Password   string `json:"password" validate:"required,max=16"`
		RememberMe bool   `json:"remember_me" validate:"boolean"`
	}
	RegisterWithEmailRequest struct {
		Fullname string `json:"fullname" validate:"required,min=4,max=170"`
		Email    string `json:"email" validate:"required,email,max=200"`
		Password string `json:"password" validate:"required,min=8,max=16"`
	}
)

func (r *RegisterWithEmailRequest) Check() ([]*validator.FormErrorField, error) {
	return nil, nil
}
