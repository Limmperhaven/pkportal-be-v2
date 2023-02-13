package controllers

import (
	"github.com/Limmperhaven/pkportal-be-v2/internal/controllers/response"
	"github.com/Limmperhaven/pkportal-be-v2/internal/errs"
	"github.com/gin-gonic/gin"
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
	response.NewErrorResponse(c, errs.NewNotImplemented())
}
