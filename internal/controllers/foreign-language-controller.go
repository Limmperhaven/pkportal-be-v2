package controllers

import (
	"errors"
	"fmt"
	"github.com/Limmperhaven/pkportal-be-v2/internal/body"
	"github.com/Limmperhaven/pkportal-be-v2/internal/controllers/response"
	"github.com/Limmperhaven/pkportal-be-v2/internal/errs"
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/mapper"
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/tpportal"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (s *ControllerStorage) CreateForeignLanguage(c *gin.Context) {
	response.NewErrorResponse(c, errs.NewNotImplemented())
}

func (s *ControllerStorage) GetForeignLanguage(c *gin.Context) {
	response.NewErrorResponse(c, errs.NewNotImplemented())
}

func (s *ControllerStorage) UpdateForeignLanguage(c *gin.Context) {
	response.NewErrorResponse(c, errs.NewNotImplemented())
}

func (s *ControllerStorage) ListForeignLanguages(c *gin.Context) {
	fls, err := s.uc.ListForeignLanguages(c)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, mapper.NewIdNameArrayToRest(fls))
}

func (s *ControllerStorage) SetForeignLanguageToUser(c *gin.Context) {
	userIdParam := c.Param("userId")
	userId, err := strconv.ParseInt(userIdParam, 10, 64)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequest(fmt.Errorf("невалидный id пользователя: %s", userIdParam)))
		return
	}
	flIdParam := c.Param("flId")
	flId, err := strconv.ParseInt(flIdParam, 10, 64)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequest(fmt.Errorf("невалидный id иностранного языка: %s", flIdParam)))
		return
	}
	err = s.uc.SetForeignLanguageToUser(c, userId, flId, false)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	c.Status(http.StatusOK)
}

func (s *ControllerStorage) SetForeignLanguageToMe(c *gin.Context) {
	userIdCtx, ok := c.Get(body.UserCtx)
	if !ok {
		response.NewErrorResponse(c, errs.NewInternal(errors.New("в контексте отсутствует userId")))
		return
	}
	userId := userIdCtx.(tpportal.User).ID
	flIdParam := c.Param("flId")
	flId, err := strconv.ParseInt(flIdParam, 10, 64)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequest(fmt.Errorf("невалидный id иностранного языка: %s", flIdParam)))
		return
	}
	err = s.uc.SetForeignLanguageToUser(c, userId, flId, false)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	c.Status(http.StatusOK)
}
