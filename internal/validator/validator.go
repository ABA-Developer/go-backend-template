package validator

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

func NewValidator() *Validator {
	return &Validator{Validator: validator.New()}
}

// Validator wrapper.
type Validator struct {
	Validator *validator.Validate
}

func (v Validator) Validate(i interface{}) (err error) {
	err = v.Validator.Struct(i)

	return
}

func ValidationErrors(err error) (message map[string]interface{}) {
	ve, ok := err.(validator.ValidationErrors)
	if !ok {
		return
	}

	message = make(map[string]interface{})

	for _, e := range ve {
		switch e.Tag() {
		case "datetime":
			message[strings.ToLower(convertCase(e.Field(), '_'))] = fmt.Sprintf("Field %s must be date & time", convertCase(e.Field(), ' '))
		case "email":
			message[strings.ToLower(convertCase(e.Field(), '_'))] = "Input must be valid email address"
		case "max":
			message[strings.ToLower(convertCase(e.Field(), '_'))] = fmt.Sprintf("Field %s must be less than %s", convertCase(e.Field(), ' '), e.Param())
		case "min":
			message[strings.ToLower(convertCase(e.Field(), '_'))] = fmt.Sprintf("Field %s must be more than %s", convertCase(e.Field(), ' '), e.Param())
		case "required":
			message[strings.ToLower(convertCase(e.Field(), '_'))] = fmt.Sprintf("Field %s can not empty!", convertCase(e.Field(), ' '))
		case "oneof":
			message[strings.ToLower(convertCase(e.Field(), '_'))] = fmt.Sprintf("Field %s must be one of [%s]", convertCase(e.Field(), ' '), e.Param())
		}
	}

	return
}

func convertCase(t string, c rune) string {
	buf := &bytes.Buffer{}

	for i, r := range t {
		if i > 0 && unicode.IsUpper(r) {
			if t[i-1] != 'I' && r != 'D' {
				buf.WriteRune(c)
			}
		}

		buf.WriteRune(r)
	}

	return buf.String()
}
