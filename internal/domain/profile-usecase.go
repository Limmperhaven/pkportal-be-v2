package domain

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Limmperhaven/pkportal-be-v2/internal/errs"
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/tpportal"
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
	user, err := tpportal.FindUser(ctx, u.st.DBSX(), userId)
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

	up := tpportal.UserProfile{
		UserID:            user.ID,
		UserEducationYear: user.EducationYear,
		FirstProfileID:    null.Int64From(req.FirstProfileId),
		SecondProfileID:   null.Int64From(req.SecondProfileId),
	}

	err = up.Upsert(ctx, u.st.DBSX(), true,
		[]string{tpportal.UserProfileColumns.UserID, tpportal.UserProfileColumns.UserEducationYear},
		boil.Whitelist(tpportal.UserProfileColumns.FirstProfileID, tpportal.UserProfileColumns.SecondProfileID),
		boil.Infer())
	if err != nil {
		return errs.NewInternal(err)
	}
	return nil
}
