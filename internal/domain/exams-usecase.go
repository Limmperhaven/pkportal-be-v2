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
)

func (u *Usecase) SetGrades(ctx context.Context, req tpportal.SetGradesRequest) error {
	user, err := tpportal.FindUser(ctx, u.st.DBSX(), req.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.NewNotFound(fmt.Errorf("пользователь с id %d не найден", req.UserId))
		}
		return errs.NewInternal(err)
	}

	isRegisteredToTD, err := tpportal.UserTestDates(
		tpportal.UserTestDateWhere.UserID.EQ(req.UserId),
		tpportal.UserTestDateWhere.TestDateID.EQ(req.TestDateId),
		tpportal.UserTestDateWhere.EducationYear.EQ(user.EducationYear),
	).Exists(ctx, u.st.DBSX())
	if err != nil {
		return errs.NewInternal(err)
	}

	if !isRegisteredToTD {
		return errs.NewBadRequest(errors.New("пользователь не участвовал в выбранном тестировании"))
	}

	exams := tpportal.UserExamResult{
		UserID:        req.UserId,
		TestDateID:    req.TestDateId,
		EducationYear: user.EducationYear,
		RussianLanguageGrade: null.Int{
			Int:   int(req.RussianLanguageGrade.Val),
			Valid: req.RussianLanguageGrade.IsValid,
		},
		MathGrade: null.Int{
			Int:   int(req.MathGrade.Val),
			Valid: req.MathGrade.IsValid,
		},
		ForeignLanguageGrade: null.Int{
			Int:   int(req.ForeignLanguageGrade.Val),
			Valid: req.ForeignLanguageGrade.IsValid,
		},
		FirstProfileGrade: null.Int{
			Int:   int(req.FirstProfileGrade.Val),
			Valid: req.FirstProfileGrade.IsValid,
		},
		SecondProfileGrade: null.Int{
			Int:   int(req.SecondProfileGrade.Val),
			Valid: req.SecondProfileGrade.IsValid,
		},
	}
	err = exams.Upsert(ctx, u.st.DBSX(), true, []string{
		tpportal.UserExamResultColumns.UserID,
		tpportal.UserExamResultColumns.TestDateID,
	}, boil.Infer(), boil.Infer())
	if err != nil {
		return errs.NewInternal(err)
	}
	return nil
}
