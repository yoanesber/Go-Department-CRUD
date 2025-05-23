package validator

import (
	"reflect"
	"strings"
	"sync"

	"gopkg.in/go-playground/validator.v9"
)

var (
	once     sync.Once
	validate *validator.Validate
)

// Init initializes the validator and registers custom validations.
func InitValidator() {
	once.Do(func() {
		validate = validator.New()

		// Register tag name function to use JSON field names if available
		// instead of struct field names
		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			tag := fld.Tag.Get("json")
			if tag == "-" || tag == "" {
				return fld.Name // fallback ke nama field struct
			}
			return strings.Split(tag, ",")[0]
		})
	})
}

// GetValidator returns the initialized validator instance.
func GetValidator() *validator.Validate {
	return validate
}
