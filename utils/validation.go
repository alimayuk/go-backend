package utils

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var customMessages = map[string]string{
	"Name.required":     "İsim zorunludur",
	"Name.min":          "İsim en az {param} karakter olmalı",
	"Name.max":          "İsim en fazla {param} karakter olabilir",
	"Title.required":    "Başlık zorunludur",
	"Title.min":         "Başlık en az {param} karakter olmalı",
	"Title.max":         "Başlık en fazla {param} karakter olabilir",
	"Email.required":    "E-posta zorunludur",
	"Email.email":       "Geçerli bir e-posta giriniz",
	"Password.required": "Şifre zorunludur",
	"Password.min":      "Şifre en az {param} karakter olmalı",
	"Role.required":     "Rol zorunludur",
	"Role.oneof":        "Rol 'admin' veya 'user' olmalıdır",
}

func FormatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrors {
			field := strings.ToLower(fieldErr.Field())
			tag := fieldErr.Tag()
			param := fieldErr.Param()

			key := fieldErr.Field() + "." + tag

			if msg, exists := customMessages[key]; exists {
				msg = strings.ReplaceAll(msg, "{param}", param)
				errors[field] = msg
			} else {
				errors[field] = "Geçersiz değer"
			}
		}
	}

	return errors
}
