package controllers

import (
	"errors"
	"fmt"
	"github.com/Limmperhaven/pkportal-be-v2/internal/body"
	"github.com/Limmperhaven/pkportal-be-v2/internal/controllers/response"
	"github.com/Limmperhaven/pkportal-be-v2/internal/errs"
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/mapper"
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/restmodels"
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/tpportal"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (s *ControllerStorage) CreateTestDate(c *gin.Context) {
	var req restmodels.CreateTestDateRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequest(err))
		return
	}
	err = s.uc.CreateTestDate(c, *mapper.NewCreateTestDateFromRest(&req))
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	c.Status(http.StatusOK)
}

func (s *ControllerStorage) GetTestDate(c *gin.Context) {
	tdIdParam := c.Param("id")
	tdId, err := strconv.ParseInt(tdIdParam, 10, 64)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequest(fmt.Errorf("невалидный id даты тестирования: %s", tdIdParam)))
		return
	}
	res, err := s.uc.GetTestDate(c, tdId)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, *mapper.NewTestDateResponseToRest(&res))
}

func (s *ControllerStorage) UpdateTestDate(c *gin.Context) {
	response.NewErrorResponse(c, errs.NewNotImplemented())
}

func (s *ControllerStorage) SetTestDatePubStatus(c *gin.Context) {
	tdStatus := c.Param("status")
	tdIdParam := c.Param("id")
	tdId, err := strconv.ParseInt(tdIdParam, 10, 64)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequest(fmt.Errorf("невалидный id даты тестирования: %s", tdIdParam)))
		return
	}
	err = s.uc.SetTestDatePubStatus(c, tdId, tdStatus)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	c.Status(http.StatusOK)
}

func (s *ControllerStorage) ListTestDates(c *gin.Context) {
	var req restmodels.ListTestDatesRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequest(err))
		return
	}
	tds, err := s.uc.ListTestDates(c, *mapper.NewListTestDatesRequestFromRest(&req), false)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, mapper.NewTestDateResponseArrayToRest(tds))
}

func (s *ControllerStorage) ListAvailableTestDates(c *gin.Context) {
	tds, err := s.uc.ListTestDates(c, tpportal.ListTestDatesRequest{}, true)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, mapper.NewTestDateResponseArrayToRest(tds))
}

func (s *ControllerStorage) SignUpUserToTestDate(c *gin.Context) {
	userIdParam := c.Param("userId")
	userId, err := strconv.ParseInt(userIdParam, 10, 64)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequest(fmt.Errorf("невалидный id пользователя: %s", userIdParam)))
		return
	}
	tdIdParam := c.Param("tdId")
	tdId, err := strconv.ParseInt(tdIdParam, 10, 64)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequest(fmt.Errorf("невалидный id даты тестирования: %s", tdIdParam)))
		return
	}
	err = s.uc.SignUpUserToTestDate(c, userId, tdId, false)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	c.Status(http.StatusOK)
}

func (s *ControllerStorage) SignUpMeToTestDate(c *gin.Context) {
	userIdCtx, ok := c.Get(body.UserCtx)
	if !ok {
		response.NewErrorResponse(c, errs.NewInternal(errors.New("в контексте отсутствует userId")))
		return
	}
	userId := userIdCtx.(tpportal.User).ID
	tdIdParam := c.Param("tdId")
	tdId, err := strconv.ParseInt(tdIdParam, 10, 64)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequest(fmt.Errorf("невалидный id даты тестирования: %s", tdIdParam)))
		return
	}
	err = s.uc.SignUpUserToTestDate(c, userId, tdId, true)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	c.Status(http.StatusOK)
}

func (s *ControllerStorage) ListCommonLocations(c *gin.Context) {
	cls, err := s.uc.ListCommonLocations(c)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, mapper.NewIdNameArrayToRest(cls))
}

func (s *ControllerStorage) SetTestDateAttended(c *gin.Context) {
	userIdParam := c.Param("userId")
	userId, err := strconv.ParseInt(userIdParam, 10, 64)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequest(fmt.Errorf("невалидный id пользователя: %s", userIdParam)))
		return
	}
	tdIdParam := c.Param("tdId")
	tdId, err := strconv.ParseInt(tdIdParam, 10, 64)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequest(fmt.Errorf("невалидный id даты тестирования: %s", tdIdParam)))
		return
	}
	err = s.uc.SetTestDateAttended(c, userId, tdId)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	c.Status(http.StatusOK)
}
