package domain

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Limmperhaven/pkportal-be-v2/internal/body"
	"github.com/Limmperhaven/pkportal-be-v2/internal/errs"
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/tpportal"
	"github.com/friendsofgo/errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func (u *Usecase) SignUp(ctx context.Context, req *tpportal.SignUpRequest) error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password+body.AppSalt), body.AppCost)
	if err != nil {
		return errs.NewInternal(err)
	}

	dob, err := time.Parse("02.01.2006", req.DateOfBirth)
	if err != nil {
		return errs.NewBadRequest(fmt.Errorf("invalid date_of_birth: %s", req.DateOfBirth))
	}

	user := tpportal.User{
		Email:             req.Email,
		HashPassword:      string(hashPassword),
		Fio:               req.Fio,
		DateOfBirth:       dob,
		Gender:            tpportal.UserGender(req.Gender),
		PhoneNumber:       req.PhoneNumber,
		ParentPhoneNumber: req.ParentPhoneNumber,
		CurrentSchool:     null.StringFrom(req.CurrentSchool),
		EducationYear:     int16(req.EducationYear),
	}
	err = user.Insert(ctx, u.st.DBSX(), boil.Infer())
	if err != nil {
		return errs.NewInternal(err)
	}
	return nil
}

func (u *Usecase) SignIn(ctx context.Context, req *tpportal.SignInRequest) (tpportal.UserWithAuth, error) {
	user, err := tpportal.Users(
		tpportal.UserWhere.Email.EQ(req.Email),
	).One(ctx, u.st.DBSX())
	if err != nil {
		if err == sql.ErrNoRows {
			return tpportal.UserWithAuth{}, errs.NewUnauthorized(errors.New("Пользователь с таким email не найден"))
		}
		return tpportal.UserWithAuth{}, errs.NewInternal(err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(req.Password+body.AppSalt))
	if err != nil {
		return tpportal.UserWithAuth{}, errs.NewUnauthorized(errors.New("Введен неверный пароль"))
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tpportal.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Id: user.ID,
	})
	signedToken, err := token.SignedString([]byte(body.AppSalt))
	if err != nil {
		return tpportal.UserWithAuth{}, errs.NewInternal(err)
	}

	return tpportal.UserWithAuth{
		User:      *user,
		AuthToken: signedToken,
	}, nil
}
