package controllers

import (
	"github.com/Limmperhaven/pkportal-be-v2/internal/controllers/response"
	"github.com/Limmperhaven/pkportal-be-v2/internal/errs"
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/mapper"
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/restmodels"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *ControllerStorage) CreateNotification(c *gin.Context) {
	var req restmodels.Notification
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequest(err))
		return
	}
	err = s.uc.CreateNotification(c, mapper.NewNotificationFromRest(&req))
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	c.Status(http.StatusOK)
}
