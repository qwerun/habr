package dto

import "errors"

type Refresh struct {
	Access      string `json:"access"`
	Refresh     string `json:"refresh"`
	FingerPrint string `json:"fingerprint"`
}

func (u *Refresh) IsValid() error {
	if u.Access == "" || u.Refresh == "" || u.FingerPrint == "" {
		return errors.New("Missing required fields")
	}
	return nil
}
