package validator

import (
	"errors"
	"fmt"

	"github.com/codenomdev/viona/pkg/response"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/labstack/echo/v4"
)

// Use a single instance of Validate, it caches struct info
var (
	validate *validator.Validate
	uni      *ut.UniversalTranslator
)

type Validation struct {
	validation *validator.Validate
	trans      ut.Translator
}

type fieldError struct {
	Namespace string      `json:"namespace"`
	Field     string      `json:"field_name"`
	Type      string      `json:"type_validation"`
	Message   string      `json:"message"`
	Error     string      `json:"error"`
	Value     interface{} `json:"value"`
}

func New() *Validation {
	//init translate validator
	en := en.New()
	uni = ut.New(en, en)

	// this is usually know or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	trans, _ := uni.GetTranslator("en")
	validate = validator.New()
	// validate.RegisterValidation("code_challenge", ValidateCodeChallenge)
	// validate.RegisterValidation("code_challenge_method", ValidateCodeChallengeMethod)
	// validate.RegisterValidation("code_verifier", ValidateCodeVerifier)

	// Register default lang validation
	en_translations.RegisterDefaultTranslations(validate, trans)

	// register translation
	// RegisterSimpleTranslation(validate, trans, "code_challenge", "{0} must be a valid Base64URL (43–128 chars, no '=')")
	// RegisterSimpleTranslation(validate, trans, "code_challenge_method", "{0} must be 'S256'")
	// RegisterSimpleTranslation(validate, trans, "code_verifier", "{0} must be a valid PKCE code_verifier (43–128 chars, [A–Z, a–z, 0–9, '-', '.', '_', '~'])")

	return &Validation{
		validation: validate,
		trans:      trans,
	}
}

// Validate struct ctx
// this func with binding request & translate message
// custom return for API response
func (v *Validation) ValidateStructCtx(ctx echo.Context, s interface{}) (response.RestError, error) {
	if err := ctx.Bind(s); err != nil {
		return nil, err
	}

	err := v.validation.StructCtx(ctx.Request().Context(), s)
	var vError validator.ValidationErrors

	if errors.As(err, &vError) {
		// translate all error at once
		errV := err.(validator.ValidationErrors)
		messages := make([]fieldError, len(errV))
		for i, err := range errV {
			// returns a map with key = namespace & value = translated error
			// NOTICE: 2 errors are returned and you'll see something surprising
			// translations are i18n aware!!!!
			// eg. '10 characters' vs '1 character'
			messages[i] = fieldError{
				Namespace: err.Namespace(),
				Field:     err.Field(),
				Type:      err.ActualTag(),
				Message:   err.Translate(v.trans),
				Error:     fmt.Sprintf("%s %s", err.Tag(), err.Param()),
				Value:     err.Value(),
			}
		}
		return response.NewHttpValidationError(
			messages,
		), nil
	}

	return nil, nil
}

// Validate struct ctx
// this func original call github.com/go-playground/validator/v10
// return default error
func (v *Validation) StructCtx(ctx echo.Context, s interface{}) error {
	if err := ctx.Bind(s); err != nil {
		return err
	}

	return v.validation.StructCtx(ctx.Request().Context(), s)
}

// Validate struct
// this func not with binding request
// custom return
func (v *Validation) ValidateStruct(s interface{}) (response.RestError, error) {
	err := v.validation.Struct(s)

	if err != nil {
		// translate all error at once
		errV := err.(validator.ValidationErrors)
		messages := make([]*fieldError, len(errV))
		for i, err := range errV {
			// returns a map with key = namespace & value = translated error
			// NOTICE: 2 errors are returned and you'll see something surprising
			// translations are i18n aware!!!!
			// eg. '10 characters' vs '1 character'
			messages[i] = &fieldError{
				Namespace: err.Namespace(),
				Field:     err.Field(),
				Type:      err.ActualTag(),
				Message:   err.Translate(v.trans),
				Error:     fmt.Sprintf("%s %s", err.Tag(), err.Param()),
				Value:     err.Value(),
			}
		}
		return response.NewHttpValidationError(
			messages,
		), nil
	}

	return nil, nil
}

// Validate struct
// this func not with binding request
// return default error original from github.com/go-playground/validator/v10
func (v *Validation) Struct(s interface{}) error {
	return v.validation.Struct(s)
}
