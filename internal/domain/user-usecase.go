package domain

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Limmperhaven/pkportal-be-v2/internal/body"
	"github.com/Limmperhaven/pkportal-be-v2/internal/config"
	"github.com/Limmperhaven/pkportal-be-v2/internal/errs"
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/tpportal"
	"github.com/friendsofgo/errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
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
	}

	var otherEducationYear int16
	if user.EducationYear == int16(10) {
		otherEducationYear = int16(9)
	} else {
		otherEducationYear = int16(10)
	}

	err = u.st.QueryTx(ctx, func(tx *sqlx.Tx) error {
		err = user.Insert(ctx, tx, boil.Infer())
		if err != nil {
			return errs.NewInternal(err)
		}
		uss := tpportal.UserStatusSlice{
			&tpportal.UserStatus{
				UserID:        user.ID,
				StatusID:      req.StatusId,
				EducationYear: user.EducationYear,
			},
			&tpportal.UserStatus{
				UserID:        user.ID,
				StatusID:      body.Registered.Int64(),
				EducationYear: otherEducationYear,
			},
		}
		err = user.AddUserStatuses(ctx, tx, true, uss...)
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

func (u *Usecase) GetUser(ctx context.Context, userId int64) (tpportal.GetUserResponse, error) {
	user, err := tpportal.Users(
		tpportal.UserWhere.ID.EQ(userId),
		qm.Load(
			qm.Rels(
				tpportal.UserRels.UserStatuses,
				tpportal.UserStatusRels.Status,
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
		qm.Load(
			qm.Rels(
				tpportal.UserRels.UserTestDates,
			),
		),
	).One(ctx, u.st.DBSX())
	if err != nil {
		if err == sql.ErrNoRows {
			return tpportal.GetUserResponse{}, errs.NewNotFound(fmt.Errorf("пользователь с id %d не найден", userId))
		}
		return tpportal.GetUserResponse{}, errs.NewInternal(err)
	}

	status := tpportal.IdName{}
	if len(user.R.UserStatuses) != 0 {
		for _, us := range user.R.UserStatuses {
			if us.EducationYear == user.EducationYear {
				status.Id = us.R.Status.ID
				status.Name = us.R.Status.Name
				break
			}
		}
	}

	firstProfile := tpportal.IdName{}
	secondProfile := tpportal.IdName{}
	if len(user.R.UserProfiles) != 0 {
		for _, up := range user.R.UserProfiles {
			if up.UserEducationYear == user.EducationYear {
				if up.R.FirstProfile != nil {
					firstProfile.Id = up.R.FirstProfile.ID
					firstProfile.Name = up.R.FirstProfile.Name
				}
				if up.R.SecondProfile != nil {
					secondProfile.Id = up.R.SecondProfile.ID
					secondProfile.Name = up.R.SecondProfile.Name
				}
				break
			}
		}
	}

	firstProfileSubject := tpportal.IdName{}
	secondProfileSubject := tpportal.IdName{}
	if len(user.R.UserProfileSubjects) != 0 {
		for _, ups := range user.R.UserProfileSubjects {
			if ups.UserEducationYear == user.EducationYear {
				if ups.R.FirstProfileSubject != nil {
					firstProfileSubject.Id = ups.R.FirstProfileSubject.ID
					firstProfileSubject.Name = ups.R.FirstProfileSubject.Name
				}
				if ups.R.SecondProfileSubject != nil {
					secondProfileSubject.Id = ups.R.SecondProfileSubject.ID
					secondProfileSubject.Name = ups.R.SecondProfileSubject.Name
				}
				break
			}
		}
	}
	foreignLanguage := tpportal.IdName{}
	if len(user.R.UserForeignLanguages) != 0 {
		for _, fl := range user.R.UserForeignLanguages {
			if fl.UserEducationYear == user.EducationYear {
				foreignLanguage.Id = fl.R.ForeignLanguage.ID
				foreignLanguage.Name = fl.R.ForeignLanguage.Name
				break
			}
		}
	}

	testDate := tpportal.GetUserResponseTestDate{}
	if len(user.R.UserTestDates) != 0 {
		for _, utd := range user.R.UserTestDates {
			if utd.EducationYear == user.EducationYear {
				testDate.TestDateId = utd.TestDateID
				testDate.IsAttended = utd.IsAttended
			}
		}
	}

	res := tpportal.GetUserResponse{
		Id:                   user.ID,
		Role:                 user.Role.String(),
		Fio:                  user.Fio,
		DateOfBirth:          u.formatDate(user.DateOfBirth),
		Gender:               user.Gender.String(),
		Email:                user.Email,
		PhoneNumber:          user.PhoneNumber,
		ParentPhoneNumber:    user.ParentPhoneNumber,
		CurrentSchool:        user.CurrentSchool.String,
		EducationYear:        int64(user.EducationYear),
		Status:               status,
		FirstProfile:         firstProfile,
		SecondProfile:        secondProfile,
		FirstProfileSubject:  firstProfileSubject,
		SecondProfileSubject: secondProfileSubject,
		ForeignLanguage:      foreignLanguage,
		TestDate:             testDate,
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
	user, err := tpportal.Users(
		tpportal.UserWhere.ID.EQ(userId),
		qm.Load(
			qm.Rels(
				tpportal.UserRels.UserStatuses,
			),
		),
	).One(ctx, u.st.DBSX())
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.NewNotFound(fmt.Errorf("пользователь с id: %d не найден", userId))
		}
		return errs.NewInternal(err)
	}
	if len(user.R.UserStatuses) == 0 {
		return errs.NewInternal(errors.New("у пользователя не хватает записи о статусе"))
	}
	for _, us := range user.R.UserStatuses {
		if us.EducationYear == user.EducationYear {
			us.StatusID = statusId
			_, err = us.Update(ctx, u.st.DBSX(), boil.Whitelist(tpportal.UserStatusColumns.StatusID))
			if err != nil {
				return errs.NewInternal(err)
			}
			return nil
		}
	}
	return errs.NewNotFound(errors.New("у пользователя не хватает записи о статусе"))
}

func (u *Usecase) UploadScreenshot(ctx context.Context, req tpportal.UploadScreenshotRequest) error {
	user, err := u.extractUserFromCtx(ctx)
	if err != nil {
		return err
	}
	fileNameS3 := uuid.New().String()
	uploadFileReq := tpportal.UploadFileRequest{
		FileKey:     fileNameS3,
		FileSize:    req.FileSize,
		FileContent: req.FileContent,
		ContentType: u.detectContentType(req.FileContent),
	}

	key, err := u.s3.UploadFile(ctx, uploadFileReq)
	if err != nil {
		return err
	}

	us, err := tpportal.UserStatuses(
		tpportal.UserStatusWhere.UserID.EQ(user.ID),
		tpportal.UserStatusWhere.EducationYear.EQ(user.EducationYear),
	).One(ctx, u.st.DBSX())
	if us.StatusID == body.Registered.Int64() {
		us.StatusID = body.AttachedScreenshot.Int64()
	}

	usc := &tpportal.UserScreenshot{
		EducationYear: user.EducationYear,
		OriginalName:  req.FileName,
		FileName:      key,
	}

	err = u.st.QueryTx(ctx, func(tx *sqlx.Tx) error {
		err := user.AddUserScreenshots(ctx, tx, true, usc)
		if err != nil {
			errDel := u.s3.DeleteFile(ctx, key)
			if errDel != nil {
				return errs.NewInternal(fmt.Errorf(
					"ошибка при добавлении файла: %s, ошибка при удалении добавленного файла из хранилища: %s",
					err.Error(), errDel.Error()))
			}
			return errs.NewInternal(fmt.Errorf("ошибка при добавлении файла для пользователя: %s", err.Error()))
		}
		_, err = us.Update(ctx, tx, boil.Whitelist(tpportal.UserStatusColumns.StatusID))
		if err != nil {
			return errs.NewInternal(err)
		}
		return nil
	})
	return err
}

func (u *Usecase) DownloadScreenshot(ctx context.Context, userId int64) (tpportal.DownloadScreenshotResponse, error) {
	user, err := tpportal.Users(
		tpportal.UserWhere.ID.EQ(userId),
		qm.Load(
			qm.Rels(tpportal.UserRels.UserScreenshots),
		),
	).One(ctx, u.st.DBSX())
	if err != nil {
		if err == sql.ErrNoRows {
			return tpportal.DownloadScreenshotResponse{}, errs.NewNotFound(fmt.Errorf("пользователь с id: %d не найден", userId))
		}
		return tpportal.DownloadScreenshotResponse{}, errs.NewInternal(err)
	}

	var fileKey, fileName string
	for _, usc := range user.R.UserScreenshots {
		if usc.EducationYear == user.EducationYear {
			fileKey = usc.FileName
			fileName = usc.OriginalName
			break
		}
	}
	if fileKey == "" {
		return tpportal.DownloadScreenshotResponse{}, errs.NewNotFound(errors.New("скриншот пользователя не найден"))
	}

	fileData, err := u.s3.DownloadFile(ctx, fileKey)
	if err != nil {
		return tpportal.DownloadScreenshotResponse{}, errs.NewInternal(fmt.Errorf("не удалось скачать файл: %s", err.Error()))
	}
	contentType := u.detectContentType(fileData)

	return tpportal.DownloadScreenshotResponse{
		FileName:    fileName,
		FileContent: fileData,
		ContentType: contentType,
	}, nil
}

func (u *Usecase) ListUsers(ctx context.Context, req tpportal.ListUsersRequest) ([]tpportal.GetUserResponse, error) {
	//queryMods := make([]qm.QueryMod, 0)
	//
	//if len(req.EducationYears) != 0 {
	//	educationYears := make([]interface{}, len(req.EducationYears))
	//	for i, ey := range req.EducationYears {
	//		educationYears[i] = ey
	//	}
	//	//queryMods = append(queryMods, qm.WhereIn(tpportal.TableNames.Users+"."+tpportal.UserColumns.EducationYear+" IN ?", educationYears...))
	//}
	//if len(req.ProfileIds) != 0 {
	//	profileIds := make([]interface{}, len(req.ProfileIds))
	//	for i, pi := range req.ProfileIds {
	//		profileIds[i] = pi
	//	}
	//	expr := qm.Expr(
	//		qm.WhereIn(tpportal.TableNames.UserProfiles+"."+tpportal.UserProfileColumns.FirstProfileID+" IN ?", profileIds...),
	//		qm.Or2(qm.WhereIn(tpportal.TableNames.UserProfiles+"."+tpportal.UserProfileColumns.SecondProfileID+" IN ?", profileIds...)),
	//	)
	//	queryMods = append(queryMods, expr)
	//}
	//if len(req.StatusIds) != 0 {
	//
	//}

	//users, err := tpportal.Users().All(ctx, u.st.DBSX())

	educationYears := make([]interface{}, len(req.EducationYears))
	for i, ey := range req.EducationYears {
		educationYears[i] = ey
	}

	users, err := tpportal.Users(
		qm.WhereIn(tpportal.TableNames.Users+"."+tpportal.UserColumns.EducationYear+" IN ?", educationYears...),
		qm.InnerJoin(fmt.Sprintf("%s ON %s = %s AND %s = %s",
			tpportal.TableNames.UserProfiles,
			tpportal.UserColumns.ID,
			tpportal.UserProfileColumns.UserID,
			tpportal.UserColumns.EducationYear,
			tpportal.UserProfileColumns.UserEducationYear,
		)),
	).All(ctx, u.st.DBSX())

	for _, user := range users {
		fmt.Println(len(user.R.UserProfiles))
	}

	return nil, err
}
