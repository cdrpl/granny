package validate_test

import (
	"testing"

	. "github.com/cdrpl/idlemon/pkg/validate"
)

func createValidatorWithErrors(errMsg string) Validator {
	v := Validator{Errors: make([]string, 1)}
	v.Errors[0] = errMsg
	return v
}

func TestCreateValidator(t *testing.T) {
	// it should make the errors slice
	v := CreateValidator()

	if v.Errors == nil {
		t.Error("CreateValidator() returns a Validator with nil Errors")
	}
}

func TestHasError(t *testing.T) {
	// it should return true if an error is present
	v := createValidatorWithErrors("has mock error")
	if !v.HasError() {
		t.Error("HasError() failed to detect error")
	}

	// it should return false if an error is not present
	v = Validator{Errors: make([]string, 0)}
	if v.HasError() {
		t.Error("HasError() falsely detected an error")
	}
}

func TestFirstError(t *testing.T) {
	// it should return true, message if errors are present
	v := createValidatorWithErrors("a mock error")
	hasErr, msg := v.FirstError()
	if !hasErr {
		t.Error("FirstError() should return true if errors are present")
	}
	if msg != "a mock error" {
		t.Error("FirstError() should return the first error message if errors are present")
	}

	// it should return false, "" if there are no errors
	v = Validator{Errors: make([]string, 0)}
	hasErr, msg = v.FirstError()
	if hasErr {
		t.Error("FirstError() should return false when no errors are present")
	}
	if msg != "" {
		t.Error("FirstError() should return an empty string when no errors are present")
	}
}

func TestRequired(t *testing.T) {
	// it should catch an emptry string
	v := Validator{make([]string, 1)}
	v.Required("", "required err msg")
	if v.Errors[1] != "required err msg" {
		t.Error("Required() did not catch the error")
	}

	// it should not return an error if the input is valid
	v.Errors = make([]string, 0)
	v.Required("theinput", "")
	if len(v.Errors) != 0 {
		t.Error("Required() falsely detected an error")
	}
}

func TestMin(t *testing.T) {
	// it should catch the error
	v := Validator{make([]string, 0)}
	v.Min("1234", "min err msg", 5)
	if v.Errors[0] != "min err msg" {
		t.Error("Min() did not catch the error")
	}

	// it should not return an error if the input is valid
	v.Errors = make([]string, 0)
	v.Min("12345", "", 5)
	if len(v.Errors) != 0 {
		t.Error("Min() falsely detected an error")
	}
}

func TestMax(t *testing.T) {
	// it should catch the error
	v := Validator{make([]string, 0)}
	v.Max("123456", "max err msg", 5)
	if v.Errors[0] != "max err msg" {
		t.Error("Max() did not catch the error")
	}

	// it should not return an error if the input is valid
	v.Errors = make([]string, 0)
	v.Max("12345", "", 5)
	if len(v.Errors) != 0 {
		t.Error("Max() falsely detected an error")
	}
}

func TestIsEmail(t *testing.T) {
	// it should catch the error
	v := Validator{make([]string, 0)}
	v.IsEmail("fakeemail@", "email is not valid")
	if v.Errors[0] != "email is not valid" {
		t.Error("IsEmail() did not catch the error")
	}

	// it should not return an error if the input is valid
	v.Errors = make([]string, 0)
	v.IsEmail("valid@email.gmail.com", "")
	if len(v.Errors) != 0 {
		t.Error("IsEmail() falsely detected an error")
	}
}
