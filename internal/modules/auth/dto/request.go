package dto

type (
	LoginWithEmailRequest struct {
		Email      string `json:"email" validate:"required,email"`
		Password   string `json:"password" validate:"required,max=16"`
		RememberMe bool   `json:"remember_me" validate:"boolean"`
	}
)
