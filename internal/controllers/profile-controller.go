package controllers

import (
	"github.com/Limmperhaven/pkportal-be-v2/internal/controllers/response"
	"github.com/Limmperhaven/pkportal-be-v2/internal/errs"
	"github.com/gin-gonic/gin"
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
	response.NewErrorResponse(c, errs.NewNotImplemented())
}

func (s *ControllerStorage) SetProfilesToUser(c *gin.Context) {
	response.NewErrorResponse(c, errs.NewNotImplemented())
}

func (s *ControllerStorage) SetProfilesToMe(c *gin.Context) {
	response.NewErrorResponse(c, errs.NewNotImplemented())
}
