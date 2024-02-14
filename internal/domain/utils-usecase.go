package domain

import (
	"context"
	"errors"
	"fmt"
	"github.com/Limmperhaven/pkportal-be-v2/internal/body"
	"github.com/Limmperhaven/pkportal-be-v2/internal/errs"
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/tpportal"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"strings"
	"time"
)

func (u *Usecase) extractUserFromCtx(ctx context.Context) (tpportal.User, error) {
	ginC, ok := ctx.(*gin.Context)
	if !ok {
		return tpportal.User{}, errs.NewInternal(errors.New("ошибка в преобразовании контекста"))
	}
	userCtx, ok := ginC.Get(body.UserCtx)
	if !ok {
		return tpportal.User{}, errs.NewInternal(errors.New("пользователь отсутствует в контексте"))
	}
	user, ok := userCtx.(tpportal.User)
	if !ok {
		return tpportal.User{}, errs.NewInternal(errors.New("невалидный формат пользователя в контексте"))
	}
	return user, nil
}

func (u *Usecase) parseDate(in string) (time.Time, error) {
	date, err := time.Parse("02.01.2006", in)
	if err != nil {
		return time.Time{}, errs.NewBadRequest(fmt.Errorf("невалидная дата: %s", in))
	}
	return date, nil
}

func (u *Usecase) hashPassword(in string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(in+body.AppSalt), body.AppCost)
	if err != nil {
		return "", errs.NewInternal(err)
	}
	return string(hash), nil
}

func (u *Usecase) formatDate(in time.Time) string {
	return in.Format("02.01.2006")
}

func (u *Usecase) parseDateTime(dateString string, timeString string) (dateTime time.Time, err error) {
	date, err := time.Parse("02.01.2006", dateString)
	if err != nil {
		return time.Time{}, errs.NewBadRequest(fmt.Errorf("невалидная дата: %s", dateString))
	}

	timeStringParts := strings.Split(timeString, ":")
	if len(timeStringParts) != 2 {
		return time.Time{}, errs.NewBadRequest(fmt.Errorf("невалидное время: %s", dateString))
	}
	hours, err := strconv.Atoi(timeStringParts[0])
	if err != nil {
		return time.Time{}, errs.NewBadRequest(fmt.Errorf("невалидное время: %s", dateString))
	}
	minutes, err := strconv.Atoi(timeStringParts[1])
	if err != nil {
		return time.Time{}, errs.NewBadRequest(fmt.Errorf("невалидное время: %s", dateString))
	}

	dateTime = time.Date(date.Year(), date.Month(), date.Day(), hours, minutes, 0, 0, time.Local)
	return dateTime, nil
}

func (u *Usecase) formatDateTime(in time.Time) (dateString, timeString string) {
	dateString = in.Format("02.01.2006")
	var minuteString string
	if in.Minute() < 10 {
		minuteString = "0" + strconv.Itoa(in.Minute())
	} else {
		minuteString = strconv.Itoa(in.Minute())
	}
	timeString = strconv.Itoa(in.Hour()) + ":" + minuteString
	return dateString, timeString
}

func (u *Usecase) detectContentType(data []byte) string {
	mimetype.SetLimit(0)
	return mimetype.Detect(data).String()
}

func (u *Usecase) int64SliceToInterfaceSlice(in []int64) []interface{} {
	out := make([]interface{}, len(in))
	for i, v := range in {
		out[i] = v
	}
	return out
}
