package domain

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/Limmperhaven/pkportal-be-v2/internal/body"
	"github.com/Limmperhaven/pkportal-be-v2/internal/errs"
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/tpportal"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/friendsofgo/errors"
	"github.com/jmoiron/sqlx"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/xuri/excelize/v2"
	"html/template"
	"os"
	"strconv"
	"time"
	"unicode/utf8"
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

func (u *Usecase) GetTestDate(ctx context.Context, tdId int64) (tpportal.TestDateResponse, error) {
	td, err := tpportal.TestDates(
		tpportal.TestDateWhere.ID.EQ(tdId),
		qm.Load(tpportal.TestDateRels.UserTestDates),
	).One(ctx, u.st.DBSX())
	if err != nil {
		if err == sql.ErrNoRows {
			return tpportal.TestDateResponse{}, errs.NewNotFound(fmt.Errorf("дата тестирования с id: %d не найдена", tdId))
		}
		return tpportal.TestDateResponse{}, errs.NewInternal(err)
	}

	tdDate, tdTime := u.formatDateTime(td.DateTime)

	var regPersons, attendedPersons int64
	if td.R.UserTestDates != nil {
		regPersons = int64(len(td.R.UserTestDates))
		for _, utd := range td.R.UserTestDates {
			if utd.IsAttended {
				attendedPersons++
			}
		}
	}

	res := tpportal.TestDateResponse{
		Id:                td.ID,
		Date:              tdDate,
		Time:              tdTime,
		Location:          td.Location,
		AttendedPersons:   attendedPersons,
		RegisteredPersons: regPersons,
		MaxPersons:        int64(td.MaxPersons),
		EducationYear:     int64(td.EducationYear),
		PubStatus:         td.PubStatus.String(),
	}
	return res, nil
}

func (u *Usecase) ListTestDates(ctx context.Context, filter tpportal.ListTestDatesRequest, availableOnly bool) ([]tpportal.TestDateResponse, error) {
	user, err := u.extractUserFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	queryMods := make([]qm.QueryMod, 0, 3)
	queryMods = append(queryMods, qm.OrderBy(tpportal.TestDateColumns.ID))
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
		if utd != nil {
			queryMods = append(queryMods, tpportal.TestDateWhere.ID.NEQ(utd.TestDateID))
		}
	} else {
		if filter.EducationYear != 0 {
			queryMods = append(queryMods,
				tpportal.TestDateWhere.EducationYear.EQ(int16(filter.EducationYear)))
		}
	}

	tds, err := tpportal.TestDates(queryMods...).All(ctx, u.st.DBSX())
	if err != nil {
		return nil, errs.NewInternal(err)
	}

	res := make([]tpportal.TestDateResponse, 0, len(tds))
	for _, td := range tds {
		if availableOnly && (td.MaxPersons == len(td.R.UserTestDates) || td.EducationYear != user.EducationYear) {
			continue
		}
		date, time := u.formatDateTime(td.DateTime)

		var regPersons, attendedPersons int64
		if td.R.UserTestDates != nil {
			regPersons = int64(len(td.R.UserTestDates))
			for _, utd := range td.R.UserTestDates {
				if utd.IsAttended {
					attendedPersons++
				}
			}
		}

		res = append(res, tpportal.TestDateResponse{
			Id:                td.ID,
			Date:              date,
			Time:              time,
			Location:          td.Location,
			RegisteredPersons: regPersons,
			AttendedPersons:   attendedPersons,
			MaxPersons:        int64(td.MaxPersons),
			EducationYear:     int64(td.EducationYear),
			PubStatus:         td.PubStatus.String(),
		})
	}
	return res, nil
}

func (u *Usecase) SignUpUserToTestDate(ctx context.Context, userId, tdId int64) error {
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
		qm.Load(
			qm.Rels(
				tpportal.UserRels.UserForeignLanguages,
			),
		),
		qm.Load(
			qm.Rels(
				tpportal.UserRels.UserTestDates,
			),
		),
	).One(ctx, u.st.DBSX())
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.NewNotFound(fmt.Errorf("пользователь с id: %d не найден", userId))
		}
		return errs.NewInternal(err)
	}

	if len(user.R.UserTestDates) >= 2 {
		return errs.NewBadRequest(errors.New("нельзя записаться более чем на 2 тестирования"))
	}

	firstProfileSet := false
	secondProfileSet := false

	if user.R.UserProfiles != nil {
		valid := false
		for _, up := range user.R.UserProfiles {
			if up.UserEducationYear == user.EducationYear {
				if up.FirstProfileID.Valid {
					firstProfileSet = true
				}
				if up.SecondProfileID.Valid {
					secondProfileSet = true
				}
				valid = true
				break
			}
		}
		if !valid {
			return errs.NewBadRequest(errors.New("для записи на тестирование выберите хотябы один профиль"))
		}
	} else {
		return errs.NewBadRequest(errors.New("для записи на тестирование выберите хотябы один профиль"))
	}
	if user.R.UserProfileSubjects != nil {
		valid := false
		for _, ups := range user.R.UserProfileSubjects {
			if ups.UserEducationYear == user.EducationYear {
				if ups.R.FirstProfileSubject == nil && firstProfileSet {
					return errs.NewBadRequest(errors.New("для записи на тестирование установите предмет 1 профиля"))
				}
				if ups.R.SecondProfileSubject == nil && secondProfileSet {
					return errs.NewBadRequest(errors.New("для записи на тестирование установите предмет 2 профиля"))
				}
				valid = true
				break
			}
		}
		if !valid {
			return errs.NewBadRequest(errors.New("для записи на тестирование выберите профильные предметы для указанных профилей"))
		}
	} else {
		return errs.NewBadRequest(errors.New("для записи на тестирование выберите профильные предметы для указанных профилей"))
	}
	if user.R.UserForeignLanguages != nil {
		valid := false
		for _, ufl := range user.R.UserForeignLanguages {
			if ufl.UserEducationYear == user.EducationYear {
				valid = true
				break
			}
		}
		if !valid {
			return errs.NewBadRequest(errors.New("для записи на тестирование выберите иностранный язык"))
		}
	} else {
		return errs.NewBadRequest(errors.New("для записи на тестирование выберите иностранный язык"))
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

	err = u.st.QueryTx(ctx, func(tx *sqlx.Tx) error {
		err := utd.Insert(ctx, tx, boil.Infer())
		if err != nil {
			return errs.NewInternal(err)
		}
		for _, us := range user.R.UserStatuses {
			if us.EducationYear == user.EducationYear {
				us.StatusID = body.SignedUpForTest.Int64()
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

	var emailBody string
	if user.EducationYear == 9 {
		emailBody = body.SignUpForTestDateMessage9Year
	} else {
		emailBody = body.SignUpForTestDateMessage10Year
	}

	emailMessage := fmt.Sprintf(emailBody, tdDate, td.Location,
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

func (u *Usecase) SetTestDateAttended(ctx context.Context, userId, tdId int64, attendance bool) error {
	user, err := tpportal.FindUser(ctx, u.st.DBSX(), userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.NewNotFound(fmt.Errorf("пользователь с id: %d не найден", userId))
		}
		return errs.NewInternal(err)
	}

	utd, err := tpportal.FindUserTestDate(ctx, u.st.DBSX(), user.ID, tdId)
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.NewNotFound(errors.New("указанная запись на тестирование не найдена"))
		}
		return errs.NewInternal(err)
	}
	if utd.TestDateID != tdId {
		return errs.NewNotFound(errors.New("указанная запись на тестирование не найдена"))
	}
	if utd.IsAttended == attendance {
		return nil
	}
	utd.IsAttended = attendance
	_, err = utd.Update(ctx, u.st.DBSX(), boil.Whitelist(tpportal.UserTestDateColumns.IsAttended))
	if err != nil {
		return errs.NewInternal(err)
	}
	return nil
}

func (u *Usecase) DownloadRegistrationList(ctx context.Context, tdId int64) (tpportal.DownloadFileResponse, error) {
	td, err := tpportal.TestDates(
		tpportal.TestDateWhere.ID.EQ(tdId),
		qm.Load(
			qm.Rels(
				tpportal.TestDateRels.UserTestDates,
				tpportal.UserTestDateRels.User,
				tpportal.UserRels.UserProfiles,
				tpportal.UserProfileRels.FirstProfile,
			),
		),
		qm.Load(
			qm.Rels(
				tpportal.TestDateRels.UserTestDates,
				tpportal.UserTestDateRels.User,
				tpportal.UserRels.UserProfiles,
				tpportal.UserProfileRels.SecondProfile,
			),
		),
		qm.Load(
			qm.Rels(
				tpportal.TestDateRels.UserTestDates,
				tpportal.UserTestDateRels.User,
				tpportal.UserRels.UserProfileSubjects,
				tpportal.UserProfileSubjectRels.FirstProfileSubject,
			),
		),
		qm.Load(
			qm.Rels(
				tpportal.TestDateRels.UserTestDates,
				tpportal.UserTestDateRels.User,
				tpportal.UserRels.UserProfileSubjects,
				tpportal.UserProfileSubjectRels.SecondProfileSubject,
			),
		),
		qm.Load(
			qm.Rels(
				tpportal.TestDateRels.UserTestDates,
				tpportal.UserTestDateRels.User,
				tpportal.UserRels.UserForeignLanguages,
				tpportal.UserForeignLanguageRels.ForeignLanguage,
			),
		),
	).One(ctx, u.st.DBSX())
	if err != nil {
		if err == sql.ErrNoRows {
			return tpportal.DownloadFileResponse{}, errs.NewNotFound(fmt.Errorf("даты тестирования с id: %d не найдено", tdId))
		}
		return tpportal.DownloadFileResponse{}, errs.NewInternal(err)
	}

	tdDate, tdTime := u.formatDateTime(td.DateTime)

	rld := tpportal.RegListData{
		TdDate:     tdDate,
		TdTime:     tdTime,
		TdLocation: td.Location,
	}

	rldu := make([]tpportal.RegListDataUser, 0, len(td.R.UserTestDates))
	if td.R.UserTestDates != nil {
		for _, utd := range td.R.UserTestDates {
			user := utd.R.User

			var firstProfile tpportal.Profile
			var secondProfile tpportal.Profile
			var firstSubject tpportal.Subject
			var secondSubject tpportal.Subject
			var foreignLanguage tpportal.ForeignLanguage

			if user.R.UserProfiles != nil {
				for _, up := range user.R.UserProfiles {
					if up.UserEducationYear == user.EducationYear {
						if up.R.FirstProfile != nil {
							firstProfile = *up.R.FirstProfile
						}
						if up.R.SecondProfile != nil {
							secondProfile = *up.R.SecondProfile
						}
						break
					}
				}
			}
			if user.R.UserProfileSubjects != nil {
				for _, ups := range user.R.UserProfileSubjects {
					if ups.UserEducationYear == user.EducationYear {
						if ups.R.FirstProfileSubject != nil {
							firstSubject = *ups.R.FirstProfileSubject
						}
						if ups.R.SecondProfileSubject != nil {
							secondSubject = *ups.R.SecondProfileSubject
						}
						break
					}
				}
			}
			if user.R.UserForeignLanguages != nil {
				for _, ufls := range user.R.UserForeignLanguages {
					if ufls.UserEducationYear == user.EducationYear {
						if ufls.R.ForeignLanguage != nil {
							foreignLanguage = *ufls.R.ForeignLanguage
						}
						break
					}
				}
			}
			rldu = append(rldu, tpportal.RegListDataUser{
				Id:                   user.ID,
				Fio:                  user.Fio,
				ForeignLanguage:      foreignLanguage.Name,
				FirstProfile:         firstProfile.Name,
				FirstProfileSubject:  firstSubject.Name,
				SecondProfile:        secondProfile.Name,
				SecondProfileSubject: secondSubject.Name,
			})

		}

	}
	rld.Users = rldu

	htmlTemplate, err := os.ReadFile("etc/template.html")
	if err != nil {

	}

	t, err := template.New("reglist").Parse(string(htmlTemplate))
	if err != nil {
		return tpportal.DownloadFileResponse{}, errs.NewInternal(fmt.Errorf("ошибка при создании шаблона html: %s", err.Error()))
	}

	outHtml := new(bytes.Buffer)
	err = t.Execute(outHtml, rld)
	if err != nil {
		return tpportal.DownloadFileResponse{}, errs.NewInternal(fmt.Errorf("ошибка при генерации html: %s", err.Error()))
	}

	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		if err != nil {
			return tpportal.DownloadFileResponse{}, errs.NewInternal(fmt.Errorf("ошибка при инициализации генератора pdf: %s", err.Error()))
		}
	}
	page := wkhtmltopdf.NewPageReader(bytes.NewReader(outHtml.Bytes()))

	page.EnableLocalFileAccess.Set(true)
	pdfg.AddPage(page)

	pdfg.MarginLeft.Set(15)
	pdfg.MarginRight.Set(15)
	pdfg.Dpi.Set(300)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)

	err = pdfg.Create()
	if err != nil {
		if err != nil {
			return tpportal.DownloadFileResponse{}, errs.NewInternal(fmt.Errorf("ошибка при генерации pdf: %s", err.Error()))
		}
	}

	b64File := base64.StdEncoding.EncodeToString(pdfg.Bytes())

	return tpportal.DownloadFileResponse{
		FileName:    "reglist_id" + strconv.Itoa(int(td.ID)),
		FileContent: b64File,
		ContentType: "application/pdf",
	}, nil
}

func (u *Usecase) ExportTestDateToXlsx(ctx context.Context, tdId int64) (tpportal.DownloadFileResponse, error) {
	td, err := tpportal.TestDates(
		tpportal.TestDateWhere.ID.EQ(tdId),
		qm.Load(
			qm.Rels(
				tpportal.TestDateRels.UserTestDates,
				tpportal.UserTestDateRels.User,
				tpportal.UserRels.UserProfiles,
				tpportal.UserProfileRels.FirstProfile,
			),
		),
		qm.Load(
			qm.Rels(
				tpportal.TestDateRels.UserTestDates,
				tpportal.UserTestDateRels.User,
				tpportal.UserRels.UserProfiles,
				tpportal.UserProfileRels.SecondProfile,
			),
		),
		qm.Load(
			qm.Rels(
				tpportal.TestDateRels.UserTestDates,
				tpportal.UserTestDateRels.User,
				tpportal.UserRels.UserProfileSubjects,
				tpportal.UserProfileSubjectRels.FirstProfileSubject,
			),
		),
		qm.Load(
			qm.Rels(
				tpportal.TestDateRels.UserTestDates,
				tpportal.UserTestDateRels.User,
				tpportal.UserRels.UserProfileSubjects,
				tpportal.UserProfileSubjectRels.SecondProfileSubject,
			),
		),
		qm.Load(
			qm.Rels(
				tpportal.TestDateRels.UserTestDates,
				tpportal.UserTestDateRels.User,
				tpportal.UserRels.UserForeignLanguages,
				tpportal.UserForeignLanguageRels.ForeignLanguage,
			),
		),
		qm.Load(
			qm.Rels(
				tpportal.TestDateRels.UserTestDates,
				tpportal.UserTestDateRels.User,
				tpportal.UserRels.UserStatuses,
				tpportal.UserStatusRels.Status,
			),
		),
	).One(ctx, u.st.DBSX())
	if err != nil {
		if err == sql.ErrNoRows {
			return tpportal.DownloadFileResponse{}, errs.NewNotFound(fmt.Errorf("даты тестирования с id: %d не найдено", tdId))
		}
		return tpportal.DownloadFileResponse{}, errs.NewInternal(err)
	}

	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Sheet1"

	sheetIndex, err := f.NewSheet(sheetName)
	if err != nil {
		return tpportal.DownloadFileResponse{}, errs.NewInternal(
			fmt.Errorf("ошибка при создании нового листа xlsx: %s", err.Error()))
	}

	f.SetCellValue("Sheet1", "A1", "ID")
	f.SetCellValue("Sheet1", "B1", "ФИО")
	f.SetCellValue("Sheet1", "C1", "Дата Рождения")
	f.SetCellValue("Sheet1", "D1", "Email")
	f.SetCellValue("Sheet1", "E1", "Профиль 1")
	f.SetCellValue("Sheet1", "F1", "Предмет профиля 1")
	f.SetCellValue("Sheet1", "G1", "Профиль 2")
	f.SetCellValue("Sheet1", "H1", "Предмет профиля 2")
	f.SetCellValue("Sheet1", "I1", "Иностранный язык")
	f.SetCellValue("Sheet1", "J1", "Школа")
	f.SetCellValue("Sheet1", "K1", "Номер телефона")
	f.SetCellValue("Sheet1", "L1", "Номер телефона законного представителя")
	f.SetCellValue("Sheet1", "M1", "Статус")

	if td.R.UserTestDates != nil {
		for i, utd := range td.R.UserTestDates {
			index := i + 2
			user := utd.R.User

			var firstProfile tpportal.Profile
			var secondProfile tpportal.Profile
			var firstSubject tpportal.Subject
			var secondSubject tpportal.Subject
			var foreignLanguage tpportal.ForeignLanguage
			var status tpportal.Status

			if user.R.UserProfiles != nil {
				for _, up := range user.R.UserProfiles {
					if up.UserEducationYear == user.EducationYear {
						if up.R.FirstProfile != nil {
							firstProfile = *up.R.FirstProfile
						}
						if up.R.SecondProfile != nil {
							secondProfile = *up.R.SecondProfile
						}
						break
					}
				}
			}
			if user.R.UserProfileSubjects != nil {
				for _, ups := range user.R.UserProfileSubjects {
					if ups.UserEducationYear == user.EducationYear {
						if ups.R.FirstProfileSubject != nil {
							firstSubject = *ups.R.FirstProfileSubject
						}
						if ups.R.SecondProfileSubject != nil {
							secondSubject = *ups.R.SecondProfileSubject
						}
						break
					}
				}
			}
			if user.R.UserForeignLanguages != nil {
				for _, ufls := range user.R.UserForeignLanguages {
					if ufls.UserEducationYear == user.EducationYear {
						if ufls.R.ForeignLanguage != nil {
							foreignLanguage = *ufls.R.ForeignLanguage
						}
						break
					}
				}
			}

			if user.R.UserStatuses != nil {
				for _, us := range user.R.UserStatuses {
					if us.EducationYear == user.EducationYear {
						status = *us.R.Status
						break
					}
				}
			}

			dob := u.formatDate(user.DateOfBirth)

			f.SetCellValue("Sheet1", "A"+strconv.Itoa(index), user.ID)
			f.SetCellValue("Sheet1", "B"+strconv.Itoa(index), user.Fio)
			f.SetCellValue("Sheet1", "C"+strconv.Itoa(index), dob)
			f.SetCellValue("Sheet1", "D"+strconv.Itoa(index), user.Email)
			f.SetCellValue("Sheet1", "E"+strconv.Itoa(index), firstProfile.Name)
			f.SetCellValue("Sheet1", "F"+strconv.Itoa(index), firstSubject.Name)
			f.SetCellValue("Sheet1", "G"+strconv.Itoa(index), secondProfile.Name)
			f.SetCellValue("Sheet1", "H"+strconv.Itoa(index), secondSubject.Name)
			f.SetCellValue("Sheet1", "I"+strconv.Itoa(index), foreignLanguage.Name)
			f.SetCellValue("Sheet1", "J"+strconv.Itoa(index), user.CurrentSchool.String)
			f.SetCellValue("Sheet1", "K"+strconv.Itoa(index), user.PhoneNumber)
			f.SetCellValue("Sheet1", "L"+strconv.Itoa(index), user.ParentPhoneNumber)
			f.SetCellValue("Sheet1", "M"+strconv.Itoa(index), status.Name)
		}
	}

	cols, err := f.GetCols(sheetName)
	if err != nil {
		return tpportal.DownloadFileResponse{}, errs.NewInternal(
			fmt.Errorf("ошибка при получении столбцов xlsx: %s", err.Error()))
	}
	for idx, col := range cols {
		largestWidth := 0
		for _, rowCell := range col {
			cellWidth := utf8.RuneCountInString(rowCell) + 2
			if cellWidth > largestWidth {
				largestWidth = cellWidth
			}
		}
		name, err := excelize.ColumnNumberToName(idx + 1)
		if err != nil {
			return tpportal.DownloadFileResponse{}, errs.NewInternal(
				fmt.Errorf("ошибка при получении названий столбцов: %s", err.Error()))
		}
		f.SetColWidth("Sheet1", name, name, float64(largestWidth))
	}

	f.SetActiveSheet(sheetIndex)
	buf, err := f.WriteToBuffer()
	if err != nil {
		return tpportal.DownloadFileResponse{}, errs.NewInternal(fmt.Errorf("ошибка при чтении файла: %s", err.Error()))
	}

	b64File := base64.StdEncoding.EncodeToString(buf.Bytes())

	return tpportal.DownloadFileResponse{
		FileName:    "test_date" + strconv.Itoa(int(td.ID)),
		FileContent: b64File,
		ContentType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	}, nil
}
