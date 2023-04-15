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

func (u *Usecase) ListForeignLanguages(ctx context.Context) ([]tpportal.IdName, error) {
	fls, err := tpportal.ForeignLanguages().All(ctx, u.st.DBSX())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFound(errors.New("иностранные языки не найдены"))
		}
		return nil, errs.NewInternal(err)
	}
	res := make([]tpportal.IdName, len(fls))
	for i, fl := range fls {
		res[i] = tpportal.IdName{
			Id:   fl.ID,
			Name: fl.Name,
		}
	}
	return res, nil
}

func (u *Usecase) SetForeignLanguageToUser(ctx context.Context, userId, flId int64, dateCheck bool) error {
	user, err := tpportal.FindUser(ctx, u.st.DBSX(), userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.NewNotFound(fmt.Errorf("пользователь с id: %d не найден", userId))
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
				return errs.NewBadRequest(errors.New("иностранные языки можно изменять не позднее чем за 3 дня до начала тестирования"))
			}
		}
	}
	ufl := tpportal.UserForeignLanguage{
		UserID:            user.ID,
		UserEducationYear: user.EducationYear,
		ForeignLanguageID: null.Int64From(flId),
	}
	err = ufl.Upsert(ctx, u.st.DBSX(), true,
		[]string{tpportal.UserForeignLanguageColumns.UserID, tpportal.UserForeignLanguageColumns.UserEducationYear},
		boil.Whitelist(tpportal.UserForeignLanguageColumns.ForeignLanguageID), boil.Infer())
	if err != nil {
		return errs.NewInternal(err)
	}
	return nil
}
