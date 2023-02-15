package controllers

import (
	"github.com/Limmperhaven/pkportal-be-v2/internal/controllers/response"
	"github.com/Limmperhaven/pkportal-be-v2/internal/errs"
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/mapper"
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/restmodels"
	"github.com/gin-gonic/gin"
	"net/http"
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
	//userIdParam := c.Param("id")
	//userId, err := strconv.ParseInt(userIdParam, 10, 64)
	//if err != nil {
	//	response.NewErrorResponse(c, errs.NewBadRequest(errors.New("Невалидный id пользователя")))
	//	return
	//}
	//
	//s.uc.
	response.NewErrorResponse(c, errs.NewNotImplemented())
}

func (s *ControllerStorage) ListUsers(c *gin.Context) {
	response.NewErrorResponse(c, errs.NewNotImplemented())
}

func (s *ControllerStorage) UpdateUser(c *gin.Context) {
	response.NewErrorResponse(c, errs.NewNotImplemented())
}

func (s *ControllerStorage) GetMe(c *gin.Context) {
	response.NewErrorResponse(c, errs.NewNotImplemented())
}

func (s *ControllerStorage) UpdateMe(c *gin.Context) {
	response.NewErrorResponse(c, errs.NewNotImplemented())
}

func (s *ControllerStorage) ListStatuses(c *gin.Context) {
	response.NewErrorResponse(c, errs.NewNotImplemented())
}

func (s *ControllerStorage) SetStatus(c *gin.Context) {
	response.NewErrorResponse(c, errs.NewNotImplemented())
}
