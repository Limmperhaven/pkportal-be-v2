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
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func (u *Usecase) SignUp(ctx context.Context, req *tpportal.SignUpRequest) error {
	checkUser, err := tpportal.Users(tpportal.UserWhere.Email.EQ(req.Email)).One(ctx, u.st.DBSX())
	if err != nil && err != sql.ErrNoRows {
		return errs.NewInternal(err)
	}
	if checkUser != nil {
		return errs.NewBadRequest(errors.New("пользователь с таким email уже зарегистрирован"))
	}

	hashPassword, err := u.hashPassword(req.Password)
	if err != nil {
		return err
	}
	dob, err := u.parseDate(req.DateOfBirth)
	if err != nil {
		return err
	}

	user := tpportal.User{
		Email:                      req.Email,
		HashPassword:               hashPassword,
		Fio:                        req.Fio,
		DateOfBirth:                dob,
		Gender:                     tpportal.UserGender(req.Gender),
		PhoneNumber:                req.PhoneNumber,
		ParentPhoneNumber:          req.ParentPhoneNumber,
		CurrentSchool:              null.StringFrom(req.CurrentSchool),
		EducationYear:              int16(req.EducationYear),
		IsActivated:                false,
		ActivationToken:            uuid.New().String(),
		ChangePasswordToken:        uuid.New().String(),
		LastActivationMailSent:     null.Time{Valid: false},
		LastChangePasswordMailSent: null.Time{Valid: false},
	}

	var otherEducationYear int16
	if user.EducationYear == int16(10) {
		otherEducationYear = int16(9)
	} else {
		otherEducationYear = int16(10)
	}

	cfg := config.Get().Server
	activationLink := cfg.Domain + "/auth/activate/" + user.ActivationToken

	err = u.st.QueryTx(ctx, func(tx *sqlx.Tx) error {
		err = user.Insert(ctx, tx, boil.Infer())
		if err != nil {
			return errs.NewInternal(err)
		}
		uss := tpportal.UserStatusSlice{
			&tpportal.UserStatus{
				UserID:        user.ID,
				StatusID:      body.Registered.Int64(),
				EducationYear: user.EducationYear,
			},
			&tpportal.UserStatus{
				UserID:        user.ID,
				StatusID:      body.Registered.Int64(),
				EducationYear: otherEducationYear,
			},
		}
		err = user.AddUserStatuses(ctx, tx, true, uss...)
		if err != nil {
			return errs.NewInternal(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = u.mail.SendTextEmail(body.CreateAccountSubject, body.CreateAccountMessage+activationLink, []string{req.Email})
	if err != nil {
		return errs.NewInternal(err)
	}

	return nil
}

func (u *Usecase) SignIn(ctx context.Context, req *tpportal.SignInRequest) (tpportal.SignInResponse, error) {
	user, err := tpportal.Users(
		tpportal.UserWhere.Email.EQ(req.Email),
		qm.Load(
			qm.Rels(
				tpportal.UserRels.UserStatuses,
				tpportal.UserStatusRels.Status,
			),
		),
	).One(ctx, u.st.DBSX())
	if err != nil {
		if err == sql.ErrNoRows {
			return tpportal.SignInResponse{}, errs.NewUnauthorized(errors.New("Пользователь с таким email не найден"))
		}
		return tpportal.SignInResponse{}, errs.NewInternal(err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(req.Password+body.AppSalt))
	if err != nil {
		return tpportal.SignInResponse{}, errs.NewUnauthorized(errors.New("Введен неверный пароль"))
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
		return tpportal.SignInResponse{}, errs.NewInternal(err)
	}

	status := tpportal.IdName{}
	if len(user.R.UserStatuses) != 0 {
		for _, us := range user.R.UserStatuses {
			if us.EducationYear == user.EducationYear {
				status.Id = us.R.Status.ID
				status.Name = us.R.Status.Name
			}
		}
	}

	return tpportal.SignInResponse{
		Id:                user.ID,
		Email:             user.Email,
		Fio:               user.Fio,
		DateOfBirth:       u.formatDate(user.DateOfBirth),
		Gender:            user.Gender.String(),
		PhoneNumber:       user.PhoneNumber,
		ParentPhoneNumber: user.ParentPhoneNumber,
		CurrentSchool:     user.CurrentSchool.String,
		EducationYear:     int64(user.EducationYear),
		IsActivated:       user.IsActivated,
		Role:              user.Role.String(),
		Status:            status,
		AuthToken:         signedToken,
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
	if user.LastChangePasswordMailSent.Time.Add(2 * time.Minute).After(time.Now()) {
		return errs.NewBadRequest(errors.New("письмо можно отправлять не чаще чем раз в 2 минуты"))
	}

	cfg := config.Get().App
	url := cfg.FrontendUrl + "/setPassword/" + user.ChangePasswordToken
	user.LastChangePasswordMailSent = null.TimeFrom(time.Now())
	_, err = user.Update(ctx, u.st.DBSX(), boil.Whitelist(tpportal.UserColumns.LastChangePasswordMailSent))
	if err != nil {
		return errs.NewInternal(err)
	}

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
