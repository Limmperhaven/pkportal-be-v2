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

func (s *ControllerStorage) CreateSubject(c *gin.Context) {
	response.NewErrorResponse(c, errs.NewNotImplemented())
}

func (s *ControllerStorage) GetSubject(c *gin.Context) {
	response.NewErrorResponse(c, errs.NewNotImplemented())
}

func (s *ControllerStorage) UpdateSubject(c *gin.Context) {
	response.NewErrorResponse(c, errs.NewNotImplemented())
}

func (s *ControllerStorage) ListSubjects(c *gin.Context) {
	var req restmodels.ListSubjectsRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequest(err))
		return
	}
	subjects, err := s.uc.ListSubjects(c, req.ProfileId)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, mapper.NewIdNameArrayToRest(subjects))
}

func (s *ControllerStorage) SetSubjectToUser(c *gin.Context) {
	userIdParam := c.Param("userId")
	userId, err := strconv.ParseInt(userIdParam, 10, 64)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequest(fmt.Errorf("невалидный id пользователя: %s", userIdParam)))
		return
	}
	var req restmodels.SetSubjectsRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequest(err))
		return
	}
	err = s.uc.SetSubjectsToUser(c, *mapper.NewSetSubjectsRequestFromRest(&req), userId, false)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	c.Status(http.StatusOK)
}

func (s *ControllerStorage) SetSubjectToMe(c *gin.Context) {
	userIdCtx, ok := c.Get(body.UserCtx)
	if !ok {
		response.NewErrorResponse(c, errs.NewInternal(errors.New("в контексте отсутствует userId")))
		return
	}
	userId := userIdCtx.(tpportal.User).ID
	var req restmodels.SetSubjectsRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequest(err))
		return
	}
	err = s.uc.SetSubjectsToUser(c, *mapper.NewSetSubjectsRequestFromRest(&req), userId, true)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	c.Status(http.StatusOK)
}
