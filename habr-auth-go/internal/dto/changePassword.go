package dto

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

type ChangePasswordRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	NewPassword string `json:"newPassword"`
}

func (u *ChangePasswordRequest) IsValid() error {
	if u.Email == "" || u.Password == "" || u.NewPassword == "" {
		return errors.New("Missing required fields")
	}
	if err := u.validatePass(); err != nil {
		return err
	}

	return nil
}

func (u *ChangePasswordRequest) validatePass() error {
	if len(u.NewPassword) < 8 {
		return errors.New("New password is short (Minimum 8 characters)")
	}
	if len(u.NewPassword) > 64 {
		return errors.New("The new password is too long")
	}
	allowedSumbols := "()*_-+=%\""
	var hasLetter, hasDigit bool
	for _, ch := range u.NewPassword {
		switch {
		case ch >= 'a' && ch <= 'z':
			hasLetter = true
		case ch >= 'A' && ch <= 'Z':
			hasLetter = true
		case unicode.IsDigit(ch):
			hasDigit = true
		case strings.ContainsRune(allowedSumbols, ch):
		default:
			return errors.New(fmt.Sprintf("The new password can only contain Latin letters, numbers and symbols %s", allowedSumbols))
		}
	}
	if !hasLetter || !hasDigit {
		return errors.New("The new password must contain letters and numbers")
	}

	return nil
}
