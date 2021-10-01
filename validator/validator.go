package validator

import (
	"discountDealer/models/userModels"
	"discountDealer/repository/userRepository"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func Validate(entity interface{}) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(entity)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

func UserExist(username string) bool {
	vUser := &userModels.User{Username: username}
	userDB := userRepository.Init()
	err := userDB.GetByFilter("username = ?", vUser.Username, vUser)
	return err == nil
}

func init() {
	validate = validator.New()
	validate.RegisterValidation("username-unique", func(fl validator.FieldLevel) bool {
		return !UserExist(fl.Field().String())
	})
}
