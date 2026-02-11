package lib

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateAndRespond(w http.ResponseWriter, s any) bool {
	if err := validate.Struct(s); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			RespondError(w, http.StatusBadRequest, "invalid request")
			return false
		}

		messages := make([]string, 0, len(validationErrors))
		for _, fe := range validationErrors {
			messages = append(messages, fmt.Sprintf("%s is %s", fe.Field(), fe.Tag()))
		}

		RespondError(w, http.StatusBadRequest, strings.Join(messages, ", "))
		return false
	}
	return true
}
