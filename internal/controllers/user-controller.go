package controllers

import (
	"github.com/Limmperhaven/pkportal-be-v2/internal/controllers/response"
	"github.com/Limmperhaven/pkportal-be-v2/internal/errs"
	"github.com/gin-gonic/gin"
)

func (s *ControllerStorage) CreateUser(c *gin.Context) {
	response.NewErrorResponse(c, errs.NewNotImplemented())
}

func (s *ControllerStorage) GetUser(c *gin.Context) {
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
