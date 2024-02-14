package domain

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Limmperhaven/pkportal-be-v2/internal/errs"
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/tpportal"
	"github.com/jmoiron/sqlx"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"time"
)

func (u *Usecase) ListProfiles(ctx context.Context) ([]tpportal.ListProfilesResponseItem, error) {
	profiles, err := tpportal.Profiles(
		qm.Load(
			qm.Rels(
				tpportal.ProfileRels.Subjects,
			),
		),
	).All(ctx, u.st.DBSX())
	if err != nil {
		return nil, errs.NewInternal(err)
	}

	res := make([]tpportal.ListProfilesResponseItem, len(profiles))
	for i, profile := range profiles {
		subjects := make([]tpportal.IdName, len(profile.R.Subjects))
		for j, subj := range profile.R.Subjects {
			subjects[j] = tpportal.IdName{
				Id:   subj.ID,
				Name: subj.Name,
			}
		}
		res[i] = tpportal.ListProfilesResponseItem{
			Id:            profile.ID,
			Name:          profile.Name,
			EducationYear: int64(profile.EducationYear),
			Subjects:      subjects,
		}
	}
	return res, nil
}

func (u *Usecase) SetProfilesToUser(ctx context.Context, req tpportal.SetProfilesToUserRequest, userId int64, dateCheck bool) error {
	user, err := tpportal.Users(
		tpportal.UserWhere.ID.EQ(userId),
		qm.Load(tpportal.UserRels.UserProfileSubjects),
	).One(ctx, u.st.DBSX())
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.NewNotFound(fmt.Errorf("пользователь с id %d не найден", userId))
		}
		return errs.NewInternal(err)
	}
	if dateCheck {
		td, err := tpportal.UserTestDates(
			tpportal.UserTestDateWhere.UserID.EQ(user.ID),
			tpportal.UserTestDateWhere.EducationYear.EQ(user.EducationYear),
			qm.Load(tpportal.UserTestDateRels.TestDate),
		).One(ctx, u.st.DBSX())
		if err != nil && err != sql.ErrNoRows {
			return errs.NewInternal(err)
		}

		if td != nil && td.R.TestDate != nil {
			if td.R.TestDate.DateTime.Before(time.Now().Add(3 * 24 * time.Hour)) {
				return errs.NewBadRequest(errors.New("профили можно изменять не позднее чем за 3 дня до начала тестирования"))
			}
		}
	}
	if req.FirstProfileId == 0 && req.SecondProfileId != 0 {
		return errs.NewBadRequest(errors.New("невозможно установить второй профиль, не установив первый"))
	}

	if req.FirstProfileId != 0 {
		firstProfile, err := tpportal.FindProfile(ctx, u.st.DBSX(), req.FirstProfileId)
		if err != nil {
			if err == sql.ErrNoRows {
				return errs.NewNotFound(fmt.Errorf("профиль с id %d не найден", req.FirstProfileId))
			}
			return errs.NewInternal(err)
		}

		if firstProfile.EducationYear != user.EducationYear {
			return errs.NewBadRequest(errors.New("профиль 1 не соответствует выбранному году обучения"))
		}
	}

	if req.SecondProfileId != 0 {
		secondProfile, err := tpportal.FindProfile(ctx, u.st.DBSX(), req.SecondProfileId)
		if err != nil {
			if err == sql.ErrNoRows {
				return errs.NewNotFound(fmt.Errorf("профиль с id %d не найден", req.SecondProfileId))
			}
			return errs.NewInternal(err)
		}

		if secondProfile.EducationYear != user.EducationYear {
			return errs.NewBadRequest(errors.New("профиль 2 не соответствует выбранному году обучения"))
		}
	}

	up := tpportal.UserProfile{
		UserID:            user.ID,
		UserEducationYear: user.EducationYear,
	}

	if req.FirstProfileId == 0 {
		up.FirstProfileID = null.Int64FromPtr(nil)
	} else {
		up.FirstProfileID = null.Int64From(req.FirstProfileId)
	}

	if req.SecondProfileId == 0 {
		up.SecondProfileID = null.Int64FromPtr(nil)
	} else {
		up.SecondProfileID = null.Int64From(req.SecondProfileId)
	}

	var userSubjects *tpportal.UserProfileSubject
	for _, us := range user.R.UserProfileSubjects {
		if us.UserEducationYear == user.EducationYear {
			userSubjects = us
			break
		}
	}

	err = u.st.QueryTx(ctx, func(tx *sqlx.Tx) error {
		if userSubjects != nil {
			_, err = userSubjects.Delete(ctx, tx)
			if err != nil {
				return errs.NewInternal(err)
			}
		}
		err = up.Upsert(ctx, tx, true,
			[]string{tpportal.UserProfileColumns.UserID, tpportal.UserProfileColumns.UserEducationYear},
			boil.Whitelist(tpportal.UserProfileColumns.FirstProfileID, tpportal.UserProfileColumns.SecondProfileID),
			boil.Infer())
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
