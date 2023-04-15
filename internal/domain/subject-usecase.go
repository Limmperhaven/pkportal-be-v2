package domain

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Limmperhaven/pkportal-be-v2/internal/errs"
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/tpportal"
	"github.com/friendsofgo/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"time"
)

func (u *Usecase) ListSubjects(ctx context.Context, profileId int64) ([]tpportal.IdName, error) {
	if profileId == 0 {
		subjects, err := tpportal.Subjects().All(ctx, u.st.DBSX())
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, errs.NewNotFound(errors.New("предметы не найдены"))
			}
			return nil, errs.NewInternal(err)
		}
		res := make([]tpportal.IdName, len(subjects))
		for i, subj := range subjects {
			res[i] = tpportal.IdName{
				Id:   subj.ID,
				Name: subj.Name,
			}
		}
		return res, nil
	}

	profile, err := tpportal.Profiles(
		tpportal.ProfileWhere.ID.EQ(profileId),
		qm.Load(tpportal.ProfileRels.Subjects),
	).One(ctx, u.st.DBSX())
	tpportal.Subjects()
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFound(fmt.Errorf("профиль с id: %d не найден", profileId))
		}
		return nil, errs.NewInternal(err)
	}

	res := make([]tpportal.IdName, len(profile.R.Subjects))
	for i, subj := range profile.R.Subjects {
		res[i] = tpportal.IdName{
			Id:   subj.ID,
			Name: subj.Name,
		}
	}
	return res, nil
}

func (u *Usecase) SetSubjectsToUser(ctx context.Context, req tpportal.SetSubjectsRequest, userId int64, dateCheck bool) error {
	user, err := tpportal.Users(
		tpportal.UserWhere.ID.EQ(userId),
		qm.Load(
			qm.Rels(
				tpportal.UserRels.UserProfiles,
				tpportal.UserProfileRels.FirstProfile,
				tpportal.ProfileRels.Subjects,
			),
		),
		qm.Load(
			qm.Rels(
				tpportal.UserRels.UserProfiles,
				tpportal.UserProfileRels.SecondProfile,
				tpportal.ProfileRels.Subjects,
			),
		),
	).One(ctx, u.st.DBSX())
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.NewNotFound(fmt.Errorf("пользователь с id: %d не найден", userId))
		}
		return errs.NewInternal(err)
	}

	if req.FirstSubjectId == 0 && req.SecondSubjectId != 0 {
		return errs.NewBadRequest(errors.New("невозможно выбрать предмет второго профиля, не выбрав предмет первого профиля"))
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
				return errs.NewBadRequest(errors.New("профильные предметы можно изменять не позднее чем за 3 дня до начала тестирования"))
			}
		}
	}
	ups := tpportal.UserProfileSubject{
		UserID:            user.ID,
		UserEducationYear: user.EducationYear,
	}
	var firstUserProfile *tpportal.Profile
	var secondUserProfile *tpportal.Profile
	for _, up := range user.R.UserProfiles {
		if up.UserEducationYear == user.EducationYear {
			firstUserProfile = up.R.FirstProfile
			secondUserProfile = up.R.SecondProfile
			break
		}
	}
	if req.FirstSubjectId == 0 {
		ups.FirstProfileSubjectID = null.Int64{Valid: false}
	} else {
		if firstUserProfile == nil {
			return errs.NewBadRequest(errors.New("для выбора предмета заполните 1 профиль"))
		}

		valid := false
		for _, subj := range firstUserProfile.R.Subjects {
			if subj.ID == req.FirstSubjectId {
				valid = true
				break
			}
		}
		if !valid {
			return errs.NewBadRequest(errors.New("выбранный предмет 1 профиля не является предметом выбранного профиля"))
		}
		ups.FirstProfileSubjectID = null.Int64From(req.FirstSubjectId)
	}
	if req.SecondSubjectId == 0 {
		ups.SecondProfileSubjectID = null.Int64{Valid: false}
	} else {
		if secondUserProfile == nil {
			return errs.NewBadRequest(errors.New("для выбора предмета заполните 2 профиль"))
		}
		valid := false
		for _, subj := range secondUserProfile.R.Subjects {
			if subj.ID == req.SecondSubjectId {
				valid = true
				break
			}
		}
		if !valid {
			return errs.NewBadRequest(errors.New("выбранный предмет 2 профиля не является предметом выбранного профиля"))
		}
		ups.SecondProfileSubjectID = null.Int64From(req.SecondSubjectId)
	}
	err = ups.Upsert(ctx, u.st.DBSX(), true,
		[]string{tpportal.UserProfileSubjectColumns.UserID, tpportal.UserProfileSubjectColumns.UserEducationYear},
		boil.Whitelist(tpportal.UserProfileSubjectColumns.FirstProfileSubjectID, tpportal.UserProfileSubjectColumns.SecondProfileSubjectID),
		boil.Infer())
	if err != nil {
		return errs.NewInternal(err)
	}

	return nil
}
