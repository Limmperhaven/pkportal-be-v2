package domain

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Limmperhaven/pkportal-be-v2/internal/body"
	"github.com/Limmperhaven/pkportal-be-v2/internal/config"
	"github.com/Limmperhaven/pkportal-be-v2/internal/errs"
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/tpportal"
	"github.com/google/uuid"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (u *Usecase) CreateUser(ctx context.Context, req *tpportal.CreateUserRequest) error {
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
		IsActivated:         req.IsActivated,
		ActivationToken:     uuid.New().String(),
		ChangePasswordToken: uuid.New().String(),
		Role:                tpportal.UserRole(req.Role),
		StatusID:            req.StatusId,
	}

	err = user.Insert(ctx, u.st.DBSX(), boil.Infer())
	if err != nil {
		return errs.NewInternal(err)
	}

	return nil
}

func (u *Usecase) GetUser(ctx context.Context, userId int64) (tpportal.GetUserResponse, error) {
	user, err := tpportal.Users(
		tpportal.UserWhere.ID.EQ(userId),
		qm.Load(
			qm.Rels(
				tpportal.UserRels.Status,
			),
		),
		qm.Load(
			qm.Rels(
				tpportal.UserRels.UserProfiles,
				tpportal.UserProfileRels.FirstProfile,
			),
		),
		qm.Load(
			qm.Rels(
				tpportal.UserRels.UserProfiles,
				tpportal.UserProfileRels.SecondProfile,
			),
		),
		qm.Load(
			qm.Rels(
				tpportal.UserRels.UserProfileSubjects,
				tpportal.UserProfileSubjectRels.FirstProfileSubject,
			),
		),
		qm.Load(
			qm.Rels(
				tpportal.UserRels.UserProfileSubjects,
				tpportal.UserProfileSubjectRels.SecondProfileSubject,
			),
		),
		qm.Load(
			qm.Rels(
				tpportal.UserRels.UserForeignLanguages,
				tpportal.UserForeignLanguageRels.ForeignLanguage,
			),
		),
	).One(ctx, u.st.DBSX())
	if err != nil {
		if err == sql.ErrNoRows {
			return tpportal.GetUserResponse{}, errs.NewNotFound(fmt.Errorf("пользователь с id %d не найден", userId))
		}
		return tpportal.GetUserResponse{}, errs.NewInternal(err)
	}

	firstProfile := tpportal.IdName{}
	secondProfile := tpportal.IdName{}
	if len(user.R.UserProfiles) != 0 {
		for _, up := range user.R.UserProfiles {
			if up.UserEducationYear == user.EducationYear {
				firstProfile.Id = up.R.FirstProfile.ID
				firstProfile.Name = up.R.FirstProfile.Name
				secondProfile.Id = up.R.SecondProfile.ID
				secondProfile.Name = up.R.SecondProfile.Name
				break
			}
		}
	}

	firstProfileSubject := tpportal.IdName{}
	secondProfileSubject := tpportal.IdName{}
	if len(user.R.UserProfileSubjects) != 0 {
		for _, ups := range user.R.UserProfileSubjects {
			if ups.UserEducationYear == user.EducationYear {
				firstProfileSubject.Id = ups.R.FirstProfileSubject.ID
				firstProfileSubject.Name = ups.R.FirstProfileSubject.Name
				secondProfileSubject.Id = ups.R.SecondProfileSubject.ID
				secondProfileSubject.Name = ups.R.SecondProfileSubject.Name
			}
		}
	}
	foreignLanguage := tpportal.IdName{}
	if len(user.R.UserForeignLanguages) != 0 {
		for _, fl := range user.R.UserForeignLanguages {
			if fl.UserEducationYear == user.EducationYear {
				foreignLanguage.Id = fl.R.ForeignLanguage.ID
				foreignLanguage.Name = fl.R.ForeignLanguage.Name
			}
		}
	}

	res := tpportal.GetUserResponse{
		Id:                user.ID,
		Role:              user.Role.String(),
		Fio:               user.Fio,
		DateOfBirth:       u.formatDate(user.DateOfBirth),
		Gender:            user.Gender.String(),
		Email:             user.Email,
		PhoneNumber:       user.PhoneNumber,
		ParentPhoneNumber: user.ParentPhoneNumber,
		CurrentSchool:     user.CurrentSchool.String,
		EducationYear:     int64(user.EducationYear),
		Status: tpportal.IdName{
			Id:   user.R.Status.ID,
			Name: user.R.Status.Name,
		},
		FirstProfile:         firstProfile,
		SecondProfile:        secondProfile,
		FirstProfileSubject:  firstProfileSubject,
		SecondProfileSubject: secondProfileSubject,
		ForeignLanguage:      foreignLanguage,
		IsActivated:          user.IsActivated,
	}

	return res, nil
}

func (u *Usecase) UpdateUser(ctx context.Context, req tpportal.UpdateUserRequest, userId int64) error {
	user, err := tpportal.Users(tpportal.UserWhere.ID.EQ(userId)).One(ctx, u.st.DBSX())
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.NewNotFound(fmt.Errorf("пользователь с id: %d не найден", userId))
		}
		return errs.NewInternal(err)
	}
	if req.Email != user.Email {
		user.Email = req.Email
		user.IsActivated = false

		cfg := config.Get().Server
		activationLink := cfg.Scheme + "://" + cfg.Domain + "/auth/activate/" + user.ActivationToken

		err = u.mail.SendTextEmail(body.CreateAccountSubject, body.CreateAccountMessage+activationLink, []string{req.Email})
		if err != nil {
			return errs.NewInternal(err)
		}
	}
	if req.DateOfBirth != "" {
		dob, err := u.parseDate(req.DateOfBirth)
		if err != nil {
			return err
		}
		user.DateOfBirth = dob
	}

	user.Fio = req.Fio
	user.Gender = tpportal.UserGender(req.Gender)
	user.PhoneNumber = req.PhoneNumber
	user.ParentPhoneNumber = req.ParentPhoneNumber
	user.CurrentSchool = null.StringFrom(req.CurrentSchool)
	user.EducationYear = int16(req.EducationYear)

	_, err = user.Update(ctx, u.st.DBSX(), boil.Infer())
	if err != nil {
		return errs.NewInternal(err)
	}

	return nil
}

func (u *Usecase) ListStatuses(ctx context.Context, request tpportal.ListStatusesRequest) ([]tpportal.IdName, error) {
	conditions := make([]qm.QueryMod, 0, 2)
	if request.AvailableFor10thClass {
		conditions = append(conditions, tpportal.StatusWhere.AvailableFor10THClass.EQ(true))
	}
	if request.AvailableFor9thClass {
		conditions = append(conditions, tpportal.StatusWhere.AvailableFor9THClass.EQ(true))
	}

	statuses, err := tpportal.Statuses(conditions...).All(ctx, u.st.DBSX())
	if err != nil {
		return nil, errs.NewInternal(err)
	}

	res := make([]tpportal.IdName, len(statuses))
	for i, status := range statuses {
		res[i] = tpportal.IdName{Id: status.ID, Name: status.Name}
	}

	return res, nil
}

func (u *Usecase) SetUserStatus(ctx context.Context, userId int64, statusId int64) error {
	user := tpportal.User{ID: userId, StatusID: statusId}
	_, err := user.Update(ctx, u.st.DBSX(), boil.Whitelist(tpportal.UserColumns.StatusID))
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.NewNotFound(fmt.Errorf("пользователь с id: %d не найден", userId))
		}
		return errs.NewInternal(err)
	}
	return nil
}
