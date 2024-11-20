package validatorx

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entrans "github.com/go-playground/validator/v10/translations/en"
)

var (
	validate   *validator.Validate
	translator ut.Translator
)

//nolint:errname
type ValidationErrors map[string]string

func (v ValidationErrors) Error() string {
	errorString := "Validation errors:\n"
	for field, err := range v {
		errorString += field + ": " + err + "\n"
	}
	return strings.TrimSpace(errorString)
}

func init() {
	validate = validator.New()
	english := en.New()
	uni := ut.New(english, english)
	translator, _ = uni.GetTranslator("en")
	if err := entrans.RegisterDefaultTranslations(validate, translator); err != nil {
		panic(fmt.Errorf("registering default translations: %w", err))
	}

	if err := registerCustomTranslations(); err != nil {
		panic(fmt.Errorf("registering custom translations: %w", err))
	}
}

func IsValidEmail(email string) bool {
	err := validate.Var(email, "email")
	return err == nil
}

func IsAbsoluteURL(url string) bool {
	err := validate.Var(url, "http_url")
	return err == nil
}

func ValidateStruct(s any) error {
	if err := validate.Struct(s); err != nil {
		var validationErrors validator.ValidationErrors
		if ok := errors.As(err, &validationErrors); !ok {
			return err
		}

		errorsMap := make(map[string]string)

		structType := reflect.ValueOf(s).Type()
		if structType.Kind() == reflect.Ptr {
			structType = structType.Elem()
		}

		for _, err := range validationErrors {
			jsonField := jsonFieldName(structType, err.StructField())
			errorsMap[jsonField] = err.Translate(translator)
		}

		return ValidationErrors(errorsMap)
	}

	return nil
}

func registerCustomTranslations() error {
	// Register translation for http_url tag
	return validate.RegisterTranslation("http_url", translator, func(ut ut.Translator) error {
		return ut.Add("http_url", "{0} must be a valid HTTP URL", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, err := ut.T("http_url", fe.Field())
		if err != nil {
			panic(fmt.Errorf("failed to translate http_url tag: %w", err))
		}
		return t
	})
}

func jsonFieldName(structType reflect.Type, fieldName string) string {
	if field, ok := structType.FieldByName(fieldName); ok {
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" {
			// Return the JSON field name, strip any tag options (e.g. omitempty)
			return strings.Split(jsonTag, ",")[0]
		}
	}
	return fieldName
}
