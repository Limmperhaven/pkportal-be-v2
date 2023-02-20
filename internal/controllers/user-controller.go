package controllers

import (
	"fmt"
	"github.com/Limmperhaven/pkportal-be-v2/internal/body"
	"github.com/Limmperhaven/pkportal-be-v2/internal/controllers/response"
	"github.com/Limmperhaven/pkportal-be-v2/internal/errs"
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/mapper"
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/restmodels"
	"github.com/friendsofgo/errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (s *ControllerStorage) CreateUser(c *gin.Context) {
	var req restmodels.CreateUserRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequest(err))
		return
	}
	err = s.uc.CreateUser(c, mapper.NewCreateUserRequestFromRest(&req))
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	c.Status(http.StatusOK)
}

func (s *ControllerStorage) GetUser(c *gin.Context) {
	userIdParam := c.Param("id")
	userId, err := strconv.ParseInt(userIdParam, 10, 64)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequest(fmt.Errorf("невалидный userId: %s", userIdParam)))
		return
	}
	user, err := s.uc.GetUser(c, userId)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	res := mapper.NewGetUserResponseToRest(&user)
	c.JSON(http.StatusOK, *res)
}

func (s *ControllerStorage) ListUsers(c *gin.Context) {
	response.NewErrorResponse(c, errs.NewNotImplemented())
}

func (s *ControllerStorage) UpdateUser(c *gin.Context) {
	userIdParam := c.Param("id")
	userId, err := strconv.ParseInt(userIdParam, 10, 64)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequest(fmt.Errorf("невалидный userId: %s", userIdParam)))
		return
	}
	var req restmodels.UpdateUserRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequest(err))
		return
	}
	err = s.uc.UpdateUser(c, *mapper.NewUpdateUserRequestFromRest(&req), userId)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	c.Status(http.StatusOK)
}

func (s *ControllerStorage) GetMe(c *gin.Context) {
	userIdCtx, ok := c.Get(body.UserId)
	if !ok {
		response.NewErrorResponse(c, errs.NewInternal(errors.New("в контексте отсутствует userId")))
		return
	}
	userId := userIdCtx.(int64)
	user, err := s.uc.GetUser(c, userId)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	res := mapper.NewGetUserResponseToRest(&user)
	c.JSON(http.StatusOK, *res)
}

func (s *ControllerStorage) ListStatuses(c *gin.Context) {
	var req restmodels.ListStatusesRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequest(err))
		return
	}
	statuses, err := s.uc.ListStatuses(c, *mapper.NewListStatusesRequestFromRest(&req))
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, mapper.NewIdNameArrayToRest(statuses))
}

func (s *ControllerStorage) SetUserStatus(c *gin.Context) {
	userIdParam := c.Param("userId")
	statusIdParam := c.Param("statusId")
	userId, err := strconv.ParseInt(userIdParam, 10, 64)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequest(fmt.Errorf("невалидный userId: %s", userIdParam)))
		return
	}
	statusId, err := strconv.ParseInt(statusIdParam, 10, 64)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequest(fmt.Errorf("невалидный statusId: %s", userIdParam)))
		return
	}
	err = s.uc.SetUserStatus(c, userId, statusId)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	c.Status(http.StatusOK)
}
