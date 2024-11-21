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

	if err := registerCustomValidations(); err != nil {
		panic(fmt.Errorf("registering custom validations: %w", err))
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
	translationFunc := func(tag, message string, paramFunc func(fe validator.FieldError) []string) error {
		return validate.RegisterTranslation(tag, translator, func(ut ut.Translator) error {
			return ut.Add(tag, message, true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, err := ut.T(tag, paramFunc(fe)...)
			if err != nil {
				panic(fmt.Errorf("translating tag %s: %w", tag, err))
			}
			return t
		})
	}

	translations := map[string]struct {
		message   string
		paramFunc func(fe validator.FieldError) []string
	}{
		"http_url":          {"Field must be a valid HTTP URL", nil},
		"empty_or_http_url": {"Field must be a valid HTTP URL", nil},
		"required":          {"Field is required", nil},
		"min": {
			"Field must be at least {0} character(s) in length",
			func(fe validator.FieldError) []string { return []string{fe.Param()} },
		},
	}

	for tag, t := range translations {
		paramFunc := t.paramFunc
		if paramFunc == nil {
			paramFunc = func(fe validator.FieldError) []string { return nil }
		}

		if err := translationFunc(tag, t.message, paramFunc); err != nil {
			return fmt.Errorf("registering %s translation: %w", tag, err)
		}
	}

	return nil
}

func registerCustomValidations() error {
	return validate.RegisterValidation("empty_or_http_url", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		if value == "" {
			return true
		}
		err := validate.Var(value, "http_url")
		return err == nil
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
