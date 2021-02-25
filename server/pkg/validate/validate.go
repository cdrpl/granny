package validate

import "regexp"

// Validator is used to validate inputs. Any errors are added to the errors slice.
type Validator struct {
	Errors []string
}

// CreateValidator will create a Validator struct and return it.
func CreateValidator() *Validator {
	return &Validator{
		Errors: make([]string, 0),
	}
}

// HasError will return true if the errors slice len is greater than 0.
func (v *Validator) HasError() bool {
	return len(v.Errors) > 0
}

// FirstError returns the true and the first error in the errors slice. Returns false and an empty string if no errors are present.
func (v *Validator) FirstError() (bool, string) {
	if v.HasError() {
		return true, v.Errors[0]
	}
	return false, ""
}

// Required will create an error if the input len is 0.
func (v *Validator) Required(input, message string) {
	if len(input) == 0 {
		v.Errors = append(v.Errors, message)
	}
}

// Min will append an error if the input len is less than min.
func (v *Validator) Min(input, message string, min int) {
	if len(input) < min {
		v.Errors = append(v.Errors, message)
	}
}

// Max will append an error if the input len is greater than max.
func (v *Validator) Max(input, message string, max int) {
	if len(input) > max {
		v.Errors = append(v.Errors, message)
	}
}

// IsEmail will append an error if the email is not in valid format
func (v *Validator) IsEmail(input, message string) {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !re.MatchString(input) {
		v.Errors = append(v.Errors, message)
	}
}
