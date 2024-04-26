package validator

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entrans "github.com/go-playground/validator/v10/translations/en"
)

var (
	Validate   *validator.Validate
	translator ut.Translator
)

func init() {
	Validate = validator.New()
	english := en.New()
	uni := ut.New(english, english)
	translator, _ = uni.GetTranslator("en")
	_ = entrans.RegisterDefaultTranslations(Validate, translator)
}

func IsValidEmail(email string) bool {
	err := Validate.Var(email, "email")
	return err == nil
}

func IsAbsoluteURL(url string) bool {
	err := Validate.Var(url, "http_url")
	return err == nil
}
