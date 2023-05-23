package domain

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Limmperhaven/pkportal-be-v2/internal/errs"
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/tpportal"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (u *Usecase) CreateNotification(ctx context.Context, req *tpportal.Notification) error {
	users, err := tpportal.Users(
		qm.WhereIn("id IN (?)", u.int64SliceToInterfaceSlice(req.UserIds)...),
	).All(ctx, u.st.DBSX())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return errs.NewInternal(err)
	}
	userEmails := make([]string, len(users))
	for i, user := range users {
		userEmails[i] = user.Email
	}
	fmt.Println(userEmails)
	err = u.mail.SendTextEmail(req.Topic, req.Message, userEmails)
	if err != nil {
		return errs.NewInternal(err)
	}
	return nil
}
