package dto

import (
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"strings"
	"unicode"
)

type RegisterRequest struct {
	Email        string `json:"email"`
	PasswordHash string `json:"password"`
	Nickname     string `json:"nickname"`
}

func (u *RegisterRequest) IsValid() error {
	if u.Email == "" || u.PasswordHash == "" || u.Nickname == "" {
		return errors.New("Missing required fields")
	}
	if err := u.validateEmail(); err != nil {
		return err
	}
	if err := u.validatePass(); err != nil {
		return err
	}
	if err := u.validateNick(); err != nil {
		return err
	}
	return nil
}

func (u *RegisterRequest) validateEmail() error {
	return validation.Validate(u.Email,
		validation.Required,
		is.Email,
	)
}

func (u *RegisterRequest) validateNick() error {
	if len(u.Nickname) < 3 {
		return errors.New("Nickname is short (Minimum 3 characters)")
	}
	if len(u.Nickname) > 25 {
		return errors.New("The nickname is too long")
	}
	allowedSumbols := "-_"
	var hasLetter, hasDigit bool
	for _, ch := range u.Nickname {
		switch {
		case ch >= 'a' && ch <= 'z':
			hasLetter = true
		case ch >= 'A' && ch <= 'Z':
			hasLetter = true
		case unicode.IsDigit(ch):
			hasDigit = true
		case strings.ContainsRune(allowedSumbols, ch):
		default:
			return errors.New("The nickname can only contain Latin letters, numbers and symbols dashes and underlines")
		}
	}
	if !hasLetter || !hasDigit {
		return errors.New("The nickname must contain letters and numbers")
	}

	return nil
}

func (u *RegisterRequest) validatePass() error {
	if len(u.PasswordHash) < 8 {
		return errors.New("Password is short (Minimum 8 characters)")
	}
	if len(u.PasswordHash) > 64 {
		return errors.New("The password is too long")
	}
	allowedSumbols := "()*_-+=%\""
	var hasLetter, hasDigit bool
	for _, ch := range u.PasswordHash {
		switch {
		case ch >= 'a' && ch <= 'z':
			hasLetter = true
		case ch >= 'A' && ch <= 'Z':
			hasLetter = true
		case unicode.IsDigit(ch):
			hasDigit = true
		case strings.ContainsRune(allowedSumbols, ch):
		default:
			return errors.New(fmt.Sprintf("The password can only contain Latin letters, numbers and symbols %s", allowedSumbols))
		}
	}
	if !hasLetter || !hasDigit {
		return errors.New("The password must contain letters and numbers")
	}

	return nil
}
