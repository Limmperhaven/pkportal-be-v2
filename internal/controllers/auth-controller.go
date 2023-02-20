package controllers

import (
	"github.com/Limmperhaven/pkportal-be-v2/internal/body"
	"github.com/Limmperhaven/pkportal-be-v2/internal/config"
	"github.com/Limmperhaven/pkportal-be-v2/internal/controllers/response"
	"github.com/Limmperhaven/pkportal-be-v2/internal/errs"
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/mapper"
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/restmodels"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *ControllerStorage) SignUp(c *gin.Context) {
	var req restmodels.SignUpRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequest(err))
		return
	}
	err = s.uc.SignUp(c, mapper.NewSignUpRequestFromRest(&req))
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	c.Status(http.StatusOK)
}

func (s *ControllerStorage) SignIn(c *gin.Context) {
	var req restmodels.SignInRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequest(err))
		return
	}
	userAuth, err := s.uc.SignIn(c, mapper.NewSignInRequestFromRest(&req))
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, *mapper.NewSignInResponseToRest(&userAuth))
}

func (s *ControllerStorage) RecoverPassword(c *gin.Context) {
	email := c.Param("email")
	err := s.uc.RecoverPassword(c, email)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	c.Status(http.StatusOK)
}

func (s *ControllerStorage) ActivateAccount(c *gin.Context) {
	token := c.Param("token")
	err := s.uc.Activate(c, token)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	c.Status(http.StatusOK)
}

func (s *ControllerStorage) ConfirmRecover(c *gin.Context) {
	token := c.Param("token")
	var req restmodels.ConfirmRecoverRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequest(err))
		return
	}
	err = s.uc.ConfirmRecover(c, token, req.Password)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	c.Status(http.StatusOK)
}

func (s *ControllerStorage) Logout(c *gin.Context) {
	cfg := config.Get().Server

	c.SetCookie(body.AuthToken,
		"",
		-1,
		"/",
		cfg.Domain,
		cfg.Scheme == "https",
		true)
}
