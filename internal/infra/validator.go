package infra

import "github.com/go-playground/validator/v10"

type Validator interface {
	Struct(i any) error
}

type validatorImpl struct {
	validator *validator.Validate
}

func NewValidatorService() Validator {
	return &validatorImpl{validator: validator.New()}
}

func (v *validatorImpl) Struct(i any) error {
	return v.validator.Struct(i)
}
