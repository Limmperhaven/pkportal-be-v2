package worker

import (
	"context"
	"fmt"
	"github.com/Limmperhaven/pkportal-be-v2/internal/body"
	"github.com/Limmperhaven/pkportal-be-v2/internal/domain"
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/tpportal"
	"github.com/Limmperhaven/pkportal-be-v2/internal/storage/stpg"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"log"
	"strconv"
	"time"
)

func NotificationWorker(ctx context.Context, mail domain.MailClient, s3 domain.S3Client) {
	log.Println("Запущен Notification worker")
	st := stpg.Gist()
	for {
		select {
		case <-time.After(time.Hour):
			tds, err := tpportal.TestDates(
				tpportal.TestDateWhere.DateTime.LTE(time.Now().Add(2*24*time.Hour)),
				tpportal.TestDateWhere.NotificationSent.EQ(false),
				qm.Load(
					qm.Rels(
						tpportal.TestDateRels.UserTestDates,
						tpportal.UserTestDateRels.User,
					),
				),
			).All(ctx, st.DBSX())
			if err != nil {
				log.Printf("ошибка notification worker: %s\n", err.Error())
			}
			if len(tds) == 0 {
				continue
			}

			for _, td := range tds {
				userEmails := make([]string, 0)
				for _, utd := range td.R.UserTestDates {
					userEmails = append(userEmails, utd.R.User.Email)
				}
				if len(userEmails) == 0 {
					continue
				}

				tdDate := td.DateTime.Format("02.01.2006")
				tdTime := strconv.Itoa(td.DateTime.Hour()) + ":" + strconv.Itoa(td.DateTime.Minute())
				emailMessage := fmt.Sprintf(body.BeforeExamMessage, tdDate, td.Location, tdTime)
				err = mail.SendTextEmail(body.BeforeExamMessageSubject, emailMessage, userEmails)
				if err != nil {
					log.Printf("ошибка notification worker: %s\n", err.Error())
				}
				td.NotificationSent = true
				_, err := td.Update(ctx, st.DBSX(), boil.Whitelist(tpportal.TestDateColumns.NotificationSent))
				if err != nil {
					log.Printf("ошибка notification worker: %s\n", err.Error())
				}
				log.Printf("notification worker: отправлены сообщения о тестировании %d абитуриентам\n", len(userEmails))
			}
		case <-ctx.Done():
			log.Println("Notification worker stopped successfully")
			return
		}
	}
}
