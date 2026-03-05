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

func (s *ControllerStorage) CreateProfile(c *gin.Context) {
	response.NewErrorResponse(c, errs.NewNotImplemented())
}

func (s *ControllerStorage) GetProfile(c *gin.Context) {
	response.NewErrorResponse(c, errs.NewNotImplemented())
}

func (s *ControllerStorage) UpdateProfile(c *gin.Context) {
	response.NewErrorResponse(c, errs.NewNotImplemented())
}

func (s *ControllerStorage) ListProfiles(c *gin.Context) {
	profiles, err := s.uc.ListProfiles(c)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, mapper.NewListProfilesResponseToRest(profiles))
}

func (s *ControllerStorage) SetProfilesToUser(c *gin.Context) {
	userIdParam := c.Param("userId")
	userId, err := strconv.ParseInt(userIdParam, 10, 64)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequest(fmt.Errorf("невалидный userId: %s", userIdParam)))
		return
	}
	var req restmodels.SetProfilesToUserRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequest(err))
		return
	}
	err = s.uc.SetProfilesToUser(c, *mapper.NewSetProfileToUserRequestFromRest(&req), userId, false)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	c.Status(http.StatusOK)
}

func (s *ControllerStorage) SetProfilesToMe(c *gin.Context) {
	userIdCtx, ok := c.Get(body.UserCtx)
	if !ok {
		response.NewErrorResponse(c, errs.NewInternal(errors.New("в контексте отсутствует userId")))
		return
	}
	userId := userIdCtx.(tpportal.User).ID
	var req restmodels.SetProfilesToUserRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequest(err))
		return
	}
	err = s.uc.SetProfilesToUser(c, *mapper.NewSetProfileToUserRequestFromRest(&req), userId, true)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	c.Status(http.StatusOK)
}
