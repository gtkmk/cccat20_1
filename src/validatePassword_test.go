package src

import "testing"

func TestValidatePassword_Valid(t *testing.T) {
	tests := []string{
		"asdFGH123",
		"asdG123456",
		"aG1aG129999",
	}

	for _, password := range tests {
		if !ValidatePassword(password) {
			t.Errorf("ValidatePassword(%s) = false; want true", password)
		}
	}
}

func TestValidatePassword_Invalid(t *testing.T) {
	tests := []string{
		"",
		"asD123",
		"12345678",
		"asdfghjkl",
		"ASDFGHJKL",
		"asddfg123456",
	}

	for _, password := range tests {
		if ValidatePassword(password) {
			t.Errorf("ValidatePassword(%s) = true; want false", password)
		}
	}
}
