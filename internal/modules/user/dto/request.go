package dto

type (
	RequestCreateNew struct {
		FullName string `json:"fullname" validate:"required,min=4,max=200"`
		Username string `json:"username" validate:"required,min=4,max=50"`
		Email    string `json:"email" validate:"required,email,max=230"`
		Password string `json:"password" validate:"required,min=8,max=16"`
	}
	GetUserByEmailRequest struct {
		Email string `json:"email" validate:"required,email"`
	}
)
