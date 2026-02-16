package lib

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Validate(s any) error {
	if err := validate.Struct(s); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return fmt.Errorf("invalid request")
		}

		messages := make([]string, 0, len(validationErrors))
		for _, fe := range validationErrors {
			messages = append(messages, fmt.Sprintf("%s is %s", fe.Field(), fe.Tag()))
		}

		return errors.New(strings.Join(messages, ", "))
	}
	return nil
}
