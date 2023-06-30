package forms

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

var formMap = map[string]string{
	"first_name": "First Name",
	"last_name":  "Last Name",
	"email":      "Email",
	"phone":      "Phone",
}

// creates new custom Form
type Form struct {
	url.Values
	Errors errors
}

func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Has verifies if form has field filled
func (f *Form) Has(field string) bool {
	x := f.Get(field)
	return x != ""
}

// Required checks for required fields
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field should not be blank")
		}
	}
}

// IsValid returns true if there are no errors in the form
func (f *Form) IsValid() bool {
	return len(f.Errors) == 0
}

// MinLength checks for min length of string
func (f *Form) MinLength(field string, length int) bool {
	currVal := f.Get(field)
	if len(currVal) < length {
		f.Errors.Add(field, fmt.Sprintf("%s field should have length more than: %d", formMap[field], length))
		return false
	}
	return true
}

// IsEmail checks if entered email is correct
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid Email address")
	}
}
