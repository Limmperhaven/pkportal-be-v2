package domain

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Limmperhaven/pkportal-be-v2/internal/body"
	"github.com/Limmperhaven/pkportal-be-v2/internal/config"
	"github.com/Limmperhaven/pkportal-be-v2/internal/errs"
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/tpportal"
	"github.com/friendsofgo/errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
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
		Email:               req.Email,
		HashPassword:        string(hashPassword),
		Fio:                 req.Fio,
		DateOfBirth:         dob,
		Gender:              tpportal.UserGender(req.Gender),
		PhoneNumber:         req.PhoneNumber,
		ParentPhoneNumber:   req.ParentPhoneNumber,
		CurrentSchool:       null.StringFrom(req.CurrentSchool),
		EducationYear:       int16(req.EducationYear),
		IsActivated:         false,
		ActivationToken:     uuid.New().String(),
		ChangePasswordToken: uuid.New().String(),
	}
	cfg := config.Get().Server
	activationLink := cfg.Scheme + "://" + cfg.Host + "/auth/activate/" + user.ActivationToken

	err = u.st.QueryTx(ctx, func(tx *sqlx.Tx) error {
		err = user.Insert(ctx, tx, boil.Infer())
		if err != nil {
			return errs.NewInternal(err)
		}

		err = u.mail.SendTextEmail(body.CreateAccountSubject, body.CreateAccountMessage+activationLink, []string{req.Email})
		if err != nil {
			return errs.NewInternal(err)
		}
		return nil
	})
	if err != nil {
		return err
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

func (u *Usecase) Activate(ctx context.Context, token string) error {
	err := u.st.QueryTx(ctx, func(tx *sqlx.Tx) error {
		user, err := tpportal.Users(tpportal.UserWhere.ActivationToken.EQ(token)).One(ctx, tx)
		if err != nil {
			if err == sql.ErrNoRows {
				return errs.NewNotFound(err)
			}
			return errs.NewInternal(err)
		}

		user.IsActivated = true
		user.ActivationToken = uuid.New().String()
		_, err = user.Update(ctx, tx, boil.Infer())
		if err != nil {
			return errs.NewInternal(err)
		}
		return nil
	})
	return err
}

func (u *Usecase) RecoverPassword(ctx context.Context, email string) error {
	user, err := tpportal.Users(tpportal.UserWhere.Email.EQ(email)).One(ctx, u.st.DBSX())
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.NewNotFound(errors.New("пользователь с таким email не найден"))
		}
		return errs.NewInternal(err)
	}

	err = u.mail.SendTextEmail(body.RecoverPasswordSubject, body.RecoverPasswordMessage, []string{user.Email})
	if err != nil {
		return errs.NewInternal(err)
	}

	return nil
}
