package domain

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Limmperhaven/pkportal-be-v2/internal/body"
	"github.com/Limmperhaven/pkportal-be-v2/internal/errs"
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/tpportal"
	"github.com/friendsofgo/errors"
	"github.com/jmoiron/sqlx"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"time"
)

func (u *Usecase) CreateTestDate(ctx context.Context, req tpportal.CreateTestDateRequest) error {
	dateTime, err := u.parseDateTime(req.Date, req.Time)
	if err != nil {
		return err
	}

	testDate := tpportal.TestDate{
		DateTime:         dateTime,
		Location:         req.Location,
		MaxPersons:       int(req.MaxPersons),
		EducationYear:    int16(req.EducationYear),
		PubStatus:        tpportal.TestDatePubStatus(req.PubStatus),
		NotificationSent: false,
	}
	err = testDate.Insert(ctx, u.st.DBSX(), boil.Infer())
	if err != nil {
		return errs.NewInternal(err)
	}
	return nil
}

func (u *Usecase) SetTestDatePubStatus(ctx context.Context, tdId int64, status string) error {
	td := tpportal.TestDate{ID: tdId, PubStatus: tpportal.TestDatePubStatus(status)}
	_, err := td.Update(ctx, u.st.DBSX(), boil.Whitelist(tpportal.TestDateColumns.PubStatus))
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.NewNotFound(fmt.Errorf("дата тестирования с id: %d не найдена", tdId))
		}
		return errs.NewInternal(err)
	}
	return nil
}

func (u *Usecase) ListTestDates(ctx context.Context, availableOnly bool) ([]tpportal.ListTestDatesResponseItem, error) {
	user, err := u.extractUserFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	queryMods := make([]qm.QueryMod, 0, 3)
	queryMods = append(queryMods, qm.Load(
		qm.Rels(tpportal.TestDateRels.UserTestDates),
	))
	if availableOnly {
		queryMods = append(queryMods, tpportal.TestDateWhere.PubStatus.EQ(tpportal.TestDatePubStatusShown))
		queryMods = append(queryMods, tpportal.TestDateWhere.DateTime.GT(time.Now().Add(3*24*time.Hour)))
		utd, err := tpportal.UserTestDates(
			tpportal.UserTestDateWhere.EducationYear.EQ(user.EducationYear),
			tpportal.UserTestDateWhere.UserID.EQ(user.ID),
		).One(ctx, u.st.DBSX())
		if err != nil && err != sql.ErrNoRows {
			return nil, errs.NewInternal(err)
		}
		queryMods = append(queryMods, tpportal.TestDateWhere.ID.NEQ(utd.TestDateID))
	}

	tds, err := tpportal.TestDates(queryMods...).All(ctx, u.st.DBSX())
	if err != nil {
		return nil, errs.NewInternal(err)
	}

	res := make([]tpportal.ListTestDatesResponseItem, 0, len(tds))
	for _, td := range tds {
		if availableOnly && (td.MaxPersons == len(td.R.UserTestDates) || td.EducationYear != user.EducationYear) {
			continue
		}
		date, time := u.formatDateTime(td.DateTime)

		res = append(res, tpportal.ListTestDatesResponseItem{
			Id:                td.ID,
			Date:              date,
			Time:              time,
			Location:          td.Location,
			RegisteredPersons: int64(len(td.R.UserTestDates)),
			MaxPersons:        int64(td.MaxPersons),
			EducationYear:     int64(td.EducationYear),
			PubStatus:         td.PubStatus.String(),
		})
	}
	return res, nil
}

func (u *Usecase) SignUpUserToTestDate(ctx context.Context, userId, tdId int64, dateCheck bool) error {
	user, err := tpportal.Users(
		tpportal.UserWhere.ID.EQ(userId),
		qm.Load(
			qm.Rels(
				tpportal.UserRels.UserStatuses,
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
	).One(ctx, u.st.DBSX())
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.NewNotFound(fmt.Errorf("пользователь с id: %d не найден", userId))
		}
		return errs.NewInternal(err)
	}

	if dateCheck {
		prevTd, err := tpportal.UserTestDates(
			tpportal.UserTestDateWhere.UserID.EQ(user.ID),
			tpportal.UserTestDateWhere.EducationYear.EQ(user.EducationYear),
			qm.Load(tpportal.UserTestDateRels.TestDate),
		).One(ctx, u.st.DBSX())
		if err != nil && err != sql.ErrNoRows {
			return errs.NewInternal(err)
		}

		if prevTd != nil && prevTd.R.TestDate != nil {
			if prevTd.R.TestDate.DateTime.Before(time.Now().Add(3 * 24 * time.Hour)) {
				return errs.NewBadRequest(errors.New("дату тестирования можно изменять не позднее чем за 3 дня до начала тестирования"))
			}
		}
	}

	td, err := tpportal.FindTestDate(ctx, u.st.DBSX(), tdId)
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.NewNotFound(fmt.Errorf("дата тестирования с id: %d не найдена", tdId))
		}
		return errs.NewInternal(err)
	}

	if td.EducationYear != user.EducationYear {
		return errs.NewBadRequest(fmt.Errorf("данная дата тестирования доступна только для %d класса", td.EducationYear))
	}

	regCount, err := tpportal.UserTestDates(tpportal.UserTestDateWhere.TestDateID.EQ(td.ID)).Count(ctx, u.st.DBSX())
	if err != nil {
		return errs.NewInternal(err)
	}

	if regCount == int64(td.MaxPersons) {
		return errs.NewBadRequest(errors.New("недостаточно мест"))
	}

	utd := tpportal.UserTestDate{
		UserID:        user.ID,
		TestDateID:    td.ID,
		EducationYear: user.EducationYear,
		IsAttended:    false,
	}
	for i := range user.R.UserStatuses {
		if user.R.UserStatuses[i].EducationYear == user.EducationYear {
			user.R.UserStatuses[i].StatusID = body.Registered.Int64()
		}
	}

	err = u.st.QueryTx(ctx, func(tx *sqlx.Tx) error {
		err := utd.Upsert(ctx, tx, true,
			[]string{tpportal.UserTestDateColumns.UserID, tpportal.UserTestDateColumns.EducationYear},
			boil.Whitelist(tpportal.UserTestDateColumns.TestDateID, tpportal.UserTestDateColumns.IsAttended),
			boil.Infer())
		if err != nil {
			return errs.NewInternal(err)
		}
		for _, us := range user.R.UserStatuses {
			if us.EducationYear == user.EducationYear {
				us.StatusID = body.Registered.Int64()
				_, err = us.Update(ctx, tx, boil.Whitelist(tpportal.UserStatusColumns.StatusID))
				if err != nil {
					return errs.NewInternal(err)
				}
				return nil
			}
		}
		return errs.NewNotFound(errors.New("у пользователя не хватает записи о статусе"))
	})
	if err != nil {
		return err
	}

	var userProfilesString string
	if user.R.UserProfiles != nil {
		for _, up := range user.R.UserProfiles {
			if up.UserEducationYear == user.EducationYear {
				if up.R.FirstProfile != nil {
					userProfilesString = up.R.FirstProfile.Name
				}
				if up.R.SecondProfile != nil {
					userProfilesString += ", " + up.R.SecondProfile.Name
				}
				break
			}
		}
	}

	var userProfileSubjectsString string
	if user.R.UserProfileSubjects != nil {
		for _, ups := range user.R.UserProfileSubjects {
			if ups.UserEducationYear == user.EducationYear {
				if ups.R.FirstProfileSubject != nil {
					userProfileSubjectsString = ups.R.FirstProfileSubject.Name
				}
				if ups.R.SecondProfileSubject != nil {
					userProfileSubjectsString += ", " + ups.R.SecondProfileSubject.Name
				}
				break
			}
		}
	}

	tdDate, tdTime := u.formatDateTime(td.DateTime)
	emailMessage := fmt.Sprintf(body.SignUpForTestDateMessage, tdDate, td.Location,
		tdTime, userProfilesString, userProfileSubjectsString)

	err = u.mail.SendTextEmail(body.SignUpForTestDateSubject, emailMessage, []string{user.Email})
	if err != nil {
		return errs.NewInternal(err)
	}

	return nil
}

func (u *Usecase) ListCommonLocations(ctx context.Context) ([]tpportal.IdName, error) {
	cls, err := tpportal.CommonLocations().All(ctx, u.st.DBSX())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errs.NewInternal(err)
	}
	res := make([]tpportal.IdName, len(cls))
	for i, cl := range cls {
		res[i] = tpportal.IdName{
			Id:   cl.ID,
			Name: cl.Name,
		}
	}
	return res, nil
}

func (u *Usecase) SetTestDateAttended(ctx context.Context, userId, tdId int64) error {
	user, err := tpportal.FindUser(ctx, u.st.DBSX(), userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.NewNotFound(fmt.Errorf("пользователь с id: %d не найден", userId))
		}
		return errs.NewInternal(err)
	}

	utd, err := tpportal.FindUserTestDate(ctx, u.st.DBSX(), user.ID, user.EducationYear)
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.NewNotFound(errors.New("указанная запись на тестирование не найдена"))
		}
		return errs.NewInternal(err)
	}
	if utd.TestDateID != tdId {
		return errs.NewNotFound(errors.New("указанная запись на тестирование не найдена"))
	}
	utd.IsAttended = true
	_, err = utd.Update(ctx, u.st.DBSX(), boil.Whitelist(tpportal.UserTestDateColumns.IsAttended))
	if err != nil {
		return errs.NewInternal(err)
	}
	return nil
}
