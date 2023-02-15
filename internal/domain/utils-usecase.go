package domain

import (
	"fmt"
	"github.com/Limmperhaven/pkportal-be-v2/internal/body"
	"github.com/Limmperhaven/pkportal-be-v2/internal/errs"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func (u *Usecase) parseDate(in string) (time.Time, error) {
	dob, err := time.Parse("02.01.2006", in)
	if err != nil {
		return time.Time{}, errs.NewBadRequest(fmt.Errorf("invalid date_of_birth: %s", in))
	}
	return dob, nil
}

func (u *Usecase) hashPassword(in string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(in+body.AppSalt), body.AppCost)
	if err != nil {
		return "", errs.NewInternal(err)
	}
	return string(hash), nil
}
