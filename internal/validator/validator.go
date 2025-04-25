package validator

import (
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

// compiled email regex pattern
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

// / is valid if no errors
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

// add a message to the non-field errors map
func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

// / add a message to the validators' field errors map
func (v *Validator) AddFieldError(key, message string) {
	// ensure map is created
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	// add the error
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// / func only adds the field error if the check is not 'ok'
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// / returns true if <value> is not empty
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// / returns true if <value> has less than <n> characters
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// / returns true if <value> is in <permittedValues>
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

// / returns true if a value contains at least n chars
func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

// / returns true if value matches a compiled regex pattern
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
