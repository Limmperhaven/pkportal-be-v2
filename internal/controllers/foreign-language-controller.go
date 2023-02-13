package controllers

import (
	"github.com/Limmperhaven/pkportal-be-v2/internal/controllers/response"
	"github.com/Limmperhaven/pkportal-be-v2/internal/errs"
	"github.com/gin-gonic/gin"
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
	response.NewErrorResponse(c, errs.NewNotImplemented())
}

func (s *ControllerStorage) SetForeignLanguageToUser(c *gin.Context) {
	response.NewErrorResponse(c, errs.NewNotImplemented())
}

func (s *ControllerStorage) SetForeignLanguageToMe(c *gin.Context) {
	response.NewErrorResponse(c, errs.NewNotImplemented())
}
