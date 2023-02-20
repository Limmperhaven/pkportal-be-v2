package domain

import (
	"context"
	"database/sql"
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
	hashPassword, err := u.hashPassword(req.Password)
	if err != nil {
		return err
	}
	dob, err := u.parseDate(req.DateOfBirth)
	if err != nil {
		return err
	}

	user := tpportal.User{
		Email:               req.Email,
		HashPassword:        hashPassword,
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
		StatusID:            body.Registered.Int64(),
	}
	cfg := config.Get().Server
	activationLink := cfg.Scheme + "://" + cfg.Domain + "/auth/activate/" + user.ActivationToken

	err = user.Insert(ctx, u.st.DBSX(), boil.Infer())
	if err != nil {
		return errs.NewInternal(err)
	}

	err = u.mail.SendTextEmail(body.CreateAccountSubject, body.CreateAccountMessage+activationLink, []string{req.Email})
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

	cfg := config.Get().Server
	url := cfg.Domain + "/someUrl/" + user.ChangePasswordToken

	err = u.mail.SendTextEmail(body.RecoverPasswordSubject, body.RecoverPasswordMessage+url, []string{user.Email})
	if err != nil {
		return errs.NewInternal(err)
	}

	return nil
}

func (u *Usecase) ConfirmRecover(ctx context.Context, token, newPassword string) error {
	user, err := tpportal.Users(tpportal.UserWhere.ChangePasswordToken.EQ(token)).One(ctx, u.st.DBSX())
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.NewNotFound(errors.New("пользователь с таким email не найден"))
		}
		return errs.NewInternal(err)
	}

	hashPassword, err := u.hashPassword(newPassword)
	if err != nil {
		return err
	}
	user.HashPassword = hashPassword
	user.ChangePasswordToken = uuid.New().String()
	_, err = user.Update(ctx, u.st.DBSX(), boil.Infer())
	if err != nil {
		return errs.NewInternal(errors.New("ошибка при обновлении пользователя"))
	}
	return nil
}
