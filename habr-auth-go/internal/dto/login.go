package dto

import "errors"

type Login struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	FingerPrint string `json:"fingerprint"`
}

func (u *Login) IsValid() error {
	if u.Email == "" || u.Password == "" || u.FingerPrint == "" {
		return errors.New("Missing required fields")
	}
	return nil
}
